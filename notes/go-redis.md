



https://redis.io/docs/latest/develop/connect/clients/go/

https://redis.uptrace.dev/guide/

https://github.com/redis/go-redis

https://redis.io/docs/latest/develop/connect/clients/go/









# 缓存原理

缓存原指CPU和内存间的SRAM，这种存储介质速度介于CPU和内存（DRAM）间，能很好的解决CPU和内存速度差距过大的问题。后来引申为凡是位于速度相差较大的两种硬件之间，用于协调两者数据传输速度差异的结构，均可称之为Cache

缓存能加快存储是因为程序访问内存具有**局部性**

1. 时间局部性(Temporal Locality):
   - 程序在一段时间内会重复访问某些特定的内存区域。
   - 这意味着如果某个数据或指令被访问过,在短时间内它很可能会被再次访问。

2. 空间局部性(Spatial Locality):
   - 程序倾向于访问与刚访问过的内存地址相邻的内存区域。
   - 这意味着如果某个数据或指令被访问过,附近的数据或指令很可能会被接下来访问。

局部性原理是计算机系统设计的基础。它使得缓存技术可以有效地提高系统性能。例如,当CPU访问某个内存地址时,不仅会将该地址的数据加载到缓存中,还会将附近的数据也加载进缓存。这样可以充分利用空间局部性,减少CPU访问主存的次数,提高系统效率。

随着时代发展，除了常见的三级缓存，在分布式等领域也开始使用多级缓存的架构。

# 多级缓存架构

多级缓存架构是一种常见的缓存设计模式,它通过将缓存划分为不同层级,以提高整体缓存系统的性能和可靠性。

在多级缓存架构中,通常包括以下几个缓存层级:

1. **本地缓存(Local Cache)**: 位于应用程序或客户端设备的本地内存中,访问延迟最低。

2. **分布式缓存(Distributed Cache)**: 位于独立的缓存服务器集群中,如 Redis、Memcached 等,提供更大的缓存容量和更高的可用性。

3. **数据库缓存(Database Cache)**: 利用数据库自身的缓存机制,如 MySQL 的查询缓存,缓存频繁执行的数据库查询。

4. **CDN 缓存(CDN Cache)**: 利用内容分发网络(CDN)缓存静态资源,如图片、CSS、JS 等,提高访问速度。

这些缓存层级的特点如下:

- 本地缓存访问速度最快,但容量有限,适合缓存热点数据。
- 分布式缓存容量更大,但访问延迟稍高,适合缓存较热的数据。
- 数据库缓存能减轻数据库的查询压力,但命中率相对较低。
- CDN 缓存适合缓存静态资源,可以大幅提高资源访问速度。

多级缓存的工作流程如下:

1. 应用程序首先查询本地缓存。
2. 如果本地缓存未命中,则查询分布式缓存。
3. 如果分布式缓存也未命中,则查询数据库,并将结果缓存到分布式缓存。
4. 对于静态资源,还可以利用 CDN 缓存进一步提高访问速度。

通过多级缓存架构,可以充分利用不同层级缓存的优势,提高整体系统的性能和可扩展性。同时,多级缓存也能提高系统的可用性和容错能力。



# 本地缓存

本地缓存是一种常见的缓存机制,它将数据缓存在应用程序所在的本地服务器或客户端设备上,以提高数据访问的效率和响应速度。相比于远程缓存(如Redis、Memcached等),本地缓存具有以下优势:

1. **访问延迟低**: 数据存储在本地,可以快速访问,无需网络传输的开销。

2. **独立性强**: 不依赖于外部缓存服务,可以独立运行,提高系统的可靠性。

3. **成本低**: 无需部署和维护额外的缓存服务,成本相对较低。

4. **适合小规模应用**: 对于访问量较小的应用场景,本地缓存已经足够满足需求。

常见的本地缓存实现方式包括:

1. **内存缓存**: 将数据存储在应用程序的内存中,如 Java 中的 HashMap、Guava Cache 等。

2. **文件缓存**: 将数据缓存到本地文件系统中,可以利用操作系统的文件 I/O 缓存机制。

3. **嵌入式数据库**: 使用嵌入式数据库(如 SQLite、LevelDB)作为缓存存储。

4. **缓存库**: 使用专门的缓存库,如 Caffeine、Ehcache 等,提供更丰富的缓存特性。

本地缓存适用于以下场景:

1. **热点数据访问**: 频繁访问的数据可以缓存在本地,减轻后端服务器的压力。

