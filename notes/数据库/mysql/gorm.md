https://github.com/go-gorm/gorm

https://gorm.io/

https://gorm.io/gen/index.html





安装

```

go get gorm.io/gen
```



# 特性

* 全功能 ORM
* 关联 (Has One，Has Many，Belongs To，Many To Many，多态，单表继承)
* Create，Save，Update，Delete，Find 中钩子方法
* 支持 `Preload`、`Joins` 的预加载
* 事务，嵌套事务，Save Point，Rollback To Saved Point
* Context、预编译模式、DryRun 模式
* 批量插入，FindInBatches，Find/Create with Map，使用 SQL 表达式、Context Valuer 进行 CRUD
* SQL 构建器，Upsert，数据库锁，Optimizer/Index/Comment Hint，命名参数，子查询
* 复合主键，索引，约束
* Auto Migration
* 自定义 Logger
* 灵活的可扩展插件 API：Database Resolver（多数据库，读写分离）、Prometheus…
* 每个特性都经过了测试的重重考验
* 开发者友好



# gen工具



ORM（对象关系映射）是一种编程技术，用于将对象模型与关系数据库之间进行映射。它的目标是通过将对象和数据库表之间的映射关系定义在代码中，使得开发者能够以面向对象的方式操作数据库，而无需直接编写SQL语句。ORM工具提供了一组API和工具，用于执行数据库操作，包括数据检索、插入、更新和删除。

在ORM中，DAO（数据访问对象）是一种设计模式，用于封装对数据库的访问。DAO提供了一组抽象接口，用于定义对数据库的操作，包括数据的创建、读取、更新和删除。它的主要目的是将数据持久化层与业务逻辑层进行解耦，使得业务逻辑层能够独立于底层数据库的细节进行开发和测试。

通过使用DAO，开发者可以将对数据库的操作封装在具体的DAO实现类中，使得业务逻辑层只需通过调用DAO接口来进行数据库操作，而不需要直接与数据库交互。这样可以提高代码的可维护性和可测试性，同时也降低了代码的耦合度。ORM工具通常与DAO模式结合使用，以提供更便捷和高效的数据库访问方式。





```go
package main

import "gorm.io/gen"

func main() {
  g := gen.NewGenerator(gen.Config{
    OutPath: "../query",
    Mode: gen.WithoutContext|gen.WithDefaultQuery|gen.WithQueryInterface, // generate mode
  })

  // gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
  g.UseDB(gormdb) // reuse your gorm db

  // Generate basic type-safe DAO API for struct `model.User` following conventions
  
  g.ApplyBasic(
  // Generate struct `User` based on table `users`
  g.GenerateModel("users"),
  
  // Generate struct `Employee` based on table `users`
 g.GenerateModelAs("users", "Employee"),


// Generate struct `User` based on table `users` and generating options
g.GenerateModel("users", gen.FieldIgnore("address"), gen.FieldType("id", "int64")),

// Generate struct `Customer` based on table `customer` and generating options
// customer table may have a tags column, it can be JSON type, gorm/gen tool can generate for your JSON data type
g.GenerateModel("customer", gen.FieldType("tags", "datatypes.JSON")),

  )
g.ApplyBasic(
// Generate structs from all tables of current database
g.GenerateAllTable()...,
)
  // Generate the code
  g.Execute()
}
```



生成文件 model 和 query

model

```go
type User struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID    int64          `gorm:"column:user_id;not null" json:"user_id"`
	Username  string         `gorm:"column:username;not null" json:"username"`
	Password  string         `gorm:"column:password;not null" json:"password"`
	Email     string         `gorm:"column:email" json:"email"`
	Gender    int16          `gorm:"column:gender" json:"gender"`
	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}

```

query 提供子查询方法实现

```go
type IUserDo interface {
	gen.SubQuery
	Debug() IUserDo
	WithContext(ctx context.Context) IUserDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserDo
	WriteDB() IUserDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserDo
    //...
}
```





## 使用生成的文件操作数据库

https://gorm.io/zh_CN/gen/query.html

```go

```

























