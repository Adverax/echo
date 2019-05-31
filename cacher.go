package echo

import (
	"bytes"
	"github.com/adverax/echo/cache"
	"github.com/adverax/echo/data"
	"github.com/adverax/echo/generic"
	"github.com/adverax/echo/sync/arbiter"
	"time"
)

type Cacher interface {
	FetchData(
		key string,
		dst interface{},
		builder func() (interface{}, error),
		lifeTime time.Duration,
	) error

	FetchHtml(
		key string,
		builder func() (tpl Template, data interface{}, err error),
		lifeTime time.Duration,
	) (html string, err error)
}

type cacher struct {
	arbiter.Arbiter
	cache cache.Cache
}

func (c *cacher) FetchData(
	key string,
	dst interface{},
	builder func() (interface{}, error),
	lifeTime time.Duration,
) error {
	c.Lock(key)
	defer c.Unlock(key)

	err := c.cache.Get(key, dst)
	if err != data.ErrNoMatch {
		return err
	}

	val, err := builder()
	if err != nil {
		return err
	}

	generic.CloneValueTo(dst, val)

	return c.cache.Set(key, val, lifeTime)
}

func (c *cacher) FetchHtml(
	key string,
	builder func() (tpl Template, data interface{}, err error),
	lifeTime time.Duration,
) (html string, err error) {
	err = c.FetchData(
		key,
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

func NewCacher(
	arbiter arbiter.Arbiter,
	cache cache.Cache,
) Cacher {
	return &cacher{
		Arbiter: arbiter,
		cache:   cache,
	}
}