2. **离线应用**: 移动应用或嵌入式设备可以利用本地缓存来支持离线功能。

3. **微服务架构**: 各个微服务可以利用本地缓存来提高内部服务之间的访问速度。

4. **临时数据存储**: 一些临时性的数据可以存储在本地缓存中,避免频繁访问数据库。

总之,本地缓存是一种简单高效的缓存机制,适用于各种应用场景。它可以与分布式缓存(如Redis)结合使用,形成多级缓存架构,进一步提高系统的性能和可扩展性。



## 分布式缓存

分布式缓存是一种在分布式系统中使用的缓存技术。它具有以下特点:

1. 分布式架构:
   - 分布式缓存系统由多个缓存节点组成,这些节点分布在不同的服务器上。
   - 通过水平扩展,可以增加缓存的总容量和处理能力。

2. 高可用性:
   - 即使某个缓存节点故障,其他节点仍然可以提供服务。
   - 通过数据复制等机制,可以保证数据的高可用性。

3. 负载均衡:
   - 分布式缓存系统可以根据当前负载情况,自动将请求分配到不同的缓存节点。
   - 这可以提高整体的吞吐量和响应速度。

4. 数据一致性:
   - 分布式缓存通常会采用缓存同步机制,保证缓存数据的一致性。
   - 比如主动更新、延迟更新等策略。

常见的分布式缓存系统有:

- Redis
- Memcached
- Hazelcast
- Infinispan
- Couchbase

这些系统通常采用键值对的数据结构,并提供丰富的数据操作API。分布式缓存可以显著提高系统的性能和可用性,是构建高性能Web应用的重要技术之一。



# go实现固定大小本地缓存（lru）

参考：https://github.com/hashicorp/golang-lru

哈希表+双向链表的实现

## Entry结构

```go
type Entry[K comparable, V any] struct {
	next, prev *Entry[K, V]
	list *LruList[K, V]
	Key K
	Value V

	ExpiresAt time.Time    //可选的淘汰时间
    ExpireBucket uint8 // 管理数据库中键过期的数据结构
}
```

## 缓存管理结构

哈希表是key-ptr，用于快速查找数据。

双向链表记录数据的使用情况，每次使用后放在链表头，淘汰时淘汰链表尾。

```go
type LRU[K comparable, V any] struct {
	size      int //缓存容量
	evictList *internal.LruList[K, V]  // 数据存储的 双向链表
	items     map[K]*internal.Entry[K, V] //key-ptr map
	onEvict   EvictCallback[K, V]  //淘汰回调函数
}
```

## 接口

```go
type LRUCache[K comparable, V any] interface {
	// Adds a value to the cache, returns true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Add(key K, value V) bool

	// Returns key's value from the cache and
	// updates the "recently used"-ness of the key. #value, isFound
	Get(key K) (value V, ok bool)

	// Checks if a key exists in cache without updating the recent-ness.
	Contains(key K) (ok bool)

	// Returns key's value without updating the "recently used"-ness of the key.
	Peek(key K) (value V, ok bool)

	// Removes a key from the cache.
	Remove(key K) bool

	// Removes the oldest entry from cache.
	RemoveOldest() (K, V, bool)

	// Returns the oldest entry from the cache. #key, value, isFound
	GetOldest() (K, V, bool)

	// Returns a slice of the keys in the cache, from oldest to newest.
	Keys() []K

	// Values returns a slice of the values in the cache, from oldest to newest.
	Values() []V

	// Returns the number of items in the cache.
	Len() int

	// Clears all cache entries.
	Purge()

	// Resizes cache, returning number evicted
	Resize(int) int
}
```

## lru逻辑实现

* 查找元素放在链表头
* 插入元素时如果不存在且插入后达到缓存最大容量，淘汰链表尾

```go
func (c *LRU[K, V]) Get(key K) (value V, ok bool) {
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value, true
	}
	return
}


func (c *LRU[K, V]) Add(key K, value V) (evicted bool) {
	// Check for existing item
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value = value
		return false
	}

	// Add new item
	ent := c.evictList.PushFront(key, value)
	c.items[key] = ent

	evict := c.evictList.Length() > c.size
	// Verify size not exceeded
	if evict {
		c.removeOldest()
	}
	return evict
}

func (c *LRU[K, V]) removeOldest() {
	if ent := c.evictList.Back(); ent != nil {
		c.removeElement(ent)
	}
}
```

