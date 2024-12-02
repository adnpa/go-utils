https://www.mongodb.com/resources/products/capabilities/replication

https://www.mongodb.com/docs/manual/tutorial/deploy-replica-set/

## 副本集replica sets

MongoDB 副本集是一种数据冗余和高可用性解决方案，通过多个 MongoDB 实例（节点）协同工作，以确保数据的可靠性和可用性。以下是关于 MongoDB 副本集的详细信息，包括配置步骤和工作原理。

### 副本集的组成

1. **主节点（Primary）**：负责处理所有写入操作和大部分读取操作。副本集中只能有一个主节点。
2. **从节点（Secondary）**：负责从主节点复制数据，提供数据冗余和备份。可以配置为支持读取操作。
3. **仲裁节点（Arbiter）**：不存储数据，仅参与选举过程，以维持副本集的奇数节点数，避免分裂脑（split-brain）问题。

### 工作原理

- **数据复制**：主节点将操作记录在 oplog（操作日志）中，从节点定期从主节点的 oplog 中复制数据。
- **故障转移**：如果主节点不可用，其他节点会通过选举机制选出新的主节点，确保系统的高可用性。
- **读写分离**：可以将读取操作分散到从节点，减轻主节点的负担。



## 拓扑结构

| Hostname             | Ip            | Port  | Type      |
| -------------------- | ------------- | ----- | --------- |
| mongodb0.example.net | 192.168.1.130 | 27017 | Primary   |
| mongodb1.example.net | 192.168.1.130 | 27018 | Secondary |
| mongodb2.example.net | 192.168.1.130 | 27019 | Arbiter   |



Hostnames：从 MongoDB 5.0 开始，仅配置了 IP 地址的节点将无法通过启动验证并且不会启动，必须配置hostname

## 步骤

### 1 创建主节点、副本节点和仲裁节点

```yaml
replication:
   # 副本集名称
   replSetName: "rs0"
net:
   bindIp: localhost,192.168.1.130
   port: 27017
```

### 2 启动

```
mongod --config mongod.conf
```

### 3 初始化配置副本集和主节点

```
mongosh --host=192.168.1.130 --port=27017
test> rs.initiate()
```

等待一会后按回车变为主节点

### 4 添加副本节点和仲裁节点

```
rs.add("192.168.1.130:27018")
rs.addArb("192.168.1.130:27019")
```

### 5 查看状态和测试

```
rs.status()
```





错误

```
MongoServerError[NewReplicaSetConfigurationIncompatible]: Reconfig attempted to install a config that would change the implicit default write concern. Use the setDefaultRWConcern command to set a cluster-wide write concern and try the reconfig again.


MongoServerError[Unauthorized]: setDefaultRWConcern may only be run against the admin database.
```

解决

```
use admin
db.runCommand({
  setDefaultRWConcern: 1,
  defaultWriteConcern: {
    w: "majority",  // 设置为你希望的写关注级别
    wtimeout: 5000  // 可选，设置写超时
  }
});
```





关闭节点

```
db.adminCommand({ shutdown: 1 });
```

