## 主要功能概览

http标准库主要提供以下六大功能特征

1. HTTP Servers
   监听请求并决定请求被如何处理
2. HTTP Clients
   发送Http请求
3. Request Handling
   提供不同结构和方法处理Http请求
4. Routing
   为不同的 URL 模式定义路由和路由处理程序
5. Middleware
   用于修改 HTTP 服务器或客户端的行为。中间件函数可以拦截请求和响应、执行其他处理并将请求传递给下一个中间件或处理程序
6. Cookies
   用于处理 HTTP cookies





# net/http

https://pkg.go.dev/net/http

* func CanonicalHeaderKey(s string) string 规范格式返回header的key
* func DetectContentType(data []byte) string 检测ContentType
* func Error(w ResponseWriter, error string, code int)  用指定错误消息和Code回复请求
* func Handle(pattern string, handler Handler)  为模式注册处理程序
* func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) 为模式注册处理方法
* func ListenAndServe(addr string, handler Handler) error 侦听Tcp网络地址并调用Serve方法
* func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error 
* func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser  限制传入请求正文的大小，类似io.LimitReader
* func NotFound(w ResponseWriter, r *Request)  使用404回复
* func ParseHTTPVersion(vers string) (major, minor int, ok bool) 解析 HTTP 版本字符串
* func ParseTime(text string) (t time.Time, err error) 解析时间标头，格式：TimeFormat、time.RFC850 和 time.ANSIC
* `func ProxyFromEnvironment(req *Request) (*url.URL, error)` 
* func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error)
* func Redirect(w ResponseWriter, r *Request, url string, code int) 重定向到url
* func Serve(l net.Listener, handler Handler) error 接收l的http连接请求，为每个服务创建goroutine并调用处理程序
* func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, ...)
* func ServeFile(w ResponseWriter, r *Request, name string) 使用指定文件或目录内容回复请求
* func ServeFileFS(w ResponseWriter, r *Request, fsys fs.FS, name string) 使用文件系统 fsys 中指定文件或目录的内容回复请求。
* func ServeTLS(l net.Listener, handler Handler, certFile, keyFile string) error
* func SetCookie(w ResponseWriter, cookie *Cookie)  将 Set-Cookie 标头添加到提供的 ResponseWriter 标头
* func StatusText(code int) string  返回Code对应文本



## Server

用于定义服务器的参数

```go
type Server struct {
　　Addr　　　　　 string
　　Handler　　　　Handler
　　ReadTimeout　　time.Duration
　　WriteTimeout　 time.Duration
　　MaxHeaderBytes int
　　TLSConfig　　　*tls.Config
　　TLSNextProto　 map[string]func(*Server, *tls.Conn, Handler)
    ConnState　　　func(net.Conn, ConnState)
　　ErrorLog　　　 *log.Logger
}
```

api

```
type Server
    func (srv *Server) Close() error 关闭所有活动网络连接，不知道被劫持的连接
    func (srv *Server) ListenAndServe() error
    func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error 使用ssl提供服务，cert为ssl证书，key为服务器上私钥  crypto可以生成证书
    func (srv *Server) RegisterOnShutdown(f func())  启动钩子函数
    func (srv *Server) Serve(l net.Listener) error
    func (srv *Server) ServeTLS(l net.Listener, certFile, keyFile string) error
    func (srv *Server) SetKeepAlivesEnabled(v bool)  是否启用keep-alives
    func (srv *Server) Shutdown(ctx context.Context) error  优雅关闭（半关闭）
```

示例

```go
s := &http.Server{
	Addr:           ":8080",
	Handler:        myHandler,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}
log.Fatal(s.ListenAndServe())
```



## 处理器和处理函数Handler

handler是一个接口，处理请求并回复响应

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

* type Handler
  * func AllowQuerySemicolons(h Handler) Handler
  * func FileServer(root FileSystem) Handler
  * func FileServerFS(root fs.FS) Handler
  * func MaxBytesHandler(h Handler, n int64) Handler
  * func NotFoundHandler() Handler
  * func RedirectHandler(url string, code int) Handler
  * func StripPrefix(prefix string, h Handler) Handler
  * func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler

示例

```go
type MyHandler struct {
}

func (handler *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

//注册
hello := handler.MyHandler{}
http.Handle("/hello", &hello)
```

处理器函数

处理器函数是和处理器具有相同行为的函数，同样接收 ResponseWriter, *Request 参数并负责处理请求和回复响应。本质上是handler的包装，会将处理器函数转换为handler

