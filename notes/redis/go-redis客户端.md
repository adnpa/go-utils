



https://redis.io/docs/latest/develop/connect/clients/go/

https://redis.uptrace.dev/guide/

https://github.com/redis/go-redis

https://redis.io/docs/latest/develop/connect/clients/go/





redis客户端

* https://github.com/redis/go-redis
  redis官方客户端
* https://github.com/redis/rueidis
  快速的 Golang Redis 客户端，支持客户端缓存、自动流水线、泛型 OM、RedisJSON、RedisBloom、RediSearch
* https://github.com/gomodule/redigo

* https://github.com/gomodule/redigo/tree/46992b0f02f74066bcdfd9b03e33bc03abd10dc7



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

Redis 集群是 Redis 提供的一种分布式解决方案，可以将数据分散到多个 Redis 实例中，以实现水平扩展和高可用性。通过 Redis 集群，用户可以在多个节点之间分配数据，从而提高性能和容量，同时减少单点故障的风险。

### Redis 集群的特性

1. **数据分片**: Redis 集群将数据分散到多个节点中，使用哈希槽（hash slots）来管理数据分布。每个 Redis 实例负责一部分哈希槽。
2. **高可用性**: Redis 集群支持主从复制，每个主节点可以有多个从节点，从节点可以在主节点故障时接管服务。
3. **自动故障转移**: 如果主节点失败，集群会自动选择一个从节点升为主节点，确保服务的持续可用性。
4. **无中心化**: Redis 集群是去中心化的，所有节点都是平等的，没有单点故障。
5. **透明化的客户端**: Redis 客户端通常会处理集群中的节点，自动路由请求到正确的节点。

### Redis 集群的架构

在 Redis 集群中，通常有以下组件：

- **主节点（Master）**: 负责处理写请求，并将数据存储在哈希槽中。
- **从节点（Replica）**: 从主节点复制数据，处理读请求，并在主节点故障时接管。



## Redis Sentinel client

Redis 哨兵模式（Sentinel）是 Redis 提供的一种高可用性解决方案，旨在监控 Redis 主从架构，确保在主节点发生故障时能够自动进行故障转移（failover）。哨兵模式提供了监控、通知和自动故障转移等功能，使 Redis 集群更加可靠和稳定。

### 哨兵模式的主要功能

1. **监控**: 哨兵会定期检查主节点和从节点的状态，以确保它们正常运行。
2. **故障检测**: 如果主节点不可用，哨兵会检测到这一故障，并标记主节点为“下线”（subjectively down）。
3. **自动故障转移**: 当确认主节点故障后，哨兵会自动选举一个从节点升级为新的主节点，并通知其他从节点进行重新配置。
4. **通知系统**: 哨兵可以通过 API 向外部系统发送故障通知，以便进行相应的处理。
5. **配置管理**: 哨兵可以提供 Redis 集群的当前配置，包括主节点和从节点的信息。

### 哨兵模式的架构

在哨兵模式中，通常有以下组件：

- **主节点（Master）**: 负责处理写请求。
- **从节点（Replica）**: 负责处理读请求并从主节点进行数据复制。
- **哨兵节点（Sentinel）**: 监控主从节点的状态，执行故障转移。

### 





## Redis Ring client

Ring 分片客户端，是采用了一致性 HASH 算法在多个 redis 服务器之间分发 key，每个节点承担一部分 key 的存储。



## Redis Univeral client

`UniversalClient` 并不是一个客户端，而是对 `Client` 、 `ClusterClient` 、 `FailoverClient` 客户端的包装。



## redis集群

```go
client := redis.NewClusterClient(&redis.ClusterOptions{
    Addrs: []string{":16379", ":16380", ":16381", ":16382", ":16383", ":16384"},

    // To route commands by latency or randomly, enable one of the following.
    //RouteByLatency: true,
    //RouteRandomly: true,
})
```












