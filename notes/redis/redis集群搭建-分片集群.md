## redis实例

创建文件夹-创建配置文件-创建redis实例

```
mkdir 7000 7001 7002
touch 7000/redis.conf
...
cd 7000
redis-server ./redis.conf
```





## 创建 Redis 集群 

手动配置

```
redis-cli --cluster create 127.0.0.1:7000 127.0.0.1:7001 \
127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 \
--cluster-replicas 1
```

create-cluster脚本，这是redis提供的脚本，一般用于学习或演示，会启动一个有 3 个主服务器和 3 个副本的 6 节点集群

```
create-cluster start
create-cluster create
create-cluster stop
```



## 与集群交互

可分别和不同节点通信 

```
$ redis-cli -c -p 7000
redis 127.0.0.1:7000> set foo bar
-> Redirected to slot [12182] located at 127.0.0.1:7002
OK
redis 127.0.0.1:7002> set hello world
-> Redirected to slot [866] located at 127.0.0.1:7000
OK
redis 127.0.0.1:7000> get foo
-> Redirected to slot [12182] located at 127.0.0.1:7002
"bar"
redis 127.0.0.1:7002> get hello
-> Redirected to slot [866] located at 127.0.0.1:7000
"world"
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





























