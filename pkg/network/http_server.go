package network

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
	router *http.ServeMux
}

type Option struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

var DefaultOption = Option{
	Addr:           ":8080",
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}

type Middleware func(http.Handler) http.Handler

func New(router *http.ServeMux) *HttpServer {
	return NewWithOption(router, DefaultOption)
}

func NewWithOption(router *http.ServeMux, option Option) *HttpServer {
	s := &http.Server{
		Addr:           option.Addr,
		ReadTimeout:    option.ReadTimeout,
		WriteTimeout:   option.WriteTimeout,
		MaxHeaderBytes: option.MaxHeaderBytes,
	}
	return &HttpServer{s, router}
}

func (s *HttpServer) Start() {
	s.server.Handler = s.router
	log.Fatal(s.server.ListenAndServe())
}

func (s *HttpServer) Handle(pattern string, handler http.Handler) {
	s.router.Handle(pattern, handler)
}

func (s *HttpServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(pattern, handler)
}

func (s *HttpServer) Stop() {
	s.server.Close()
}

func CookieHandler(w http.ResponseWriter, req *http.Request) {
	cookie := http.Cookie{
		Name:  "username",
		Value: "abc",
	}
	http.SetCookie(w, &cookie)
	fmt.Fprintln(w, "succ")
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r) // 调用下一个处理程序
	})
}
