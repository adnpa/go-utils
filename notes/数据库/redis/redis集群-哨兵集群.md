https://redis.io/docs/latest/operate/oss_and_stack/management/sentinel/

为非集群（分片） Redis 提供高可用性



## 拓扑结构

| 节点 | IP            | PORT  |
| ---- | ------------- | ----- |
| s1   | 192.168.1.130 | 27001 |
| s2   | 192.168.1.130 | 27002 |
| s3   | 192.168.1.130 | 27003 |



## 准备配置

**新建配置文件**

27001

```
port 27001
sentinel announce-ip 192.168.1.130
sentinel monitor mymaster 192.168.1.130 7001 2
sentinel down-after-millseconds mymaster 5000
sentinel failover-timeout mymaster 60000
dir "/home/hz/test-sentinel/s1"
```

27002

```
port 27002
sentinel announce-ip 192.168.1.130
sentinel monitor mymaster 192.168.1.130 7001 2
sentinel down-after-millseconds mymaster 5000
sentinel failover-timeout mymaster 60000
dir "/home/hz/test-sentinel/s2"
```

27003

```
port 27003
sentinel announce-ip 192.168.1.130
sentinel monitor mymaster 192.168.1.130 7001 2
sentinel down-after-millseconds mymaster 5000
sentinel failover-timeout mymaster 60000
dir "/home/hz/test-sentinel/s3"
```



## 启动

```sh
redis-sentinel s1/sentinel.conf
redis-sentinel s2/sentinel.conf
redis-sentinel s3/sentinel.conf

redis-server s1/sentinel.conf --sentinel
redis-server s2/sentinel.conf --sentinel
redis-server s3/sentinel.conf --sentinel
```



## 测试

关闭节点，自动将从节点提升为master

```
3053:X 30 Nov 2024 15:37:41.021 # +switch-master mymaster 192.168.1.130 7003 192.168.1.130 7002
3053:X 30 Nov 2024 15:37:41.022 * +slave slave 192.168.1.130:7001 192.168.1.130 7001 @ mymaster 192.168.1.130 7002
3053:X 30 Nov 2024 15:37:41.023 * +slave slave 192.168.1.130:7003 192.168.1.130 7003 @ mymaster 192.168.1.130 7002
3053:X 30 Nov 2024 15:37:41.024 * Sentinel new configuration saved on disk
3053:X 30 Nov 2024 15:37:51.038 # +sdown slave 192.168.1.130:7003 192.168.1.130 7003 @ mymaster 192.168.1.130 7002
2883:X 30 Nov 2024 15:41:09.172 # +sdown sentinel fe3c9743cffcc7dae7331c8daf88acda93d64644 192.168.1.130 27003 @ mymaster 192.168.1.130 7002
```



选举过程：看谁先发现宕机





RedisTemplate连接哨兵

`RedisTemplate` 是 Spring Framework 中用于与 Redis 数据库进行交互的一个类，它提供了高层次的 API 来执行 Redis 操作。`RedisTemplate` 使得开发者可以方便地进行 Redis 数据存取、序列化和其他操作。