```go
func hello(w http.ResponseWriter, r *http.Request) {
　　fmt.Fprintf(w, "Hello!")
}

//注册
http.HandleFunc("/hello", hello)

//串联chaining
//为避免代码重复，可以用函数式技术来分隔日志记录、安全检查和错误处理这些横切关注点
func log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Println("Handler function called - " + name)
		h(w, r)
	}
}
func log(h http.Handler) http.Handler {  //返回handler
　　return http.HandlerFunc (func(w http.ResponseWriter, r *http.Request)
{
　　　　fmt.Printf("Handler called - %T\n", h)
　　　　h.ServeHTTP (w, r)
　　})
}
http.HandleFunc("/hello", log(hello))
```









## 多路复用器ServeMux

Web 多路复用器（Multiplexer）是一种用于处理 HTTP 请求**路由**的技术。它允许你将多个 HTTP 请求映射到不同的处理程序或处理函数上，以便根据请求的路径、方法或其他条件来执行相应的逻辑。ServeMux是一个特殊的handler实现。

* type ServeMux
  * func NewServeMux() *ServeMux
  * func (mux *ServeMux) Handle(pattern string, handler Handler)  为pattern注册处理器
  * func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) 为pattern注册处理方法
  * func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)  返回req对应的handler和pattern
  * func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)  分发处理req

```go
mux := http.NewServeMux()
files := http.FileServer(http.Dir("./public"))
mux.Handle("/static/", http.StripPrefix("/static/", files))
mux.HandleFunc("/", index)
```

如果被绑定的URL不是以/ 结尾，那么它只会与完全相同的URL匹配；但如果被绑定的URL以/ 结尾，那么即使
请求的URL只有前缀部分与被绑定URL相同，ServeMux 也会认定这两个URL是匹配的。

ServeMux一个缺陷是无法使用变量实现URL模式匹配，例如/thread/123显示123号帖子非常困难。因此可以考虑通过第三方的多路复用器。

### HttpRouter

https://github.com/julienschmidt/httprouter

这种情况下，处理器函数接收三个参数Params，Params里包括URL里的具名参数

```go
package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", p.ByName("name"))
}
func main() {
	mux := httprouter.New()
	mux.GET("/hello/:name", hello)
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

```



## 处理请求



## Request

Request结构主要包含URL、Header、Body、From、PostForm、MultipartForm等字段，提供方法如下

```
type Request
    func NewRequest(method, url string, body io.Reader) (*Request, error)
    func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error)
    func ReadRequest(b *bufio.Reader) (*Request, error)
    func (r *Request) AddCookie(c *Cookie)
    func (r *Request) BasicAuth() (username, password string, ok bool)
    func (r *Request) Clone(ctx context.Context) *Request
    func (r *Request) Context() context.Context
    func (r *Request) Cookie(name string) (*Cookie, error)
    func (r *Request) Cookies() []*Cookie
    func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
    func (r *Request) FormValue(key string) string
    func (r *Request) MultipartReader() (*multipart.Reader, error)
    func (r *Request) ParseForm() error
    func (r *Request) ParseMultipartForm(maxMemory int64) error
    func (r *Request) PathValue(name string) string
    func (r *Request) PostFormValue(key string) string
    func (r *Request) ProtoAtLeast(major, minor int) bool
    func (r *Request) Referer() string
    func (r *Request) SetBasicAuth(username, password string)
    func (r *Request) SetPathValue(name, value string)
    func (r *Request) UserAgent() string
    func (r *Request) WithContext(ctx context.Context) *Request
    func (r *Request) Write(w io.Writer) error
    func (r *Request) WriteProxy(w io.Writer) error
```

通过Request 结构的方法，用户还可以对请求报文中的cookie、引用URL以及用户代理进行访问。

URL格式及其结构

```go
// scheme://[userinfo@]host/path[?query][#fragment]
// 浏览器发送的请求会剔除fragment，这里无法获取
type URL struct {
　　Scheme　 string
　　Opaque　 string
　　User　　 *Userinfo
　　Host　　 string
　　Path　　 string
　　RawQuery string  //查询参数，需要语法分析才能获取具体值
　　Fragment string
}
```

header字段可以用map的访问方式或Get方法获取

```go
func headers(w http.ResponseWriter, r *http.Request) {
　　h := r.Header
    h := r.Header["Accept-Encoding"]  // 字符串切片形式返回 [gzip, deflate]
    h := r.Header.Get("Accept-Encoding")  // 字符串形式返回 gzip, deflate
　　fmt.Fprintln(w, h)
}
```

