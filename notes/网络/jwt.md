https://pkg.go.dev/github.com/golang-jwt/jwt/v5

https://github.com/golang-jwt/jwt

https://golang-jwt.github.io/jwt/usage/create/

https://datatracker.ietf.org/doc/html/rfc7519

https://docs.authing.cn/v2/concepts/jwt-token.html

https://en.wikipedia.org/wiki/JSON_Web_Token

https://jwt.io/introduction

https://www.ruanyifeng.com/blog/2018/07/json_web_token-tutorial.html





```go
go get -u github.com/golang-jwt/jwt/v5
```





# 原理

![jwt.drawio](./../../../../img/jwt.drawio.png)

JWT (JSON Web Token) 是一种开放标准,用于在网络应用程序中以紧凑和安全的方式传输信息。它的工作原理如下:

1. **Token 结构**: xxxxx.yyyyy.zzzzz
   
   JWT 由三个部分组成:头部(Header)、有效载荷(Payload)和签名(Signature)。
   
   - 头部描述了 JWT 的类型和使用的加密算法。
   - 有效载荷包含了需要传输的声明(claims)信息,如用户 ID、角色等。
   - 签名部分用于验证 JWT 的完整性和来源。
   
2. **Token 生成**:
   
   - 当用户进行身份验证(如登录)时,服务器会根据预共享的密钥(secret)对头部和有效荷载进行签名,生成完整的 JWT。
   - 生成的 JWT 会返回给客户端(如浏览器)。
   
3. **Token 传输**:
   
   - 客户端在后续的请求中,访问**受保护的资源**时，会将 JWT 放在 HTTP 头部的 `Authorization` 字段中发送给服务器。
   - 通常使用 `Bearer` 方案,如 `Authorization: Bearer <token>`.
   
4. **Token 验证**:
   
   - 当服务器收到带有 JWT 的请求时,会使用预共享的密钥对 JWT 的签名部分进行验证。
   - 验证通过后,服务器就可以信任 JWT 中包含的声明信息,进而授权用户访问相应的资源。
   
5. **Token 更新**:
   
   - JWT 通常有一个有效期(expiration time),过期后需要重新获取新的 JWT。
   - 客户端可以在 JWT 即将过期时,向服务器发送刷新请求,以获取新的 JWT。

JWT 的工作原理具有以下特点:

- **无状态**:JWT 是自包含的,服务器不需要保存任何会话信息。
- **安全性**:签名可以防止 JWT 被篡改,确保了数据的完整性和真实性。
- **可扩展性**:由于无状态,JWT 非常适合微服务和分布式系统架构。
- **跨域支持**:JWT 可以跨域传输,为单点登录(SSO)等应用场景提供支持。

总之,JWT 提供了一种简单、安全且高效的方式,用于在应用程序间传输经过身份验证的信息。





简而言之，它是一个经过签名的 JSON 对象，可以执行一些有用的操作（例如身份验证）。它通常用于 Oauth 2 中的承载令牌。令牌由三部分组成，以 . 分隔。前两部分是 JSON 对象，已进行 base64url 编码。最后一部分是签名，以相同的方式编码。

第一部分称为标题。它包含验证最后一部分（签名）所需的信息。例如，签名使用哪种加密方法以及使用什么密钥。

中间的部分是有趣的部分。它称为声明，包含您真正关心的内容。有关保留密钥以及添加您自己的密钥的正确方法的信息，请参阅 RFC 7519







# go jwt

## 生成令牌

* **exp (Expiration Time)**: 令牌的过期时间。这是一个 Unix 时间戳,表示从 1970 年 1 月 1 日 00:00:00 UTC 开始的秒数。当前时间超过此时间戳时,令牌将被视为无效。
* **iat (Issued At)**: 令牌的签发时间。这也是一个 Unix 时间戳,表示令牌创建的时间。
* **nbf (Not Before)**: 令牌的生效时间。这是一个 Unix 时间戳,表示在此时间之前,该令牌无法被接受进行验证。
* **iss (Issuer)**: 令牌的签发者。这通常是一个字符串,表示签发该令牌的实体。
* **sub (Subject)**: 令牌的主体,通常是一个字符串,表示该令牌代表的用户或实体。
* **aud (Audience)**: 令牌的受众,通常是一个字符串或字符串数组,表示该令牌被设计用于的一个或多个接收者。

