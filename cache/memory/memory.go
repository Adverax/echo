// Copyright 2019 Adverax. All Rights Reserved.
// This file is part of project
//
//      http://github.com/adverax/echo
//
// Licensed under the MIT (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://github.com/adverax/echo/blob/master/LICENSE
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memory

import (
	"container/list"
	"fmt"
	"github.com/adverax/echo/data"
	"hash/fnv"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type Options struct {
	// Max size of cache
	MaxSize int64
	// Bucket count (2 ^ n). Default 4
	Buckets uint8
	// The number of items to prune when memory is low. Default 500
	ItemsToPrune uint32
	// The size of the queue for items which should be deleted. If the queue fills
	// up, calls to Delete() will block. Default 1024
	DeleteBuffer uint32
	// The size of the queue for items which should be promotes. If the queue fills
	// up, calls to Get/Set/GetMulti/IsExists() will block. Default 1024
	PromoteBuffer uint32
	// Give a large cache with a high read / write ratio, it's usually unnecessary
	// to promote an item on every Get. GetsPerPromote specifies the number of Gets
	// a key must have before being promoted. Default 3
	GetsPerPromote int32
}

type Cache struct {
	Options
	list        *list.List
	size        int64
	buckets     []*bucket
	bucketMask  uint32
	deletables  chan *entry
	promotables chan *entry
	donec       chan struct{}
}

func New(options Options) *Cache {
	if options.MaxSize == 0 {
		options.MaxSize = 16000000
	}

	if options.Buckets == 0 {
		options.Buckets = 4
	}

	if options.ItemsToPrune == 0 {
		options.ItemsToPrune = 500
	}

	if options.DeleteBuffer == 0 {
		options.DeleteBuffer = 1024
	}

	if options.PromoteBuffer == 0 {
		options.PromoteBuffer = 1024
	}

	if options.GetsPerPromote == 0 {
		options.GetsPerPromote = 3
	}

	c := &Cache{
		list:       list.New(),
		Options:    options,
		bucketMask: uint32(options.Buckets) - 1,
		buckets:    make([]*bucket, options.Buckets),
	}

	for i := 0; i < int(options.Buckets); i++ {
		c.buckets[i] = &bucket{
			lookup: make(map[string]*entry),
		}
	}

	c.restart()
	return c
}

// Get an item from the cache. Returns nil if the item wasn't found or expired.
func (c *Cache) Get(key string, dst interface{}) error {
	item := c.get(key)
	if item == nil {
		return data.ErrNoMatch
	}

	assign(dst, item.value)
	return nil
}

// Get multiple values from cache.
func (c *Cache) GetMulti(dict map[string]interface{}) (notFound []string, err error) {
	for key, val := range dict {
		item := c.get(key)
		if item != nil {
			assign(val, item.value)
		} else {
			notFound = append(notFound, key)
		}
	}
	return nil, nil
}

// Set the value in the cache for the specified duration
func (c *Cache) Set(key string, value interface{}, duration time.Duration) error {
	c.set(key, value, duration)
	return nil
}

// Remove the item from the cache.
func (c *Cache) Delete(key string) error {
	item := c.bucket(key).delete(key)
	if item != nil {
		c.deletables <- item
	}
	return nil
}

// Check if value exists or not.
func (c *Cache) IsExists(key string) (bool, error) {
	item := c.get(key)
	if item == nil {
		return false, nil
	}
	return !item.expired(), nil
}

// Increase cached int value by key, as a counter.
func (c *Cache) Increase(key string) error {
	return c.bucket(key).increase(key)
}

// Decrease cached int value by key, as a counter.
func (c *Cache) Decrease(key string) error {
	return c.bucket(key).decrease(key)
}

// This isn't thread safe. It's meant to be called from non-concurrent tests
func (c *Cache) Clear() error {
	for _, bucket := range c.buckets {
		bucket.clear()
	}
	c.size = 0
	c.list = list.New()
	return nil
}

// Stops the background worker. Operations performed on the cache after Stop
// is called are likely to panic
func (c *Cache) Stop() {
	close(c.promotables)
	<-c.donec
}

func (c *Cache) restart() {
	c.deletables = make(chan *entry, c.DeleteBuffer)
	c.promotables = make(chan *entry, c.PromoteBuffer)
	c.donec = make(chan struct{})
	go c.worker()
}

func (c *Cache) deleteItem(bucket *bucket, item *entry) {
	bucket.delete(item.key) //stop other GETs from getting it
	c.deletables <- item
}

func (c *Cache) get(key string) *entry {
	item := c.bucket(key).get(key)
	if item == nil {
		return nil
	}
	if item.expires > time.Now().UnixNano() {
		c.promote(item)
		return item
	}
	return nil
}

func (c *Cache) set(key string, value interface{}, duration time.Duration) *entry {
	item, existing := c.bucket(key).set(key, value, duration)
	if existing != nil {
		c.deletables <- existing
	}
	c.promote(item)
	return item
}

func (c *Cache) bucket(key string) *bucket {
	h := fnv.New32a()
	h.Write([]byte(key))
	i := h.Sum32() & c.bucketMask
	return c.buckets[i]
}

func (c *Cache) promote(item *entry) {
	c.promotables <- item
}

func (c *Cache) worker() {
	defer close(c.donec)

	for {
		select {
		case item, ok := <-c.promotables:
			if ok == false {
				goto drain
			}
			if c.doPromote(item) && c.size > c.MaxSize {
				c.gc()
			}
		case item := <-c.deletables:
			c.doDelete(item)
		}
	}

drain:
	for {
		select {
		case item := <-c.deletables:
			c.doDelete(item)
		default:
			close(c.deletables)
			return
		}
	}
}

func (c *Cache) doDelete(item *entry) {
	if item.element == nil {
		item.promotions = -2
	} else {
		c.size -= item.size
		c.list.Remove(item.element)
	}
}

func (c *Cache) doPromote(item *entry) bool {
	//already deleted
	if item.promotions == -2 {
		return false
	}
	if item.element != nil { //not a new item
		if item.needPromote(c.GetsPerPromote) {
			c.list.MoveToFront(item.element)
			item.promotions = 0
		}
		return false
	}

	c.size += item.size
	item.element = c.list.PushFront(item)
	return true
}

func (c *Cache) gc() {
	element := c.list.Back()
	for i := 0; uint32(i) < c.ItemsToPrune; i++ {
		if element == nil {
			return
		}
		prev := element.Prev()
		item := element.Value.(*entry)
		c.bucket(item.key).delete(item.key)
		c.size -= item.size
		c.list.Remove(element)
		item.promotions = -2
		element = prev
	}
}

type bucket struct {
	sync.RWMutex
	lookup map[string]*entry
}

func (b *bucket) get(key string) *entry {
	b.RLock()
	res, _ := b.lookup[key]
	b.RUnlock()
	return res
}

func (b *bucket) set(key string, value interface{}, duration time.Duration) (*entry, *entry) {
	expires := time.Now().Add(duration).UnixNano()
	item := newItem(key, value, expires)
	b.Lock()
	existing, _ := b.lookup[key]
	b.lookup[key] = item
	b.Unlock()
	return item, existing
}

func (b *bucket) increase(key string) error {
	b.Lock()
	defer b.Unlock()

	entry, ok := b.lookup[key]
	if !ok {
		return fmt.Errorf("key %q does dot exists", key)
	}

	switch val := entry.value.(type) {
	case int:
		entry.value = val + 1
	case int8:
		entry.value = val + 1
	case int16:
		entry.value = val + 1
	case int32:
		entry.value = val + 1
	case int64:
		entry.value = val + 1
	case uint:
		entry.value = val + 1
	case uint8:
		entry.value = val + 1
	case uint16:
		entry.value = val + 1
	case uint32:
		entry.value = val + 1
	case uint64:
		entry.value = val + 1
	default:
		return fmt.Errorf("invalid type of cache value %q", key)
	}

	return nil
}

func (b *bucket) decrease(key string) error {
	b.Lock()
	defer b.Unlock()

	entry, ok := b.lookup[key]
	if !ok {
		return fmt.Errorf("key %q does dot exists", key)
	}

	switch val := entry.value.(type) {
	case int:
		entry.value = val - 1
	case int8:
		entry.value = val - 1
	case int16:
		entry.value = val - 1
	case int32:
		entry.value = val - 1
	case int64:
		entry.value = val - 1
	case uint:
		entry.value = val - 1
	case uint8:
		entry.value = val - 1
	case uint16:
		entry.value = val - 1
	case uint32:
		entry.value = val - 1
	case uint64:
		entry.value = val - 1
	default:
		return fmt.Errorf("invalid type of cache value %q", key)
	}

	return nil
}

func (b *bucket) delete(key string) *entry {
	b.Lock()
	item, _ := b.lookup[key]
	delete(b.lookup, key)
	b.Unlock()
	return item
}

func (b *bucket) clear() {
	b.Lock()
	b.lookup = make(map[string]*entry)
	b.Unlock()
}

type Sized interface {
	Size() int64
}

type entry struct {
	key        string
	group      string
	promotions int32
	expires    int64
	size       int64
	value      interface{}
	element    *list.Element
}

func newItem(key string, value interface{}, expires int64) *entry {
	size := int64(1)
	if sized, ok := value.(Sized); ok {
		size = sized.Size()
	}

	return &entry{
		key:        key,
		value:      value,
		promotions: 0,
		size:       size,
		expires:    expires,
	}
}

func (i *entry) needPromote(getsPerPromote int32) bool {
	i.promotions++
	return i.promotions == getsPerPromote
}

func (i *entry) expired() bool {
	expires := atomic.LoadInt64(&i.expires)
	return expires < time.Now().UnixNano()
}

func assign(dst, src interface{}) {
	s := reflect.ValueOf(src)
	d := reflect.ValueOf(dst)
	d.Elem().Set(s)
}