## 使用该库

```go
```



































# go-redis api



## 连接和使用

```go
package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

    // 操作字符串
	ctx := context.Background()

	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)

    // 操作map
	session := map[string]string{"name": "John", "surname": "Smith", "company": "Redis", "age": "29"}
	for k, v := range session {
		err := client.HSet(ctx, "user-session:123", k, v).Err()
		if err != nil {
			panic(err)
		}
	}

	userSession := client.HGetAll(ctx, "user-session:123").Val()
	fmt.Println(userSession)
}
```



## ssh连接

```go
//使用ssh连接远程redis服务
	key, err := os.ReadFile("C:\\Users\\hz\\.ssh\\id_ed25519")
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal(err)
	}
	sshConfig := &ssh.ClientConfig{
		User:            "admin",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}
	sshCli, err := ssh.Dial("tcp", "47.109.148.21:22", sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:           cfg.DB,
		Password:     cfg.Password,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		Dialer: func(ctx context2.Context, network, addr string) (net.Conn, error) {
			return sshCli.Dial(network, addr)
		},
		ReadTimeout:  -2, //点进去看disables SetWriteDeadline的值 和版本有关 网上说-1
		WriteTimeout: -2,
	})
```



访问远程redis服务

```
ssh -L local_port:remote_redis_host:remote_redis_port user@remote_server_host

ssh -L 12001:47.109.148.21:6379 admin@47.109.148.21
```



## Redis Cluster client



## Redis Sentinel client



## Redis Ring client



## Redis Univeral client





## redis集群

```go
client := redis.NewClusterClient(&redis.ClusterOptions{
    Addrs: []string{":16379", ":16380", ":16381", ":16382", ":16383", ":16384"},

    // To route commands by latency or randomly, enable one of the following.
    //RouteByLatency: true,
    //RouteRandomly: true,
})
```







## 热key问题





## 客户端库

* https://github.com/redis/go-redis
  redis官方客户端
* https://github.com/redis/rueidis
  快速的 Golang Redis 客户端，支持客户端缓存、自动流水线、泛型 OM、RedisJSON、RedisBloom、RediSearch
* https://github.com/gomodule/redigo
  
* https://github.com/gomodule/redigo/tree/46992b0f02f74066bcdfd9b03e33bc03abd10dc7





# 应用

## 共享session

通常我们在开发后台管理系统时，会使用 Session 来保存用户的会话(登录)状态，这些 Session 信息会被保存在服务器端，但这只适用于单系统应用，如果是分布式系统此模式将不再适用。

例如用户一的 Session 信息被存储在服务器一，但第二次访问时用户一被分配到服务器二，这个时候服务器并没有用户一的 Session 信息，就会出现需要重复登录的问题，问题在于分布式系统每次会把请求随机分配到不同的服务器。

因此，我们需要借助 Redis 对这些 Session 信息进行统一的存储和管理，这样无论请求发送到那台服务器，服务器都会去同一个 Redis 获取相关的 Session 信息，这样就解决了分布式系统下 Session 存储的问题。





## 缓存

* https://github.com/DiceDB/dice
  兼容 redis 的、反应灵敏的、可扩展的、高可用性的**统一缓存**，针对现代硬件进行了优化。
* https://github.com/eko/gocache
   一个完整的 Go 缓存库，为您提供多种管理缓存的方式
* https://github.com/go-redis/cache
  小项目

## 缓存雪崩、击穿、穿透

好的,让我们来详细介绍一下 Redis 中常见的几种缓存问题:缓存雪崩、缓存击穿和缓存穿透。

1. **缓存雪崩**:
   - 缓存雪崩指的是在某一时刻发生大规模的缓存失效,导致大量请求打到数据库,引发数据库崩溃的情况。
   - 这通常发生在缓存服务器重启或大量缓存key在某一时间点集中失效的情况下。
   - 解决方案包括:
     - 利用Redis集群提高缓存可用性
     - 为缓存设置合理的过期时间,加入随机因素
     - 给关键数据设置永不过期的缓存

2. **缓存击穿**:
   - 缓存击穿指的是一个存在但在某个时间点刚好失效的 key,这个时候所有的请求都会打到数据库。
   - 这通常发生在访问一些热点数据的情况下。
   - 解决方案包括:
     - 利用互斥锁(Mutex)保证同一时间只有一个线程去查询数据库
     - 设置热点数据永不过期或者过期时间很长

