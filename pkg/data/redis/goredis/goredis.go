package goredis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

import myredis "github.com/adnpa/go-utils/pkg/data/redis"

//import redis "github.com/adnpa/go-utils/pkg/data/redis"

//type

type pool struct {
	delegate    *redis.Client
	withContext bool
	ctx         context.Context
}

func (p *pool) Get(ctx context.Context) (myredis.Conn, error) {
	c := p.delegate
	if ctx != nil {
		p.ctx = ctx
		p.withContext = true
	}
	return &conn{c}, nil
}

type conn struct {
	delegate *redis.Client
	ctx      context.Context
}

func (c *conn) Get(name string) (string, error) {
	value, err := c.delegate.Get(c.ctx, name).Result()
	return value, noErrNil(err)
}

func (c *conn) Set(name, value string) (bool, error) {
	err := c.delegate.Set(c.ctx, name, value, 0).Err()
	return err == nil, nil
}

func (c *conn) SetNX(name string, value string, expiry time.Duration) (bool, error) {
	err := c.delegate.SetEx(c.ctx, name, value, expiry).Err()
	return err == nil, nil
}

func (c *conn) Eval(script *myredis.Script, keysAndArgs ...interface{}) (interface{}, error) {
	keys := make([]string, script.KeyCount)
	args := keysAndArgs

	if script.KeyCount > 0 {
		for i := 0; i < script.KeyCount; i++ {
			keys[i] = keysAndArgs[i].(string)
		}

		args = keysAndArgs[script.KeyCount:]
	}

	//使用脚本的 SHA1 哈希值执行已缓存的 Lua 脚本
	result, err := c.delegate.EvalSha(c.ctx, script.Hash, keys, args).Result()
	if err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT") {
		//执行指定的 Lua 脚本
		result, err = c.delegate.Eval(c.ctx, script.Src, keys, args...).Result()
	}
	return result, noErrNil(err)
}

func (c *conn) PTTL(name string) (time.Duration, error) {
	result, err := c.delegate.PTTL(c.ctx, name).Result()
	return result, noErrNil(err)
}

func (c *conn) Close() func(c *conn) error {
	return nil
}

func noErrNil(err error) error {
	if !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}
