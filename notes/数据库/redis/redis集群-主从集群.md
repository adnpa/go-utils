https://redis.io/docs/latest/operate/oss_and_stack/management/replication/



## 拓扑结构

| IP            | PORT | 角色   |
| ------------- | ---- | ------ |
| 192.168.1.130 | 7001 | master |
| 192.168.1.130 | 7002 | slave  |
| 192.168.1.130 | 7003 | slave  |



## 准备配置

```
# 守护进行模式启动
daemonize no

# 集群相关配置
# 是否以集群模式启动
cluster-enabled yes

# 集群节点回应最长时间，超过该时间被认为下线
cluster-node-timeout 15000

# 生成的集群节点配置文件名，文件名需要修改
cluster-config-file nodes_6379.conf
```



7001

```
bind 192.168.1.130
replica-announce-ip 192.168.1.130
port 7001
dir /home/hz/testredis/7001
save 3600 1 300 100 60 10000
appendonly no
```

7002

```
bind 192.168.1.130
replica-announce-ip 192.168.1.130
slaveof 192.168.1.130 7001
port 7002
dir /home/hz/testredis/7002
save 3600 1 300 100 60 10000
appendonly no
```

7003

```
bind 192.168.1.130
replica-announce-ip 192.168.1.130
slaveof 192.168.1.130 7001
port 7003
dir /home/hz/testredis/7003
save 3600 1 300 100 60 10000
appendonly no
```

相关linux命令

```sh
echo 7001 7002 7003 | xargs -t -n 1 cp /etc/redis/redis.conf

sed -i -e 's/6379/7001/g' -e 's/dir .\//dir \/tmp\/7001\//g' 7001/redis.conf
sed -i -e 's/6379/7002/g' -e 's/dir .\//dir \/tmp\/7002\//g' 7001/redis.conf
sed -i -e 's/6379/7003/g' -e 's/dir .\//dir \/tmp\/7003\//g' 7001/redis.conf
```

## 开启主从关系

配置中 slaveof 永久生效，命令行中使用重启失效

## 启动

```sh
redis-server redis.conf
```

## 测试

```sh
hz@hz:~$ redis-cli -h 192.168.1.130 -p 7001
192.168.1.130:7001> set num 123
OK
192.168.1.130:7001> exit


hz@hz:~$ redis-cli -h 192.168.1.130 -p 7002
192.168.1.130:7002> get num
"123"
192.168.1.130:7002> set abc 111
(error) READONLY You can't write against a read only replica.
```





## 简介

Redis 主从集群是一种常用的架构，用于提高数据的可用性和读取性能。在这种架构中，一个主节点（Master）负责处理写操作，而一个或多个从节点（Slave）则负责复制主节点的数据并处理读取请求。以下是关于 Redis 主从集群的详细介绍。

### 1. **基本概念**

- **主节点（Master）**: 负责接收和处理所有的写请求，并将数据变化同步到从节点。
- **从节点（Slave）**: 负责从主节点复制数据，处理读取请求，减轻主节点的负担。
- **数据复制**: 从节点定期与主节点进行同步，确保数据一致性。

### 2. **架构优点**

- **读写分离**: 通过将读操作分配给从节点，主节点能够专注于写操作，提高整体性能。
- **高可用性**: 如果主节点发生故障，可以手动或自动将从节点提升为新的主节点，从而保证系统的可用性。
- **负载均衡**: 将读取请求分散到多个从节点，可以有效分担主节点的压力。

### 3. **配置主从集群**

以下是配置 Redis 主从集群的基本步骤：

#### 1. 启动 Redis 实例

启动多个 Redis 实例，至少一个主节点和一个从节点。可以在不同的端口上启动，例如：

- 主节点：`127.0.0.1:6379`
- 从节点：`127.0.0.1:6380`

#### 2. 配置主节点

默认情况下，主节点的配置无需修改。但可以在 `redis.conf` 中设置一些参数，比如：

```plaintext
# 开启持久化
save 900 1
```

#### 3. 配置从节点

在从节点的 `redis.conf` 配置文件中，设置为从属于主节点：

```plaintext
# 设置从属主节点
slaveof 127.0.0.1 6379
```

#### 4. 启动从节点

启动从节点，使其开始从主节点复制数据。

### 4. **使用命令管理主从关系**

- **查看主从状态**:
  ```bash
  redis-cli -p 6379 info replication
  ```
  该命令将返回主节点和从节点的状态信息。

- **手动切换主从**:
  如果需要手动将某个节点设置为从节点，可以使用：
  ```bash
  redis-cli -p 6380 slaveof 127.0.0.1 6379
  ```