3. **缓存穿透**:
   - 缓存穿透指的是客户端请求的数据在缓存和数据库中都不存在,这将导致每次请求都要打到数据库。
   - 这种情况通常发生在客户端输入一些无意义的数据,从而查询不到结果。
   - 解决方案包括:
     - 对查不到的key设置默认值缓存
     - 采用布隆过滤器等方式提前保留一些可能存在的key
     - 限制频率过高的请求

总的来说,缓存问题的解决方案需要从缓存设计、服务架构等多个层面进行优化。对于特定的业务场景,可以采取不同的预防和补救措施。

关键是要提前预测可能出现的缓存问题,并制定相应的应对策略。通过合理的缓存设计和系统架构,可以大大提高系统的稳定性和可用性。



## 数据库和缓存如何保证一致性

数据库和缓存的一致性是一个需要权衡的问题。通常有以下几种常见的解决方案:

1. **Cache Aside（旁路缓存）**:
   - 这是最简单的方式,先操作缓存,然后再操作数据库。
   - 写操作先更新数据库,然后**删除缓存**;读操作先查缓存,缓存没有再查数据库并更新缓存。
   - 优点是实现简单,可以快速响应客户端请求。
   - 缺点是需要处理缓存失效导致的数据不一致问题。

2. **Read/Write Through**:
   - 读写都先访问缓存,缓存中没有再访问数据库。
   - **写操作先更新数据库,然后更新缓存。**
   - 优点是缓存中的数据和数据库中的数据保持一致。
   - 缺点是响应时间会稍微长一些,因为需要访问数据库。

3. **Write Behind（异步缓存）**:
   - **写操作先更新缓存,然后异步更新数据库。**
   - 优点是可以提高写入性能,且能保证最终一致性。
   - 缺点是需要处理数据库写入失败的情况,且需要复杂的事务管理。

4. **Refresh Ahead**:
   - 缓存数据在临近过期时,提前从数据库加载新数据到缓存。
   - 优点是可以避免缓存过期导致的缓存击穿问题。
   - 缺点是需要根据业务特点**预测**数据的访问模式。

除此之外,还可以使用分布式锁、二级缓存、Cache Writeback等技术来解决一致性问题。

选择哪种方案需要根据具体的业务场景和性能需求来权衡。一般来说,对于对实时性要求不高但对一致性要求高的业务,可以选择Read/Write Through;而对于对实时性要求高但对一致性要求相对较低的业务,可以选择Cache Aside。

总的来说,保证数据库和缓存的一致性需要采取多种技术手段,并结合具体业务特点进行权衡和选择。







## 分布式缓存







## 多级缓存







## 秒杀



## 分布式锁

https://github.com/go-redsync/redsync

### 为什么用Redis

Redis 作为一个中央存储，用于保存锁的状态。它的高性能和原子操作使其成为实现分布式锁的理想选择。

### 向 Redis 请求的流程

以下是实现分布式锁的典型流程：

1. **获取锁**：

   - 节点向 Redis 发送请求，通常使用 `SET` 命令设置一个键（例如锁的名称），并附加一个唯一的值（例如 UUID）和过期时间。

     ```
     SET lock_name unique_value EX 10 NX
     ```

   - 这里，`EX 10` 设置锁的过期时间为 10 秒，`NX` 表示仅在键不存在时设置。

2. **检查锁的状态**：

   - 如果返回值为 `OK`，则表示锁成功获取；如果返回 `nil`，则表示锁已被其他节点占用。

3. **执行临界区操作**：

   - 一旦锁被成功获取，节点可以执行对共享资源的操作。

4. **释放锁**：

   - 操作完成后，节点必须释放锁。释放锁的过程通常涉及**检查当前锁的值是否与节点持有的一致，然后使用 `DEL` 命令删除锁**。
   - 因为涉及多个操作，所以用lua脚本

5. **处理失败与重试**：

   - 如果节点在获取锁时失败，通常会实现重试机制，可能会使用指数退避策略来延迟重试，避免对 Redis 的频繁请求。

6. **延长锁时间**：

   * PEXPIRE

### 关键点

- **原子性**：Redis 的操作是原子的，使用 Lua 脚本可以确保获取和释放锁的过程不会被其他操作干扰。
- **超时机制**：为防止死锁，设置过期时间是非常重要的。即使节点因故障未能释放锁，锁也会在超时后自动释放。
- **一致性**：确保在分布式环境中，所有节点能够正确地获取和释放锁，保持数据的一致性。

