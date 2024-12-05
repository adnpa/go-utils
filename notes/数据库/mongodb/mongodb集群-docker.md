

https://www.mongodb.com/resources/products/fundamentals/clusters

https://www.mongodb.com/resources/products/fundamentals/mongodb-cluster-setup

[配置文件](https://www.mongodb.com/docs/manual/reference/configuration-options/)



mongo组件

| [`mongod`](https://www.mongodb.com/zh-cn/docs/manual/reference/program/mongod/#mongodb-binary-bin.mongod) | 核心数据库流程                                               |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| [`mongos`](https://www.mongodb.com/zh-cn/docs/manual/reference/program/mongos/#mongodb-binary-bin.mongos) | [分片集群](https://www.mongodb.com/zh-cn/docs/manual/reference/glossary/#std-term-sharded-cluster)的控制器和查询路由器 |
| [`mongosh`](https://www.mongodb.com/zh-cn/docs/mongodb-shell/#mongodb-binary-bin.mongosh) | 交互式 shell                                                 |



https://hub.docker.com/_/mongo



https://www.mongodb.com/resources/products/compatibilities/deploying-a-mongodb-cluster-with-docker



1. Create a Docker network.
2. Start three instances of MongoDB.
3. Initiate the Replica Set.



## 副本集

### 1 创建docker网络

```
docker network create mongoCluster
```

### 2 启动节点

* -d 表示此容器应以分离模式运行（在后台）。 
* -p 表示端口映射。您机器上端口 27017 上的任何传入请求都将被重定向到容器中的端口 27017。 
* --name 表示容器的名称。这将成为此机器的主机名。 
* --network 表示要使用哪个 Docker 网络。同一网络中的所有容器都可以互相看到。 
* mongo:5 是 Docker 将使用的映像。此映像是 MongoDB 社区服务器版本 5（由 Docker 维护）。您也可以使用 MongoDB Enterprise 自定义映像。

```
docker run -d --rm -p 27017:27017 --name mongo1 --network mongoCluster mongo:5 mongod --replSet myReplicaSet --bind_ip localhost,mongo1
```

### 3 初始化副本集

```
docker exec -it mongo1 mongosh --eval "rs.initiate({
 _id: \"myReplicaSet\",
 members: [
   {_id: 0, host: \"mongo1\"},
   {_id: 1, host: \"mongo2\"},
   {_id: 2, host: \"mongo3\"}
 ]
})"
```

### 4 测试

```
docker exec -it mongo1 mongosh --eval "rs.status()"
```

一键部署

```yaml
# Use root/example as user/password credentials
version: '3.1'

services:

  mongo:
    image: mongo
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: example

  mongo-express:
    image: mongo-express
    restart: always
    ports:
      - 8081:8081
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: root
      ME_CONFIG_MONGODB_ADMINPASSWORD: example
      ME_CONFIG_MONGODB_URL: mongodb://root:example@mongo:27017/
      ME_CONFIG_BASICAUTH: false
```



```
docker run --name some-mongo -v /my/custom:/etc/mongo -d mongo --config /etc/mongo/mongod.conf

```

