body是io.Reader接口，可以使用Read方法获取内容

```go
func body(w http.ResponseWriter, r *http.Request) {
　　len := r.ContentLength
　　body := make([]byte, len)
　　r.Body.Read(body)
　　fmt.Fprintln(w, string(body))
}
```

### 表单处理

Get请求参数以键值对在url中传递，而Post请求则包含在请求体中。HTML表单中的enctype字段决定Post以哪种形式传递。

* application/x-www-formurlencoded 默认值，编码为长查询字符串，以&和=分割不同键值对和键、值，例:`k1=v1&k2=v2`，传送简单文本时使用
* multipart/form-data 转换为MIME报文，传送大量数据（上传文件等）时使用
* text/plain

使用FormValue和FormFile可以更方便地处理POST方法提交的表单

1. **手动调用**ParseForm 方法或者ParseMultipartForm 方法，对请求进行语法分析
2. 根据步骤1调用的方法，访问相应的Form 字段、PostForm 字段或MultipartForm 字段。

```go
func process(w http.ResponseWriter, r *http.Request) {
　　r.ParseForm()
　　fmt.Fprintln(w, r.Form)
}

//直接访问 只会取一个值
fmt.Fprintln(w,r.FormValue("hello"))
```

PostForm字段用于获取仅包含在表单的字段，因为From字段会同时包含URL键和表单键。

```go
fmt.Fprintln(w, r.PostForm)

//直接获取
fmt.Fprintln(w, r.PostFormValue("hello"))
```

MultipartForm字段用于获取multipart/form-data编码的数据

```go
r.ParseMultipartForm(1024)
fmt.Fprintln(w, r.MultipartForm)
```

总结

![image-20240513154805993](./../../../../img/image-20240513154805993.png)

处理文件

multipart/form-data 编码，这种功能会用到file类型的input标签。

```go
func process(w http.ResponseWriter, r *http.Request) {
　　r.ParseMultipartForm(1024)
　　fileHeader := r.MultipartForm.File["uploaded"][0]
　　file, err := fileHeader.Open()
　　if err == nil {
　　　　data, err := ioutil.ReadAll(file)
　　　　if err == nil {
　　　　　　fmt.Fprintln(w, string(data))
　　　　}
　　}
}

//直接获取
file, _, err := r.FormFile("uploaded")
data, err := ioutil.ReadAll(file)
```

处理带有JSON的POST请求

除了HTML表单，还可以使用jQuery、Angular这些客户端库和框架来发送Post请求。ParseForm无法从Angular这些客户端获取，但可以从jQuery这些JavaScript库获取，因为jQuery会和HTML表单一样编码。









## 回复响应Response

### ResponseWriter

```go
type ResponseWriter interface {
    Header() Header  // 设置header
    Write([]byte) (int, error)  //设置返回的状态码
    WriteHeader(statusCode int)
}
```

ResponseWriter 是一个接口，处理器可以通过这个接口创建HTTP响应。http.Response是非导出的，用户只能使用ResponseWriter来使用。

若没有设置内容类型，会检测前512字节来确定

实现重定向示例

```go
package main

import (
	"fmt"
	"net/http"
)

func writeExample(w http.ResponseWriter, r *http.Request) {
	str := `<html>
<head><title>Go Web Programming</title></head>
<body><h1>Hello World</h1></body>
</html>`
	w.Write([]byte(str))
}
func writeHeaderExample(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	fmt.Fprintln(w, "No such service, try next door")
}
func headerExample(w http.ResponseWriter, r *http.Request) {
    //重定向到http://google.com
	w.Header().Set("Location", "http://google.com")
	w.WriteHeader(302)
}
func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/write", writeExample)
	http.HandleFunc("/writeheader", writeHeaderExample)
	http.HandleFunc("/redirect", headerExample)
	server.ListenAndServe()
}

```

向客户端返回json数据

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Post struct {
	User    string
	Threads []string
}

func writeExample(w http.ResponseWriter, r *http.Request) {
	str := `<html>
<head><title>Go Web Programming</title></head>
<body><h1>Hello World</h1></body>
</html>`
	w.Write([]byte(str))
}
func writeHeaderExample(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	fmt.Fprintln(w, "No such service, try next door")
}
func headerExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "http://google.com")
	w.WriteHeader(302)
}
func jsonExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	post := &Post{
		User:    "Sau Sheong",
		Threads: []string{"first", "second", "third"},
	}
	json, _ := json.Marshal(post)
	w.Write(json)
}
func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/write", writeExample)
	http.HandleFunc("/writeheader", writeHeaderExample)
	http.HandleFunc("/redirect", headerExample)
	http.HandleFunc("/json", jsonExample)
	server.ListenAndServe()
}

