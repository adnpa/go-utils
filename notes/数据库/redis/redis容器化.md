https://www.docker.com/blog/how-to-use-the-redis-docker-official-image/

https://redis.io/learn/operate/orchestration/docker





目前，Redis Cluster 不支持 NATted 环境以及重新映射 IP 地址或 TCP 端口的一般环境。 

Docker 使用一种称为端口映射的技术：在 Docker 容器内运行的程序可能会暴露与程序认为正在使用的端口不同的端口。这对于在同一台服务器上同时使用相同端口运行多个容器非常有用。 

要使 Docker 与 Redis Cluster 兼容，您需要使用 Docker 的**主机网络模式**。有关更多信息，请参阅 Docker 文档中的 --net=host 选项。

