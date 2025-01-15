

https://go.dev/doc/tutorial/web-service-gin

https://pkg.go.dev/github.com/gin-gonic/gin



## 配置路由

```go

func main() {
	r := gin.Default()

    //跨域配置
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:9000"},                   // 允许的来源
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},                  // 允许的方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 允许暴露的响应头
		AllowCredentials: true,                                                // 是否允许使用凭证
		MaxAge:           12 * 3600,                                           // 预检请求的缓存时间
	}

	r.Use(cors.New(corsConfig))
    
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	apiGroup := r.Group("/api")
	apiGroup.POST("/category/:typ", api.Category)
	apiGroup.POST("/book/:id", api.BookInfo)
	apiGroup.POST("/top", api.Top)
	apiGroup.POST("/bookshelf", api.Bookshelf)

	r.Run("127.0.0.1:10000") // 监听并在 0.0.0.0:8080 上启动服务
}
```



## 获取参数

### 获取路径参数

路径参数是 URL 路径的一部分，通常用于标识特定的资源。例如，在 URL `/api/users/123` 中，`123` 是一个路径参数，通常表示用户 ID。

通常用于**GET请求**

```go
r.GET("/api/users/:id", func(c *gin.Context) {
    // 获取路径参数
    userID := c.Param("id")
    c.JSON(http.StatusOK, gin.H{"user_id": userID})
})
```



### 获取查询参数

查询参数是 URL 中以问号 (`?`) 开始的部分，通常用于过滤或排序结果。例如，在 URL `/api/users?age=30&name=Alice` 中，`age` 和 `name` 是查询参数。

通常用于**GET请求**

```go
r.GET("/api/users", func(c *gin.Context) {
    // 获取查询参数":""},
":""},
":""} 
    age := c.Query("age")
    name := c.Query("name")
    c.JSON(http.StatusOK, gin.H{
        "age":  age,
        "name": name,
    })
})
```



### 获取header参数

```go
r.GET("/api/example", func(c *gin.Context) {
    // 获取特定的头部参数
    customHeader := c.GetHeader("X-Custom-Header")
    authHeader := c.GetHeader("Authorization")

    // 返回头部参数
    c.JSON(http.StatusOK, gin.H{
        "custom_header": customHeader,
        "auth_header":   authHeader,
    })
})
```



### 获取body参数

Get请求参数以键值对在url中传递，而Post请求则包含在请求体中。HTML表单中的enctype字段决定Post以哪种形式传递。

* application/x-www-formurlencoded 默认值，编码为长查询字符串，以&和=分割不同键值对和键、值，例:`k1=v1&k2=v2`，传送**简单文本**时使用
* multipart/form-data 转换为MIME报文，传送大量数据（**上传文件**等）时使用
* text/plain，简单的文本数据提交

```go
r.POST("/form/urlencoded", func(c *gin.Context) {
    name := c.PostForm("name")
    age := c.PostForm("age")

    c.JSON(http.StatusOK, gin.H{
        "name": name,
        "age":  age,
    })
})
//
r.POST("/form/multipart", func(c *gin.Context) {
    // 处理表单字段
    name := c.PostForm("name")
    age := c.PostForm("age")

    // 处理文件
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
        return
    }

    // 保存文件
    if err := c.SaveUploadedFile(file, "./"+file.Filename); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "name": name,
        "age":  age,
        "file": file.Filename,
    })
})
//
r.POST("/form/plain", func(c *gin.Context) {
    var input string
    if err := c.Bind(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": input,
    })
})
```



```go
r.POST("/form/multipart", func(c *gin.Context) {
    // 处理表单字段
    name := c.PostForm("name")
    age := c.PostForm("age")

    // 处理文件
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
        return
    }

    // 保存文件
    if err := c.SaveUploadedFile(file, "./"+file.Filename); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "name": name,
        "age":  age,
        "file": file.Filename,
    })
})
```



```go
r.POST("/form/plain", func(c *gin.Context) {
    var input string
    if err := c.Bind(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": input,
    })
})
```









































