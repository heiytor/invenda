package cache

import (
	"context"
	"time"

	rediscache "github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type SetOption func(i *rediscache.Item)

func WithTTL(ttl time.Duration) SetOption {
	return func(i *rediscache.Item) {
		i.TTL = ttl
	}
}

type Cache interface {
	// Get gets the cache value for the given key. Missing keys is not treated as an error.
	Get(ctx context.Context, key string, value interface{}) error

	// Set puts the key with the specified value into cache.
	Set(ctx context.Context, key string, value interface{}, opts ...SetOption) error

	// Delete deletes cached value with the specified key.
	Delete(ctx context.Context, key string) error
}

type cache struct {
	cache *rediscache.Cache
}

var _ Cache = (*cache)(nil) // Ensures that cache implements Cache

func Connect(uri string, pool int) (Cache, error) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	if pool > 0 {
		opt.PoolSize = pool
	}

	cache := &cache{
		cache: rediscache.New(&rediscache.Options{
			Redis: redis.NewClient(opt),
		}),
	}

	return cache, nil
}

func (c *cache) Get(ctx context.Context, key string, value interface{}) error {
	err := c.cache.Get(ctx, key, value)
	if err == rediscache.ErrCacheMiss {
		return nil
	}

	return err
}

func (c *cache) Set(ctx context.Context, key string, value interface{}, opts ...SetOption) error {
	i := &rediscache.Item{Ctx: ctx, Key: key, Value: value}
	for _, opt := range opts {
		opt(i)
	}

	return c.cache.Set(i)
}

func (c *cache) Delete(ctx context.Context, key string) error {
	if err := c.cache.Get(ctx, key, nil); err == rediscache.ErrCacheMiss {
		return nil
	}

	return c.cache.Delete(ctx, key)
}
