package memory

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/adverax/echo/cache"
	"github.com/adverax/echo/sync/arbiter"
	"sync"
	"time"
)

type Manager interface {
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
	// Invalidate dependencies
	Invalidate(key, val string) error
}

type engine struct {
	mx sync.Mutex
	arbiter.Arbiter
	cache.Cache
	deps map[string][]string
}

func (engine *engine) Assert(
	key string,
	dependencies map[string]string,
) error {
	engine.mx.Lock()
	defer engine.mx.Unlock()

	for k, v := range dependencies {
		s := makeKey(k, v)
		if deps, ok := engine.deps[s]; ok {
			engine.deps[s] = append(deps, key)
		} else {
			engine.deps[s] = []string{key}
		}
	}

	return nil
}

func (engine *engine) Invalidate(
	key, val string,
) error {
	engine.mx.Lock()
	defer engine.mx.Unlock()

	id := makeKey(key, val)
	if deps, ok := engine.deps[id]; ok {
		for _, key := range deps {
			err := engine.Cache.Delete(key)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func makeKey(key, val string) string {
	s := key + "=" + val
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func New(
	arbiter arbiter.Arbiter,
	cache cache.Cache,
) Manager {
	return &engine{
		Arbiter: arbiter,
		Cache:   cache,
		deps:    make(map[string][]string, 4096),
	}
}