* func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine)
* func Dir(root string, listDirectory bool) http.FileSystem
* func DisableBindValidation()
* func DisableConsoleColor()
* func EnableJsonDecoderDisallowUnknownFields()
* func EnableJsonDecoderUseNumber()
* func ForceConsoleColor()
* func IsDebugging() bool
* func Mode() string
* func SetMode(value string)



# 账户管理Account





# Context

最重要的部分，用于在中间件之间传递变量、管理流程、验证请求的 JSON 并呈现 JSON 响应。

```go

```





验证json

```
func (c *Context) ShouldBindWith(obj any, b binding.Binding) error  使用指定的绑定引擎绑定传递的结构指针

func (c *Context) AbortWithStatusJSON(code int, jsonObj any)
func (c *Context) AsciiJSON(code int, obj any)
func (c *Context) BindJSON(obj any) error
func (c *Context) IndentedJSON(code int, obj any)
func (c *Context) JSON(code int, obj any)
func (c *Context) JSONP(code int, obj any)
func (c *Context) PureJSON(code int, obj any)
func (c *Context) SecureJSON(code int, obj any)
func (c *Context) ShouldBindBodyWithJSON(obj any) error  等价c.ShouldBindWith(obj, binding.JSON)
func (c *Context) ShouldBindJSON(obj any) error
```





# Engine

引擎是框架的实例，它包含复用器、中间件和配置设置。使用 New() 或 Default() 创建 Engine 实例





#  HandlerFunc

1. 请求处理：Gin框架提供了处理请求的方法和函数，包括路由处理函数、参数绑定、验证和渲染模板等。



# 错误处理



* type Error
  * func (msg Error) Error() string
  * func (msg *Error) IsType(flags ErrorType) bool
  * func (msg *Error) JSON() any
  * func (msg *Error) MarshalJSON() ([]byte, error)
  * func (msg *Error) SetMeta(data any) *Error
  * func (msg *Error) SetType(flags ErrorType) *Error
  * func (msg *Error) Unwrap() error







# Request





## 获取参数



json类型

![image-20240515214739699](./../../../../img/image-20240515214739699.png)

```go
	if err = c.ShouldBindJSON(&form); err != nil {
```





# Response













老式爵士乐唱片的数据存储库





## web开发一般步骤

1. Design API endpoints.
2. Create a folder for your code.
3. Create the data.
4. Write a handler to return all items.
5. Write a handler to add a new item.
6. Write a handler to return a specific item.



## api

/albums

* Get请求获取列表
* Post请求添加新专辑

/albums/:id

* 以Get请求获取专辑，json形式返回















## 使用gin的开源项目

- [gorush](https://github.com/appleboy/gorush): A push notification server written in Go.
- [fnproject](https://github.com/fnproject/fn): The container native, cloud agnostic serverless platform.
- [photoprism](https://github.com/photoprism/photoprism): Personal photo management powered by Go and Google TensorFlow.
- [krakend](https://github.com/devopsfaith/krakend): Ultra performant API Gateway with middlewares.
- [picfit](https://github.com/thoas/picfit): An image resizing server written in Go.
- [gotify](https://github.com/gotify/server): A simple server for sending and receiving messages in real-time per web socket.
- [cds](https://github.com/ovh/cds): Enterprise-Grade Continuous Delivery & DevOps Automation Open Source Platform.















# 中间件机制





```go
v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		v1.POST("/post", controller.CreatePostHandler) // 创建帖子

		v1.POST("/vote", controller.VoteHandler) // 投票

		v1.POST("/comment", controller.CommentHandler)    // 评论
		v1.GET("/comment", controller.CommentListHandler) // 评论列表

		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}
```











## 按照功能将id、score等数据存到redis

需要排序的，使用zset

```go
pipeline.ZAdd(ctx, enums.KeyPostTimeZSet, redis.Z{
    Score:  now,
    Member: postId,
})
```

需要取一个社区全部帖子，将id都存起来

```
pipeline.SAdd(ctx, communityKey, postId)
```

