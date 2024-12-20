https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/



## 两种持久化方式

* RDB
  * 备份方式：周期性写入磁盘，生成一个 .rdb 文件
  * 适用场景：冷备份，在服务器意外宕机时快速恢复数据
  * 优点：
    * rdb文件是非常紧凑的单文件时间点表示，非常适合备份
    * 更适合大数据重启
    * 支持副本故障转移后的部分重新同步
  * 缺点：
* AOF
  * 备份方式：实时记录执行的所有写操作
  * 适用场景：热备份，在服务器故障时尽量避免数据丢失
  * 优点：
  * 缺点：

另外，redis支持同时开启 RDB 和 AOF 两种持久化方式。在这种情况下,Redis 会优先使用 AOF 文件进行数据恢复



## 使用RDB进行快照

使用命令

* save 主进程执行RDB，阻塞所有操作
* bgsave 开启子进程执行RDB

原理

1. fork进程
2. 子进程开始将数据写道临时的rdb文件
3. 完成后替换旧的

相关配置

```
# 900秒内至少1个key修改则执行快照
save 900 1

# 是否压缩 建议不开启 压缩消耗cpu 而磁盘不值钱
rdbcompression yes 
# RDB文件名称
dbfilename dump.rdb
# 文件保存目录
dir ./
```



## 使用AOF写入日志

相关配置

```
# 是否开启
appendonly yes
# AOF文件名
appendfilename "appendonly.aof"
# 记录频率（写回策略）
appendfsync always
appendfsync everysec
appendfsync no

# 重写
# 上次文件增加超过多少重写
auto-aof-rewrite-percentage 100
# AOF体积最小多少触发重写
auto-aof-rewrite-min-size 64mb
```

写回策略



重写日志

可以通过bgrewriteaof命令重写





重写日志原理

* `>=7.0`
  1. fork进程
  2. 子进程写入 base AOF文件
  3. 父进程写入另一个新的增量AOF文件
  4. 子进程完成重写，父进程收到信号，合并后持久化
* `<7.0`
  1. fork进程
  2. 子进程写入 base AOF文件
  3. 父进程写入到内存中的AOF缓冲区
  4. 子进程完成重写，父进程收到信号，将内存中的命令插入到子进程新建文件









![image-20241201223541557](./test.assets/image-20241201223541557.png)



