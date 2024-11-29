package cache

import (
	"context"
	"errors"
	"github.com/adnpa/go-utils/pkg/basic"
	"github.com/adnpa/go-utils/pkg/data/cache/local"
	"github.com/redis/go-redis/v9"
	"reflect"
	"time"
)

type (
	MarshalFunc   func(v interface{}) ([]byte, error)
	UnmarshalFunc func(data []byte, v interface{}) error
)

type redisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd

	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type Options struct {
	Redis        redisClient
	LocalCache   local.Local
	StatsEnabled bool

	Marshal   MarshalFunc
	Unmarshal UnmarshalFunc
}

type Item struct {
	Ctx context.Context
	Key string
	Val interface{}
	TTL time.Duration

	Do             func(*Item) (interface{}, error)
	SetXX          bool
	SetNX          bool
	SkipLocalCache bool
}

// Cache ===============================================================

type Cache struct {
	opt *Options

	marshal   MarshalFunc
	unmarshal UnmarshalFunc

	hits uint64
	miss uint64
}

func (c *Cache) Get(ctx context.Context, key string, val interface{}) error {
	return c.get(ctx, key, val, false)
}

func (c *Cache) GetSkipLocalCache(ctx context.Context, key string, val interface{}) error {
	return c.get(ctx, key, val, true)
}

func (c *Cache) get(ctx context.Context, key string, val interface{}, skipLocalCache bool) error {
	localCache := c.opt.LocalCache
	if !skipLocalCache && localCache != nil {
		result, ok := localCache.Get(key)
		if ok {
			err := setVal(val, result)
			if err == nil {
				return nil
			}
		}
	}

	redisClient := c.opt.Redis
	if redisClient == nil {
		return errors.New("redis client is nil")
	}

	bytes, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	err = basic.UnmarshalGob(bytes, val)
	return err
}

func (c *Cache) Set(item *Item) error {
	_, err := c.set(item)
	return err
}

func (c *Cache) set(item *Item) (bool, error) {
	if c.opt.LocalCache != nil && !item.SkipLocalCache {
		c.opt.LocalCache.Set(item.Key, item.Val)
	}

	if c.opt.Redis == nil {
		return true, nil
	}

	bytes, err := basic.MarshalGob(item.Val)
	if err != nil {
		return false, err
	}

	if item.SetXX {
		return true, c.opt.Redis.SetXX(item.Ctx, item.Key, bytes, item.TTL).Err()
	}
	if item.SetNX {
		return true, c.opt.Redis.SetNX(item.Ctx, item.Key, bytes, item.TTL).Err()
	}
	return true, c.opt.Redis.Set(item.Ctx, item.Key, bytes, item.TTL).Err()
}

func (c *Cache) Del(ctx context.Context, key string) error {
	if c.opt.LocalCache != nil {
		c.opt.LocalCache.Del(key)
	}
	if c.opt.Redis == nil {
		return errors.New("redis client is nil")
	}
	return c.opt.Redis.Del(ctx, key).Err()
}

func setVal(target interface{}, data interface{}) error {
	val := reflect.ValueOf(target)
	if val.Type().Kind() != reflect.Pointer {
		return errors.New("gob: attempt to decode into a non-pointer")
	}

	if val.IsValid() {
		if val.Kind() == reflect.Pointer && !val.IsNil() {
		} else if !val.CanSet() {
			return errors.New("not settable")
		}
	}

	val.Set(reflect.ValueOf(data))
	return nil
}

// stat start-----------------------------------------------------------

type Stats struct {
	Hits uint64
	Miss uint64
}

func (c *Cache) Stats() *Stats {
	if c.opt.StatsEnabled {
		return &Stats{Hits: c.hits, Miss: c.miss}
	}
	return nil
}

// stat end-----------------------------------------------------------
