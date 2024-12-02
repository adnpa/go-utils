https://www.mongodb.com/docs/manual/administration/deploy-manage-self-managed-sharded-clusters/	

## 分片集群sharded cluster



MongoDB 分片集群是一种将数据分布在多个服务器上的架构，旨在提高性能和可扩展性。以下是一些关于 MongoDB 分片集群的关键点：

### 1. 分片概念
- **分片**：将数据集分割成更小的部分（称为分片），每个分片可以存储在不同的服务器上。
- **数据划分**：通过使用 shard key（分片键）来决定数据的存储位置。分片键是文档中的一个字段，MongoDB 根据该字段的值将数据分配到不同的分片。

### 2. 组件
- **配置服务器**：存储集群的元数据和路由信息。配置服务器通常以复制集的形式部署。
- **路由服务（mongos）**：客户端与集群之间的中介，负责将请求路由到相应的分片。
- **数据分片**：实际存储数据的服务器，可以是单个实例，也可以是复制集。

### 3. 优势
- **水平扩展**：可以通过增加更多的分片轻松扩展数据库容量和性能。
- **负载均衡**：数据分布在不同的分片上，减少单台服务器的负载。
- **高可用性**：通过复制集的方式，提供容错能力。

### 4. 使用场景
适合于需要处理大量数据并且需要高并发的应用场景，如社交网络、电商平台和在线游戏等。

### 5. 注意事项
- **选择合适的分片键**：不当的分片键可能导致数据不均匀分布（热点问题），影响性能。
- **监控和维护**：分片集群的复杂性增加了运维的难度，需要定期监控集群状态。

通过合理配置和优化，MongoDB 分片集群可以为大规模应用提供强大的支持。

## 组件

* 分片，数据集群
* mongos，控制器和路由器，在客户端和集群间提供接口
* config servers，配置服务器，存储集群的元数据和配置，3.4开始必须部署为副本集



## 系统host

`/etc/hosts`

```
192.168.1.130 cfg1.example.net
192.168.1.130 cfg1.example.net
192.168.1.130 cfg1.example.net

192.168.1.130 s1-mongo1.example.net
192.168.1.130 s1-mongo2.example.net
192.168.1.130 s1-mongo3.example.net
```



## 创建配置服务器副本集

| Member   | Hostname           | Port  |
| -------- | ------------------ | ----- |
| Member 0 | `cfg1.example.net` | 27001 |
| Member 1 | `cfg2.example.net` | 27002 |
| Member 2 | `cfg3.example.net` | 27003 |

配置

```
sharding:
  clusterRole: configsvr
replication:
  replSetName: "confrs"
net:
  bindIp: localhost,192.168.1.130
```

启动-初始化

```
mongod --config 27001/mongod.conf
mongod --config 27002/mongod.conf
mongod --config 27003/mongod.conf

mongosh --host 192.168.1.130 --port 27001
rs.initiate(
  {
    _id: "confrs",
    configsvr: true,
    members: [
      { _id : 0, host : "cfg1.example.net:27001" },
      { _id : 1, host : "cfg2.example.net:27002" },
      { _id : 2, host : "cfg3.example.net:27003" }
    ]
  }
)
```



## 创建分片副本集

| Hostname              | Ip            | Port  | Type      |
| --------------------- | ------------- | ----- | --------- |
| s1-mongo1.example.net | 192.168.1.130 | 27004 | Primary   |
| s1-mongo2.example.net | 192.168.1.130 | 27005 | Secondary |
| s1-mongo3.example.net | 192.168.1.130 | 27006 | Arbiter   |
| s2-mongo1.example.net | 192.168.1.130 | 27007 | Primary   |
| s2-mongo2.example.net | 192.168.1.130 | 27008 | Secondary |
| s2-mongo3.example.net | 192.168.1.130 | 27009 | Arbiter   |

配置

```
sharding:
    clusterRole: shardsvr
replication:
    replSetName: "shardedrs"
net:
    bindIp: localhost,192.168.1.130
```

启动-初始化

```
mongod --config 27004/mongod.conf
mongod --config 27005/mongod.conf
mongod --config 27006/mongod.conf

mongosh --host 192.168.1.130 --port 27004
rs.initiate(
  {
    _id : "shardedrs",
    members: [
      { _id : 0, host : "s1-mongo1.example.net:27004" },
      { _id : 1, host : "s1-mongo2.example.net:27005" },
      { _id : 2, host : "s1-mongo3.example.net:27006" }
    ]
  }
)
```





## 启动mongos

配置

```
sharding:
  configDB: confrs/cfg1.example.net:27001,cfg2.example.net:27002,cfg3.example.net:27003
net:
  bindIp: localhost,192.168.1.130
  port: 27007
```

启动-添加分片

```
mongos --config mongos.conf

mongosh --host 192.168.1.130 --port 27007
sh.addShard( "shardedrs/s1-mongo1.example.net:27004,s1-mongo2.example.net:27005,s1-mongo3.example.net:27006")
```





## 集合分片



## 键分片