```



### CooKie

Cookie是一种在互联网上广泛使用的小型文本文件。它由网站发送到用户的浏览器，然后**存储在用户的计算机或移动设备**上。每当用户访问相同的网站时，浏览器会将存储的 cookie 发送回服务器，以便网站可以根据先前的用户活动来识别用户。

Cookie 在网站和用户之间起到了交互和跟踪用户活动的作用。它们可以用于存储用户的偏好设置、登录凭据、购物车内容等信息。网站可以使用 cookie 来提供个性化的用户体验，例如记住用户的语言选择或主题偏好。

此外，广告商和分析服务提供商也使用 cookie 来跟踪用户的浏览行为和兴趣，以便提供相关广告和统计数据。然而，随着隐私保护意识的提高，现代浏览器通常允许用户控制对 cookie 的接受与拒绝，并提供了更强大的隐私保护功能。

* 会话 cookie  没有设置Expires（http1.1废除，但为了兼容仍使用）的cookie，浏览器关闭时自动移除
* 持久 cookie

```go
type Cookie struct {
	Name  string
	Value string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite SameSite
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}
```

cookie只有两个方法

* func (c *Cookie) String() string  返回序列化后的cookie
* func (c *Cookie) Valid() error

设置cookie

```go
func setCookie(w http.ResponseWriter, r *http.Request) {
　　c1 := http.Cookie{
　　　　Name:　　 "first_cookie",
　　　　Value:　　"Go Web Programming",
　　　　HttpOnly: true,
　　}
　　c2 := http.Cookie{
　　　　Name:　　 "second_cookie",
　　　　Value:　　"Manning Publications Co",
　　　　HttpOnly: true,
　　}
　　w.Header().Set("Set-Cookie", c1.String())
　　w.Header().Add("Set-Cookie", c2.String())
    
    //快捷方法
    http.SetCookie(w, &c1)
}
```

获取cookie

```go
func getCookie(w http.ResponseWriter, r *http.Request) {
　　c1, err := r.Cookie("first_cookie")
　　if err != nil {
　　　　fmt.Fprintln(w, "Cannot get the first cookie")
　　}
　　cs := r.Cookies()
　　fmt.Fprintln(w, c1)
　　fmt.Fprintln(w, cs)
}

func getCookie(w http.ResponseWriter, r *http.Request) {
　　h := r.Header["Cookie"]
　　fmt.Fprintln(w, h)
}
```

使用cookie实现闪现消息

```go
```







## Header

handler负责处理请求并回复响应，接口如下，实现的ServeHTTP方法需要将返回写入ResponseWriter

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

* type Header
  * func (h Header) Add(key, value string)
  * func (h Header) Clone() Header
  * func (h Header) Del(key string)
  * func (h Header) Get(key string) string
  * func (h Header) Set(key, value string)
  * func (h Header) Values(key string) []string
  * func (h Header) Write(w io.Writer) error
  * func (h Header) WriteSubset(w io.Writer, exclude map[string]bool) error



## Client







## ConnState





## cookie

https://datatracker.ietf.org/doc/html/rfc6265

Cookie 代表在 HTTP 响应的 Set-Cookie 标头或 HTTP 请求的 Cookie 标头中发送的 HTTP cookie。



# HTTP2支持

1.6以上版本，以https模式启动服务器将默认使用http2，以下讨论低版本支持http2

https://pkg.go.dev/golang.org/x/net/http2

```go
go get "golang.org/x/net/http2"
```

使用

```go
http2.ConfigureServer(&server, &http2.Server{})
//检测是否成功 curl -I --http2 --insecure https://localhost:8080/
```









# html

https://pkg.go.dev/html/template



* func HTMLEscape(w io.Writer, b []byte)
* func HTMLEscapeString(s string) string
* func HTMLEscaper(args ...any) string
* func IsTrue(val any) (truth, ok bool)
* func JSEscape(w io.Writer, b []byte)
* func JSEscapeString(s string) string
* func JSEscaper(args ...any) string
* func URLQueryEscaper(args ...any) string





## Template

https://pkg.go.dev/text/template#section-documentation

### 模板语法

* 文本 会直接复制到结果
* Action 用于执行数据或控制结构 `{{ Action }}`

文本和空格

用-和空格组合可以去除文本的空格

```go
"{{23 -}} < {{- 45}}"
//生成"23<45"
```







### api

```
type Template
    func Must(t *Template, err error) *Template  创建模板帮助函数，包装错误处理过程
    func New(name string) *Template 从不同的文件系统（不仅是操作系统的文件系统）读取
    func ParseFS(fs fs.FS, patterns ...string) (*Template, error) 从不同的文件系统（不仅是操作系统的文件系统）读取
    func ParseFiles(filenames ...string) (*Template, error)  创建模板并从文件中解析

    func ParseGlob(pattern string) (*Template, error)  创建一个新模板并从模式标识的文件中解析模板定义

    func (t *Template) AddParseTree(name string, tree *parse.Tree) (*Template, error)  
    func (t *Template) Clone() (*Template, error)
    func (t *Template) DefinedTemplates() string
    func (t *Template) Delims(left, right string) *Template
    func (t *Template) Execute(wr io.Writer, data any) error  执行将已解析的模板应用于指定的数据对象，并将输出写入 wr
    func (t *Template) ExecuteTemplate(wr io.Writer, name string, data any) error
    func (t *Template) Funcs(funcMap FuncMap) *Template
    func (t *Template) Lookup(name string) *Template
    func (t *Template) Name() string
    func (t *Template) New(name string) *Template
    func (t *Template) Option(opt ...string) *Template
    func (t *Template) Parse(text string) (*Template, error)
    func (t *Template) ParseFS(fs fs.FS, patterns ...string) (*Template, error)
    func (t *Template) ParseFiles(filenames ...string) (*Template, error)
    func (t *Template) ParseGlob(pattern string) (*Template, error)
    func (t *Template) Templates() []*Template
