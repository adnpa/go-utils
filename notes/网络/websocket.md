## mosquitto

https://github.com/eclipse-mosquitto/mosquitto

mqtt





# gorilla/websocket使用

## 升级处理器

升级(Upgrade)是 WebSocket 协议的概念,而不是 Go 编程的概念。

在 WebSocket 协议中,客户端和服务器之间的连接初始化时使用的是 HTTP 协议。但是,在连接建立后,双方可以将 HTTP 连接升级为 WebSocket 连接,以便使用 WebSocket 协议进行数据交换。

这个升级过程由 WebSocket 协议本身定义,不是 Go 编程语言特有的。在 Go 中,我们使用标准库 `net/http` 中提供的 `Upgrade` 函数来实现这个升级过程。

具体来说:

1. 客户端首先发送一个 HTTP 请求,请求头中包含 `Upgrade: websocket` 字段。
2. 服务端收到这个请求后,如果支持 WebSocket 协议,就会返回一个 HTTP 101 Switching Protocols 响应,表示同意升级连接。
3. 此时,HTTP 连接就被升级为 WebSocket 连接,后续的数据交换就可以使用 WebSocket 协议进行。

所以,升级(Upgrade)是 WebSocket 协议的概念,而 Go 编程只是提供了相应的 API 来实现这个过程。Go 程序员需要了解 WebSocket 协议的工作原理,然后使用 Go 标准库提供的功能来构建基于 WebSocket 的应用程序。

Upgrader 指定用于将 HTTP 连接升级为 WebSocket 连接的参数，

```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    // WebSocket连接建立成功后的处理逻辑
    handleWebSocket(conn)
})
```

实现 WebSocket 连接的读写逻辑。在 WebSocket 连接建立之后,我们可以通过 `conn.ReadMessage()` 和 `conn.WriteMessage()` 进行数据的收发。

```go
func handleWebSocket(conn *websocket.Conn) {
    defer conn.Close()
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        log.Println("Received: ", string(p))
        if err = conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }
    }
}
```

## 读取消息

读取二进制数据，连接超时后调用这个方法会报错

```go
func (c *Conn) ReadMessage() (messageType int, p []byte, err error)
```







## 写入消息



# 基于websoekcet构建IM应用

## Hub

- Hub 是 WebSocket 服务器端的核心组件,负责管理所有已连接的客户端。
- Hub 维护了所有客户端的连接状态,跟踪客户端的订阅情况。
- Hub 负责路由消息,将来自某个客户端的消息转发给订阅了该主题的其他客户端。
- Hub 可以主动向特定客户端或频道推送消息。
- Hub 处理客户端的各种操作命令,如加入/退出频道等。

```go
// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
```



## client

- 代表客户端,也就是连接到 WebSocket 服务器的用户。
- 客户端负责发送和接收 WebSocket 消息。建立连接时向hub注册，并创建读写协程
  * 读协程在连接超时或发生其他错误时退出
  * 写协程在send chan关闭或发生其他错误时退出
    由于写一般没有那么频繁，在写进程一并处理向websocket发送ping帧的逻辑
- 客户端可以订阅感兴趣的频道或主题,接收相关的消息推送。
- 客户端可以向服务器发送各种操作命令,如加入/退出频道、发送消息等。

```go
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
```



## 运行程序

## main

```go
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
    //提供html页面，作为客户端入口
	http.HandleFunc("/", serveHome)
    //
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	}) 
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```