- **提升从节点**:
  如果主节点故障，可以将从节点提升为主节点：
  ```bash
  redis-cli -p 6380 replicaof no one
  ```

### 5. **数据一致性**

- **异步复制**: 默认情况下，Redis 的数据复制是异步的，这意味着从节点可能会有短暂的延迟，数据在主节点写入后，从节点不会立即更新。
- **持久化**: 配置持久化策略（如 RDB 或 AOF），确保数据在重启后能够恢复。

### 6. **监控与管理**

- **监控工具**: 使用 Redis 自带的监控命令（如 `INFO`）来监控主从节点的状态。
- **故障检测**: 定期检查主节点和从节点的健康状态，确保系统稳定运行。

### 总结

Redis 主从集群架构通过读写分离和数据复制，提高了系统的可用性和性能。配置主从关系相对简单，但在生产环境中需要注意数据一致性和故障处理策略。通过合理的监控和管理，可以有效利用 Redis 主从集群的优势。





## 为什么使用这种架构

读多写少，读压力大





## Docker搭建过程

https://medium.com/@ahmettuncertr/redis-cluster-using-docker-1c8458a93d4b

https://www.youtube.com/watch?v=MCQhCBF8KFk





### 创建redis-cluster网络

```
docker network create redis-cluster
```



```
docker run -e ALLOW_EMPTY_PASSWORD -e REDIS_NODES=redis1,.  --name redis1 --network redis-cluster -p 6379:6379 -d bitnami/redis-cluster:latest

docker run -e ALLOW_EMPTY_PASSWORD -e REDIS_NODES=redis1,...  --name redis2 --network redis-cluster -p 6380:6379 -d bitnami/redis-cluster:latest

...
```



```
docker exec -it redis1 redis-cli --cluster create redis1:6379 redis1:6380 redis1:6381 --cluster-replicas 1 --cluster-yes
```







### 1. **准备工作**

确保你已经安装了 Docker 和 Docker Compose。你可以在命令行中运行以下命令检查安装状态：

```bash
docker --version
docker-compose --version
```

### 2. **创建 Docker Compose 文件**

创建一个 `docker-compose.yml` 文件，配置 Redis 主从节点。

```yaml
version: '3.8'

services:
  redis-master:
    image: redis:latest
    ports:
      - "6379:6379"
    volumes:
      - master-data:/data
    command: ["redis-server", "--appendonly", "yes"]

  redis-slave:
    image: redis:latest
    ports:
      - "6380:6379"
    volumes:
      - slave-data:/data
    command: ["redis-server", "--slaveof", "redis-master", "6379", "--appendonly", "yes"]

volumes:
  master-data:
  slave-data:
```

### 3. **启动 Redis 集群**

在包含 `docker-compose.yml` 文件的目录下，运行以下命令启动 Redis 主从集群：

```bash
docker-compose up -d
```

这将启动两个 Redis 实例：一个主节点和一个从节点。

### 4. **验证集群状态**

你可以使用 `redis-cli` 连接到主节点和从节点，验证主从关系是否配置正确。

#### 连接到主节点

```bash
docker exec -it <container_id_of_master> redis-cli
```

在 Redis CLI 中，运行以下命令查看主节点信息：

```bash
INFO replication
```

你应该看到类似以下的输出：

```
# Replication
role:master
connected_slaves:1
slave0:ip=172.18.0.4,port=6379,state=online,offset=12345,lag=0
```

#### 连接到从节点

```bash
docker exec -it <container_id_of_slave> redis-cli
```

在从节点的 Redis CLI 中，运行以下命令：

```bash
INFO replication
```

你应该看到类似以下的输出：

```
# Replication
role:slave
master_host:redis-master
master_port:6379
master_link_status:up
```

### 5. **测试主从同步**

在主节点中插入一些数据：

```bash
docker exec -it <container_id_of_master> redis-cli set testkey "Hello, Redis!"
```

然后在从节点中检查数据：

```bash
docker exec -it <container_id_of_slave> redis-cli get testkey
```

应该能够看到从节点返回的值为 `Hello, Redis!`。

### 6. **停止和清理**

要停止并清理 Docker Compose 创建的容器，可以运行以下命令：

```bash
docker-compose down
```

### 总结

通过上述步骤，你可以轻松使用 Docker 搭建一个 Redis 主从集群。这个集群支持数据的读写分离，并且可以通过 Docker 容器快速部署和管理。你可以根据需要扩展更多从节点，只需在 `docker-compose.yml` 中添加新的服务配置即可。

```dockerfile
```