```



```go
```











# Web服务



## SOAP

SOAP（Simple Object Access Protocol）是一种用于在网络上进行通信的**协议**，它被广泛应用于 Web 服务中。

SOAP 是一种基于 XML 的协议，它定义了一种格式和规范，用于在不同的系统之间进行结构化数据的交换。SOAP 允许应用程序在网络上通过 HTTP、SMTP 或其他协议发送和接收消息。

SOAP 消息由 SOAP Envelope（SOAP 信封）包裹，其中包含了 Header（头部）和 Body（主体）两个主要部分。Header 可选，用于传递与消息处理相关的元数据。Body 部分包含实际的数据和调用的方法。

SOAP 的主要优点是其兼容性和扩展性。它可以在不同的平台和编程语言之间进行通信，因为消息的格式是基于 XML 的。此外，SOAP 支持使用 Web Services Description Language（WSDL）定义 Web 服务的接口和操作，使得服务的调用者可以动态了解和使用服务。

然而，随着时间的推移，RESTful Web Services 和更轻量级的通信协议（如 JSON）变得更加流行，SOAP 的使用逐渐减少。RESTful 风格的 Web 服务通常更简单、易于使用，并且更适合于互联网应用的需求。

```

```







## RESTful 

RESTful Web Services（Representational State Transfer）是一种基于 REST 架构原则设计的 Web 服务。

RESTful Web Services 是通过 HTTP 协议进行通信的，它使用标准的 HTTP 方法（如 GET、POST、PUT、DELETE）来对资源进行操作，并使用 URL 定位资源。服务的响应通常以 **JSON 或 XML 格式**返回数据。

与 SOAP 相比，RESTful Web Services 具有以下优势：

1. 简化和轻量级：RESTful Web Services 使用简单的 HTTP 方法和资源 URL，减少了协议和消息的复杂性。相比之下，SOAP 使用 XML 格式的消息和复杂的协议规范。

2. 更好的可读性和可理解性：RESTful Web Services 的资源 URL 和 HTTP 方法具有自我描述的特性，使得 API 更易读、理解和使用。

3. 增加可伸缩性和性能：RESTful Web Services 通常使用无状态的通信方式，每个请求都是独立的，不需要在服务端保存状态信息。这使得服务更具可伸缩性，并且在分布式系统中更易于实现负载均衡。

4. 灵活性和兼容性：RESTful Web Services 可以使用不同的数据格式，如 JSON、XML 等。这使得它们能够与各种不同的客户端和服务端技术进行交互，具有更好的兼容性。

5. 前后端分离：RESTful Web Services 与前端应用程序的分离更为紧密，使得前端开发人员可以更灵活地使用不同的技术栈和框架。

总的来说，RESTful Web Services 更加简单、轻量级、易于理解和使用，适用于大多数互联网应用的需求。它具有更好的可伸缩性和性能，并与现代前端开发的实践更为契合。





## 创建web服务























































