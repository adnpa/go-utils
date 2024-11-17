https://pkg.go.dev/slices@go1.22.3

https://pkg.go.dev/database/sql/driver@go1.22.3

https://pkg.go.dev/database/sql

https://go.dev/doc/tutorial/database-access



# 基础





# 数据库优化

https://www.zhihu.com/question/36431635/answer/3107541973?utm_psn=1840002132541046784

## 硬件配置

* 磁盘扩容
* 机械硬盘换成固态
* Cpu核数
* 内存扩容，让Buffer Pool能吃进更多数据

## 参数配置

https://www.cnblogs.com/duanxz/p/3875760.html

* 保证从内存读取
  innodb_buffer_pool_size = 750M

* 数据预热
  将磁盘中的数据提前放入 BufferPool 内存缓冲池内

* 降低磁盘的写入次数

* （1）增大 [redo log](https://zhida.zhihu.com/search?content_id=595092286&content_type=Answer&match_order=1&q=redo+log&zhida_source=entity)，减少落盘次数：

  redo log 是重做日志，用于保证数据的一致，减少落盘相当于减少了系统 IO 操作。

  [innodb_log_file_size](https://zhida.zhihu.com/search?content_id=595092286&content_type=Answer&match_order=1&q=innodb_log_file_size&zhida_source=entity) 设置为 0.25 * innodb_buffer_pool_size

  （2）通用查询日志、慢查询日志可以不开 ，binlog 可开启。

  通用查询和慢查询日志也是要落盘的，可以根据实际情况开启，如果不需要使用的话就可以关掉。binlog 用于恢复和主从复制，这个可以开启。

  查看相关参数的命令：

  ```text
  # 慢查询日志
  show variables like 'slow_query_log%'
  # 通用查询日志
  show variables like '%general%';
  # 错误日志
  show variables like '%log_error%'
  # 二进制日志
  show variables like '%binlog%';
  ```

  （3）写 redo log 策略 [innodb_flush_log_at_trx_commit](https://zhida.zhihu.com/search?content_id=595092286&content_type=Answer&match_order=1&q=innodb_flush_log_at_trx_commit&zhida_source=entity) 设置为 0 或 2

  对于不需要强一致性的业务，可以设置为 0 或 2。

  - 0：每隔 1 秒写日志文件和刷盘操作（写日志文件 LogBuffer --> OS cache，刷盘 OS cache --> 磁盘文件），最多丢失 1 秒数据
  - 1：事务提交，立刻写日志文件和刷盘，数据不丢失，但是会频繁 IO 操作
  - 2：事务提交，立刻写日志文件，每隔 1 秒钟进行刷盘操作

* 系统调优参数

  * **back_log**
    back_log值可以指出在MySQL暂时停止回答新请求之前的短时间内多少个请求可以被存在堆栈中。也就是说，如果MySQL的连接数据达到max_connections时，新来的请求将会被存在堆栈中，以等待某一连接释放资源，该堆栈的数量即back_log，如果等待连接的数量超过back_log，将不被授予连接资源。可以从默认的50升至500。

  * **wait_timeout**
    数据库连接闲置时间，闲置连接会占用内存资源。可以从默认的8小时减到半小时。

  * **max_user_connection**
    最大连接数，默认为0无上限，最好设一个合理上限。

  * **thread_concurrency**
    并发线程数，设为CPU核数的两倍。

  * **skip_name_resolve**
    禁止对外部连接进行DNS解析，消除DNS解析时间，但需要所有远程主机用IP访问。

  * **[key_buffer_size](https://zhida.zhihu.com/search?content_id=595092286&content_type=Answer&match_order=1&q=key_buffer_size&zhida_source=entity)**
    索引块的缓存大小，增加会提升索引处理速度，对MyISAM表性能影响最大。对于内存4G左右，可设为256M或384M，通过查询show status like 'key_read%'，保证key_reads / key_read_requests在0.1%以下最好。

  * **innodb_buffer_pool_size**
    缓存数据块和索引块，对InnoDB表性能影响最大。通过查询show status like 'Innodb_buffer_pool_read%'，保证 (Innodb_buffer_pool_read_requests – Innodb_buffer_pool_reads) / Innodb_buffer_pool_read_requests越高越好。

  * **innodb_additional_mem_pool_size**
    InnoDB存储引擎用来存放数据字典信息以及一些内部数据结构的内存空间大小，当数据库对象非常多的时候，适当调整该参数的大小以确保所有数据都能存放在内存中提高访问效率，当过小的时候，MySQL会记录Warning信息到数据库的错误日志中，这时就需要该调整这个参数大小。

  * **innodb_log_buffer_size**
    InnoDB存储引擎的事务日志所使用的缓冲区，一般来说不建议超过32MB。

  * **[query_cache_size](https://zhida.zhihu.com/search?content_id=595092286&content_type=Answer&match_order=1&q=query_cache_size&zhida_source=entity)**
    缓存MySQL中的ResultSet，也就是一条SQL语句执行的结果集，所以仅仅只能针对select语句。当某个表的数据有任何变化，都会导致所有引用了该表的select语句在Query Cache中的缓存数据失效。所以，当我们数据变化非常频繁的情况下，使用Query Cache可能得不偿失。根据命中率(Qcache_hits/(Qcache_hits+Qcache_inserts)*100))进行调整，一般不建议太大，256MB可能已经差不多了，大型的配置型静态数据可适当调大。可以通过命令show status like 'Qcache_%'查看目前系统Query catch使用大小。

  * **read_buffer_size**
    MySQL读入缓冲区大小。对表进行顺序扫描的请求将分配一个读入缓冲区，MySQL会为它分配一段内存缓冲区。如果对表的顺序扫描请求非常频繁，可以通过增加该变量值以及内存缓冲区大小来提高其性能。

  * **[sort_buffer_size](https://zhida.zhihu.com/search?content_id=595092286&content_type=Answer&match_order=1&q=sort_buffer_size&zhida_source=entity)**
    MySQL执行排序使用的缓冲大小。如果想要增加ORDER BY的速度，首先看是否可以让MySQL使用索引而不是额外的排序阶段。如果不能，可以尝试增加sort_buffer_size变量的大小。

  * **read_rnd_buffer_size**
    MySQL的随机读缓冲区大小。当按任意顺序读取行时(例如按照排序顺序)，将分配一个随机读缓存区。进行排序查询时，MySQL会首先扫描一遍该缓冲，以避免磁盘搜索，提高查询速度，如果需要排序大量数据，可适当调高该值。但MySQL会为每个客户连接发放该缓冲空间，所以应尽量适当设置该值，以避免内存开销过大。

  * **record_buffer**
    每个进行一个顺序扫描的线程为其扫描的每张表分配这个大小的一个缓冲区。如果你做很多顺序扫描，可能想要增加该值。

  * **thread_cache_size**
    保存当前没有与连接关联但是准备为后面新的连接服务的线程，可以快速响应连接的线程请求而无需创建新的。

  * **table_cache**
    类似于thread_cache _size，但用来缓存表文件，对InnoDB效果不大，主要用于MyISAM。


## 表结构设计

### 设计聚合表

设计聚合表，一般针对于统计分析功能，或者实时性不高的需求（报表统计，数据分析等系统），这是一种空间 + 时延性换时间的思想。

### 设计冗余字段

为减少关联查询，创建合理的冗余字段（创建冗余字段还需要注意数据一致性问题），当然，如果冗余字段过多，对系统复杂度和插入性能会有影响。

### 分表

分表分为垂直拆分和水平拆分两种。

* 垂直拆分，适用于字段太多的大表，比如：一个表有100多个字段，那么可以把表中经常不被使用的字段或者存储数据比较多的字段拆出来。
* 水平拆分，比如：一个表有5千万数据，那按照一定策略拆分成十个表，每个表有500万数据。这种方式，除了可以解决查询性能问题，也可以解决数据写操作的热点征用问题。

### 字段设计

数据库中的表越小，在它上面执行的查询也就会越快。因此，在创建表的时候，为了获得更好的性能，我们可以将表中字段的宽度设得尽可能小。

- 使用可以存下数据最小的数据类型，合适即可
- 尽量使用TINYINT、SMALLINT、MEDIUM_INT作为整数类型而非INT，如果非负则加上UNSIGNED；
- VARCHAR的长度只分配真正需要的空间；
- 对于某些文本字段，比如"省份"或者"性别"，使用枚举或整数代替字符串类型；在MySQL中， ENUM类型被当作数值型数据来处理，而数值型数据被处理起来的速度要比文本类型快得多
- 尽量使用TIMESTAMP而非DATETIME；
- 单表不要有太多字段，建议在20以内；
- 尽可能使用 not null 定义字段，null 占用4字节空间，这样在将来执行查询的时候，数据库不用去比较NULL值。
- 用整型来存IP。
- 尽量少用 text 类型，非用不可时最好考虑拆表。



## SQL语句及索引

### 使用 EXPLAIN 分析 SQL

这里对explain的结果进行简单说明：

- select_type：查询类型

- - SIMPLE 简单查询
  - PRIMARY 最外层查询
  - UNION union后续查询
  - SUBQUERY 子查询

- type：查询数据时采用的方式

- - ALL 全表**（性能最差）**
  - index 基于索引的全表
  - range 范围 （< > in）
  - ref 非唯一索引单值查询
  - const 使用主键或者唯一索引等值查询

- possible_keys：可能用到的索引

- key：真正用到的索引

- rows：预估扫描多少行记录

- key_len：使用了索引的字节数

- Extra：额外信息

- - Using where 索引回表
  - Using index 索引直接满足条件
  - Using filesort 需要排序
  - Using temprorary 使用到临时表

对于以上的几个列，我们重点关注的是type，最直观的反映出SQL的性能。



### SQL语句尽可能简单

一条sql只能在一个cpu运算；大语句拆小语句，减少锁时间；一条大sql可以堵死整个库。

### 对于连续数值，使用 BETWEEN 不用 IN

SELECT id FROM t WHERE num BETWEEN 1 AND 5；

### SQL 语句中 IN 包含的值不应过多

MySQL对于IN做了相应的优化，即将IN中的常量全部存储在一个数组里面，而且这个数组是排好序的。如果数值较多，需要在内存进行排序操作，产生的消耗也是比较大的。

### SELECT 语句必须指明字段名称

SELECT * 增加很多不必要的消耗（CPU、IO、内存、网络带宽）；减少了使用覆盖索引的可能性。

### 当只需要一条数据的时候，使用 limit 1

limit 相当于截断查询。

例如：对于select * from user limit 1; 虽然进行了全表扫描，但是limit截断了全表扫描，从0开始取了1条数据。

### 排序字段加索引

排序的字段建立索引在排序的时候也会用到

### 如果限制条件中其他字段没有索引，尽量少用or

### 尽量用 union all 代替 union

union和union all的差别就在于union会对数据做一个distinct的动作，而这个distanct动作的速度则取决于现有数据的数量，数量越大则时间也越慢。而对于几个数据集，要确保数据集之间的数据互相不重复，基本是O(n)的算法复杂度。

### 区分 in 和 exists、not in 和 not exists

如果是exists，那么以外层表为驱动表，先被访问，如果是IN，那么先执行子查询。所以IN适合于外表大而内表小的情况；EXISTS适合于外表小而内表大的情况。

### 使用合理的分页方式以提高分页的效率

limit m n，其中的m偏移量尽量小。m越大查询越慢。

### 避免使用 % 前缀模糊查询

例如：like '%name'或者like '%name%'，这种查询会导致索引失效而进行全表扫描。但是可以使用like 'name%'，这种会使用到索引。

### 避免在 where 子句中对字段进行表达式操作

这种不会使用到索引：

```text
select user_id,user_project from user_base where age*2=36;
```

可以改为：

```text
select user_id,user_project from user_base where age=36/2;
```

任何对列的操作都将导致表扫描，它包括数据库函数、计算表达式等等，查询时要尽可能将操作移至等号右边。

### 避免隐式类型转换

where 子句中出现的 column 字段要和数据库中的字段类型对应

### 必要时可以使用 force index 来强制查询走某个索引

有的时候 MySQL 优化器采取它认为合适的索引来检索 SQL 语句，但是可能它所采用的索引并不是我们想要的。这时就可以采用 forceindex 来强制优化器使用我们制定的索引。

### 使用联合索引时注意范围查询

对于联合索引来说，如果存在范围查询，比如between、>、<等条件时，会造成后面的索引字段失效。

### 某些情况下，可以使用连接代替子查询

因为使用 join，MySQL 不会在内存中创建临时表。

### 使用JOIN的优化

使用小表驱动大表，例如使用inner join时，优化器会选择小表作为驱动表

### 小表驱动大表，即小的数据集驱动大的数据集

如：以 A，B 两表为例，两表通过 id 字段进行关联。

```text
#当 B 表的数据集小于 A 表时，用 in 优化 exist；使用 in ，两表执行顺序是先查 B 表，再查 A 表
select * from A where id in (select id from B)

#当 A 表的数据集小于 B 表时，用 exist 优化 in；使用 exists，两表执行顺序是先查 A 表，再查 B 表
select * from A where exists (select 1 from B where B.id = A.id)
```


上面都是一些常规的优化方法，我们还可以使用：主从和分库。











# 相关库

* https://github.com/go-mysql-org/go-mysql
* https://github.com/go-mysql-org/go-mysql-elasticsearch































