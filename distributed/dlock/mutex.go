package dlock

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/adnpa/go-utils/pkg/data/redis"
	"github.com/hashicorp/go-multierror"
	"time"
)

// 默认实现的随机获取value的方法
func genValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// DelayFunc 设置重试时间
type DelayFunc func(tries int) time.Duration

// Mutex 分布式锁对象
type Mutex struct {
	name   string
	expiry time.Duration //有效时间
	until  time.Time     //过期时间

	tries     int //最大重试次数
	delayFunc DelayFunc

	driftFactor   float64
	timeoutFactor float64

	quorum int //法定人数（获取锁的最小节点数量）

	genValueFunc func() (string, error)
	value        string //唯一标识符

	shuffle       bool //访问连接池时随机打乱顺序，以减少热点问题
	failFast      bool //快速失败/快速返回
	setNXOnExtend bool //延长锁的有效期时 如果锁已过期 尝试重新设置锁

	pools []redis.Pool
}

// Name returns mutex name (i.e. the Redis key).
func (m *Mutex) Name() string {
	return m.name
}

// Value returns the current random value. The value will be empty until a lock is acquired (or WithValue option is used).
func (m *Mutex) Value() string {
	return m.value
}

// Until returns the time of validity of acquired lock. The value will be zero value until a lock is acquired.
func (m *Mutex) Until() time.Time {
	return m.until
}

// TryLock only attempts to lock m once and returns immediately regardless of success or failure without retrying.
func (m *Mutex) TryLock() error {
	return m.TryLockContext(context.Background())
}

// TryLockContext only attempts to lock m once and returns immediately regardless of success or failure without retrying.
func (m *Mutex) TryLockContext(ctx context.Context) error {
	return m.lockContext(ctx, 1)
}

// Lock locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *Mutex) Lock() error {
	return m.LockContext(context.Background())
}

// LockContext locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *Mutex) LockContext(ctx context.Context) error {
	return m.lockContext(ctx, m.tries)
}

func (m *Mutex) lockContext(ctx context.Context, tries int) error {
	if ctx == nil {
		ctx = context.Background()
	}

	value, err := m.genValueFunc()
	if err != nil {
		return err
	}

	var timer *time.Timer
	for i := 0; i < tries; i++ {

		//重试 退避
		if i != 0 {
			if timer == nil {
				timer = time.NewTimer(m.delayFunc(i))
			} else {
				timer.Reset(m.delayFunc(i))
			}

			select {
			case <-ctx.Done():
				timer.Stop()
				//提前提出 获取失败
				return ErrFailed
			case <-timer.C:
				//延迟
			}
		}

		//获取锁
		start := time.Now()
		n, err := func() (int, error) {
			ctx, cancel := context.WithTimeout(ctx, time.Duration(int64(float64(m.expiry)*m.timeoutFactor)))
			defer cancel()
			return m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
				return m.acquire(ctx, pool, value)
			})
		}()
		now := time.Now()

		until := now.Add(m.expiry - now.Sub(start) - time.Duration(int64(float64(m.expiry)*m.timeoutFactor)))
		if n >= m.quorum && now.Before(until) {
			m.value = value
			m.until = until
			return nil
		}

		//释放自身占有的锁 tries为1是尝试获取锁 不会
		_, _ = func() (int, error) {
			ctx, cancel := context.WithTimeout(ctx, time.Duration(int64(float64(m.expiry)*m.timeoutFactor)))
			defer cancel()
			return m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
				return m.release(ctx, pool, value)
			})
		}()
		if i == tries-1 && err != nil {
			return err
		}

	}

	return m.lockContext(ctx, m.tries)
}

// 尝试锁
func (m *Mutex) acquire(ctx context.Context, pool redis.Pool, value string) (bool, error) {
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	reply, err := conn.SetNX(m.name, value, m.expiry)
	if err != nil {
		return false, err
	}
	return reply, nil
}

// Unlock unlocks m and returns the status of unlock.
func (m *Mutex) Unlock() (bool, error) {
	return m.UnlockContext(context.Background())
}

// UnlockContext unlocks m and returns the status of unlock.
func (m *Mutex) UnlockContext(ctx context.Context) (bool, error) {
	n, err := m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
		return m.release(ctx, pool, m.value)
	})
	if n < m.quorum {
		return false, err
	}
	return true, nil
}

// 释放锁
func (m *Mutex) release(ctx context.Context, pool redis.Pool, value string) (bool, error) {
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	status, err := conn.Eval(deleteScript, m.name, value)
	if err != nil {
		return false, err
	}
	if status == int64(-1) {
		return false, ErrLockAlreadyExpired
	}
	return status != int64(0), nil
}

// Extend resets the mutex's expiry and returns the status of expiry extension.
func (m *Mutex) Extend() (bool, error) {
	return m.ExtendContext(context.Background())
}

func (m *Mutex) ExtendContext(ctx context.Context) (bool, error) {
	start := time.Now()
	n, err := m.actOnPoolsAsync(func(pool redis.Pool) (bool, error) {
		return m.touch(ctx, pool, m.value, int(m.expiry/time.Millisecond))
	})
	if n < m.quorum {
		return false, err
	}
	now := time.Now()
	until := now.Add(m.expiry - now.Sub(start) - time.Duration(int64(float64(m.expiry)*m.driftFactor)))
	if now.Before(until) {
		m.until = until
		return true, nil
	}
	return false, ErrExtendFailed
}

func (m *Mutex) touch(ctx context.Context, pool redis.Pool, value string, expiry int) (bool, error) {
	conn, err := pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	touchScript := touchScript
	if m.setNXOnExtend {
		touchScript = touchWithSetNXScript
	}

	status, err := conn.Eval(touchScript, m.name, value, expiry)
	if err != nil {
		return false, err
	}
	return status != int64(0), nil
}

// 在多个 Redis 连接池上并发执行 并收集结果
func (m *Mutex) actOnPoolsAsync(actFn func(redis.Pool) (bool, error)) (int, error) {
	type result struct {
		node     int
		statusOK bool
		err      error
	}

	ch := make(chan result, len(m.pools))
	for node, pool := range m.pools {
		go func(node int, pool redis.Pool) {
			r := result{node: node}
			r.statusOK, r.err = actFn(pool)
			ch <- r
		}(node, pool)
	}

	var (
		n     = 0
		taken []int
		err   error
	)

	for range m.pools {
		r := <-ch
		if r.statusOK {
			n++
		} else if errors.Is(r.err, ErrLockAlreadyExpired) {
			err = multierror.Append(err, ErrLockAlreadyExpired)
		} else if r.err != nil {
			err = multierror.Append(err, &RedisError{Node: r.node, Err: r.err})
		} else {
			taken = append(taken, r.node)
			err = multierror.Append(err, &ErrNodeTaken{Node: r.node})
		}

		if m.failFast {
			// fast return
			if n >= m.quorum {
				return n, err
			}

			// fail fast
			if len(taken) >= m.quorum {
				return n, &ErrTaken{Nodes: taken}
			}
		}
	}

	if len(taken) >= m.quorum {
		return n, &ErrTaken{Nodes: taken}
	}
	return n, err
}

// lua scripts -----------------------------------------------------------
// 释放锁
var deleteScript = redis.NewScript(1, `
	local val = redis.call("GET", KEYS[1])
	if val == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	elseif val == false then
		return -1
	else
		return 0
	end
`)

// 延长锁
var touchScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end
`)

// 延长锁 过期时重新设置
// nx确保没有其他客户端持有锁
var touchWithSetNXScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	elseif redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2], "NX") then
		return 1
	else
		return 0
	end
`)