```go
```





## 解析令牌

```
func (p *Parser) Parse(tokenString string, keyFunc Keyfunc) (*Token, error)  Parse 解析、验证、验证签名并返回解析后的令牌。 keyFunc 将接收解析后的令牌并应返回用于验证的密钥。
func (p *Parser) ParseUnverified(tokenString string, claims Claims) (token *Token, parts []string, err error)  解析令牌但不验证签名
func (p *Parser) ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error)  与 Parse 类似，但提供了一个实现 Claims 接口的默认对象 用于实现
```











在 Gin 框架中使用 JWT (JSON Web Token) 进行身份验证和授权的步骤如下:

1. **引入 JWT 库**:
   - 在您的 Gin 项目中,添加 JWT 库依赖,例如 `github.com/dgrijalva/jwt-go`。

2. **定义 JWT 密钥**:
   - 创建一个用于签名 JWT 的密钥,通常使用一个安全的随机字符串。
   - 将此密钥保存在一个安全的地方,不要将其直接硬编码在代码中。

3. **创建登录处理程序**:
   - 在登录时,验证用户的凭证(如用户名和密码)。
   - 如果验证成功,使用 JWT 库生成一个 JWT 令牌,并将其返回给客户端。
   - 例如:
     ```go
     func Login(c *gin.Context) {
       // 验证用户凭证
       // ...
       
       // 生成 JWT 令牌
       token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
         "user_id": user.ID,
         "exp":     time.Now().Add(time.Hour * 24).Unix(), // 设置过期时间为 24 小时
       })
       
       // 使用密钥签名 JWT 令牌
       tokenString, err := token.SignedString([]byte(jwtSecret))
       if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
         return
       }
       
       c.JSON(http.StatusOK, gin.H{"token": tokenString})
     }
     ```

4. **创建中间件验证 JWT 令牌**:
   - 在需要授权的路由上使用中间件来验证 JWT 令牌。
   - 从 HTTP 头部中提取 JWT 令牌,并使用签名密钥进行验证。
   - 如果验证成功,将用户信息存储在 gin.Context 中供后续处理使用。
   - 例如:
     ```go
     func JWTAuthMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
         // 从 HTTP 头部获取 JWT 令牌
         tokenString := c.GetHeader("Authorization")
         
         // 使用签名密钥验证 JWT 令牌
         token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
           // 验证签名算法
           if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
             return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
           }
           return []byte(jwtSecret), nil
         })
         
         if err != nil || !token.Valid {
           c.AbortWithStatus(http.StatusUnauthorized)
           return
         }
         
         // 将用户信息存储在 gin.Context 中
         claims := token.Claims.(jwt.MapClaims)
         c.Set("user_id", claims["user_id"])
         
         c.Next()
       }
     }
     ```

5. **在路由中使用中间件**:
   - 在需要授权的路由组中使用 `JWTAuthMiddleware()` 中间件。
   - 例如:
     ```go
     r := gin.Default()
     
     // 使用 JWT 中间件的路由组
     authorized := r.Group("/")
     authorized.Use(JWTAuthMiddleware())
     {
       authorized.GET("/profile", GetUserProfile)
       authorized.POST("/posts", CreatePost)
     }
     ```

通过以上步骤,您可以在 Gin 框架中集成 JWT 进行身份验证和授权。客户端在每个需要授权的请求中,需要在 HTTP 头部中包含 JWT 令牌。服务端会验证令牌的有效性,并将用户信息存储在 `gin.Context` 中供后续处理使用。