package cacher

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"sort"
	"time"

	"github.com/adverax/echo/data"
	"github.com/adverax/echo/generic"
)

type Storage interface {
	// Lock key
	Lock(key string)
	// Unlock key
	Unlock(key string)

	// Get value by key. Returns data.ErrNoMatch, if has no key.
	Get(key string, dst interface{}) error
	// Set value with key and expire time.
	Set(key string, val interface{}, timeout time.Duration) error
	// Delete cached value by key.
	Delete(key string) error

	// Assert depended key
	Assert(key string, dependencies map[string]string) error
	// Invalidate depended data
	Invalidate(key, val string) error
}

type Template interface {
	Execute(w io.Writer, data interface{}) error
}

type Cacher interface {
	FetchData(
		class string,
		dependencies map[string]string,
		dst interface{},
		builder func() (interface{}, error),
		lifeTime time.Duration,
	) error

	FetchHtml(
		class string,
		dependencies map[string]string,
		builder func() (tpl Template, data interface{}, err error),
		lifeTime time.Duration,
	) (html string, err error)

	// Invalidate depended data
	Invalidate(key, val string) error
}

type cacher struct {
	Storage
}

func (c *cacher) FetchData(
	class string,
	dependencies map[string]string,
	dst interface{},
	builder func() (interface{}, error),
	lifeTime time.Duration,
) error {
	key := c.makeKey(class, dependencies)

	c.Lock(key)
	defer c.Unlock(key)

	err := c.Storage.Get(key, dst)
	if err != data.ErrNoMatch {
		return err
	}

	val, err := builder()
	if err != nil {
		return err
	}

	generic.CloneValueTo(dst, val)

	err = c.Storage.Set(key, val, lifeTime)
	if err != nil {
		return err
	}

	return c.Assert(key, dependencies)
}

func (c *cacher) FetchHtml(
	class string,
	dependencies map[string]string,
	builder func() (tpl Template, data interface{}, err error),
	lifeTime time.Duration,
) (html string, err error) {
	err = c.FetchData(
		class,
		dependencies,
		&html,
		func() (interface{}, error) {
			tpl, params, err := builder()
			if err != nil {
				return nil, err
			}

			var buf bytes.Buffer
			err = tpl.Execute(&buf, params)
			if err != nil {
				return nil, err
			}

			return buf.String(), nil
		},
		lifeTime,
	)
	return
}

// Create normalized key from dependencies
func (c *cacher) makeKey(
	class string,
	dependencies map[string]string,
) string {
	keys := make([]string, len(dependencies))
	for key := range dependencies {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	hasher := md5.New()
	hasher.Write([]byte(class))
	for _, key := range keys {
		val := dependencies[key]
		item := key + "=" + val + ";"
		hasher.Write([]byte(item))
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func New(
	storage Storage,
) Cacher {
	return &cacher{
		Storage: storage,
	}
}
