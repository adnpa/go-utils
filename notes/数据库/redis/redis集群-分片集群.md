https://redis.io/docs/latest/operate/oss_and_stack/management/scaling/



## 拓扑结构

| IP            | PORT | 角色   |
| ------------- | ---- | ------ |
| 192.168.1.130 | 7001 | master |
| 192.168.1.130 | 7002 | master |
| 192.168.1.130 | 7003 | master |
| 192.168.1.130 | 8001 | slave  |
| 192.168.1.130 | 8002 | slave  |
| 192.168.1.130 | 8003 | slave  |



## 准备配置

创建文件夹-创建配置文件-创建redis实例

```
port 7001

bind 0.0.0.0
replica-announce-ip 192.168.1.130

cluster-enabled yes
cluster-config-file /tmp/6379/nodes.conf
cluster-node-timeout 5000

dir /home/hz/testredis/7001

damonize yes
protected-mode no
databases 1
logfile /tmp/6379/run.log
```

说明：

1. **cluster-enabled <yes/no>**：如果设置为 `yes`，则启用 Redis 集群支持，实例将作为集群节点运行。否则，实例将以常规的单节点模式启动。

2. **cluster-config-file <filename>**：尽管名称为“配置文件”，但该文件不是用户可编辑的，而是 Redis 集群节点在每次状态更改时自动持久化集群配置的文件。该文件列出了集群中的其他节点、它们的状态、持久变量等。通常，在接收到消息后，这个文件会被重写并刷新到磁盘。

3. **cluster-node-timeout <milliseconds>**：Redis 集群节点可以在被视为故障之前的最大不可用时间。如果主节点在指定时间内无法访问，它将被其从节点进行故障转移。此参数还控制 Redis 集群中的其他重要内容，特别是每个节点在指定时间内无法访问大多数主节点时，将停止接受查询。

4. **cluster-slave-validity-factor <factor>**：如果设置为零，从节点将始终认为自己是有效的，因此将始终尝试故障转移主节点，而不考虑与主节点断开连接的时间。如果值为正，则最大断开时间将根据节点超时值乘以指定的因子进行计算。如果节点是从节点，并且与主节点的链接断开超过指定时间，则不会尝试故障转移。例如，如果节点超时设置为 5 秒，有效性因子设置为 10，则从节点如果与主节点断开超过 50 秒，将不会尝试故障转移。请注意，任何非零值可能导致在主节点故障后 Redis 集群不可用，除非有能够故障转移的从节点。在这种情况下，集群将在原主节点重新加入集群时恢复可用性。

5. **cluster-migration-barrier <count>**：主节点与其从节点保持连接的最小数量，以便另一个从节点可以迁移到不再由任何从节点覆盖的主节点。

6. **cluster-require-full-coverage <yes/no>**：如果设置为 `yes`（默认值），则在某些键空间没有任何节点覆盖时，集群将停止接受写入。如果设置为 `no`，即使只有对部分键的请求能够处理，集群仍将继续提供查询服务。

7. **cluster-allow-reads-when-down <yes/no>**：如果设置为 `no`（默认值），当集群被标记为故障时，Redis 集群中的节点将停止处理所有流量，无论是无法访问主节点的法定人数还是未满足完全覆盖。这可以防止从一个对集群变化不了解的节点读取潜在的不一致数据。设置为 `yes` 可以允许在故障状态下从节点读取，这对于希望优先考虑读取可用性的应用程序非常有用，同时仍然希望防止不一致的写入。它也可以用于只有一个或两个分片的 Redis 集群，因为它允许节点在主节点故障但无法自动故障转移时继续处理写入请求。

## 启动实例

```sh
redis-server 7001/redis.conf
redis-server 7002/redis.conf
redis-server 7003/redis.conf
redis-server 8001/redis.conf
redis-server 8002/redis.conf
redis-server 8003/redis.conf
```



## 创建 Redis 集群 

手动配置

```
redis-cli --cluster create 192.168.1.130:7001 192.168.1.130:7002 192.168.1.130:7003\
 192.168.1.130:8001 192.168.1.130:8002 192.168.1.130:8003 \
--cluster-replicas 1
```

create-cluster脚本，这是redis提供的脚本，一般用于学习或演示，会启动一个有 3 个主服务器和 3 个副本的 6 节点集群

```
create-cluster start
create-cluster create
create-cluster stop
```



## 测试

重定向由客户端实现

`-c` 选项用于启用集群模式的命令行界面（CLI），即 `redis-cli`。使用这个选项时，`redis-cli` 能够自动处理集群中的重定向（如 `MOVED` 和 `ASK` 错误）

```
hz@hz:~/test-cluster$ redis-cli -h 192.168.1.130 -p 7001
192.168.1.130:7001> set num 123
OK
192.168.1.130:7001> exit

hz@hz:~/test-cluster$ redis-cli -h 192.168.1.130 -p 8002
192.168.1.130:8002> get num
(error) MOVED 2765 192.168.1.130:7001

hz@hz:~/test-cluster$ redis-cli -c -h 192.168.1.130 -p 8001
192.168.1.130:8001> get num
-> Redirected to slot [2765] located at 192.168.1.130:7001
"123"
```

不使用-c选项发现，无论向哪一个slave发送get都不能成功



## 集群命令

cluster help查看

```sh
redis-cli -p 7001 cluster nodes
```



```
cluster nodes
```







## 使用 redis-rb-cluster 编写示例应用程序 

```
```



## 对集群进行重新分片 

重新分片（resharding）是指在集群中重新分配数据分片，以实现负载均衡或添加/移除节点

重新分片基本上意味着**将哈希槽从一组节点移动到另一组节点**。与集群创建一样，它是使用 redis-cli 实用程序完成的。

命令

```
redis-cli --cluster reshard <host>:<port> 
--cluster-from <node-id> 
--cluster-to <node-id> 
--cluster-slots <number of slots> 
--cluster-yes
```





```
redis-cli --cluster reshard 127.0.0.1:7000
```

选择移动多少槽

```
How many slots do you want to move (from 1 to 16384)?
```





## 更有趣的示例应用程序 



## 测试故障转移(failover)

让单个进程崩溃

```
redis-cli -p 7002 debug segfault
```

## 手动故障转移 



## 添加新节点 

### 1添加实例

和第一部分一样，创建文件夹，配置文件，启动

### 2主节点，重新分片

```
redis-cli --cluster add-node 127.0.0.1:7006 127.0.0.1:7000
```

### 3从节点，复制

不指定，会将新节点添加为副本较少的主节点中的随机主节点的副本。

```
redis-cli --cluster add-node 127.0.0.1:7006 127.0.0.1:7000 --cluster-slave
```

指定

```
redis-cli --cluster add-node 127.0.0.1:7006 127.0.0.1:7000 --cluster-slave --cluster-master-id 3c3a0c74aae0b56170ccb03a76b60cfe7dc1912e
```

简洁方式

```
cluster replicate 3c3a0c74aae0b56170ccb03a76b60cfe7dc1912e
```



## 删除节点 

**移除主节点必须把其中的数据转移**

```
redis-cli --cluster del-node 127.0.0.1:7000 `<node-id>`
```



```
redis-cli --cluster call 127.0.0.1:7000 cluster forget `<node-id>`
```





## 副本迁移 

在 Redis 集群中，你可以随时使用以下命令重新配置副本以从其他主节点复制数据

```
CLUSTER REPLICATE <master-node-id>
```





## 升级 Redis 集群中的节点 



## 迁移到 Redis 集群





























