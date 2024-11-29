# gob

https://pkg.go.dev/encoding/gob

https://go.dev/blog/gob

`gob` 是 Go 语言中的一种数据编码格式，主要用于在 Go 程序之间高效地序列化和反序列化数据。它可以将 Go 数据结构转换为字节流，便于存储或通过网络传输。

典型用途是传输远程过程调用 (RPC) 的参数和结果，例如 net/rpc 提供的参数和结果。 该实现为流中的每种数据类型编译自定义编解码器，当使用单个编码器传输值流时效率最高，从而分摊了编译成本。

注意：由于gob只支持有导出字段的结构，一般使用较少





# json

https://pkg.go.dev/encoding/json@go1.22.3

https://go.dev/blog/json

json（JavaScript Object Notation）是一种简单的数据交换格式，语义上表示JavaScript的对象和列表，常用于后端和浏览器的JavaScript项目做通信，后来作为一种数据交换格式成为一种标准。

## 基础

编码和解码

```go
func Marshal(v interface{}) ([]byte, error)
func Unmarshal(data []byte, v interface{}) error
```

注意

* 当类型和json不完全匹配时，只会解码指定类型有的数据
* 结构体属性名必须是导出的（大写开头，exported）
  json里的属性名同名即可，如果不同名需要使用标签映射 `json:"Name"`

## 通用json

当不知道json的类型时，可以使用 interface{} 解码，将返回 string-interface{} 的map数据

```go
var f interface{}
err := json.Unmarshal(b, &f)
f = map[string]interface{}{
    "Name": "Wednesday",
    "Age":  6,
    "Parents": []interface{}{
        "Gomez",
        "Morticia",
    },
}
```

## 引用类型

解码传入的参数是指定类型的指针，初始为nil，若匹配Unmarshal 会分配，因此可以使用接收器类型，业务只需要判断msg是否为空

```go
type IncomingMessage struct {
    Cmd *Command
    Msg *Message
}
```

## 流编码器和解码器

包装 io.Reader 和 io.Writer 接口，用于支持 json 流

```go
dec := json.NewDecoder(os.Stdin)
enc := json.NewEncoder(os.Stdout)
for {
    var v map[string]interface{}
    if err := dec.Decode(&v); err != nil {
        log.Println(err)
        return
    }
    for k := range v {
        if k != "Name" {
            delete(v, k)
        }
    }
    if err := enc.Encode(&v); err != nil {
        log.Println(err)
    }
```

由于Readers和Writers的普遍性，这些编码器和解码器类型可以用于广泛的场景，例如读取和写入 HTTP 连接、WebSocket 或文件。

* NewDecoder

* UseNumber 解码到interface{}时，对数字类型使用Number而不是float64

* DisallowUnknownFields 遇到有不匹配的字段时返回错误

* Decode 

* Buffered 获取解码结果的reader，需要调用一次decode才能用

* More 判断当前数组是否还有更多对象

* InputOffset 输入流相对当前decoder的字节偏移量，最近返回的标记的结尾位置，以及下一个标记的开头。

* Token 返回下一个json token（json值），到达末尾返回nil, [io.EOF]

  使用看源码，用于逐个处理json值





















































# xml

https://pkg.go.dev/encoding/xml@go1.22.3



## 分析xml

1. 创建存储xml数据的结构
2. func Unmarshal(data []byte, v any) error  解封数据

```xml
<?xml version="1.0" encoding="utf-8"?>
<post id="1">
　<content>Hello World!</content>
　<author id="2">Sau Sheong</author>
</post>
```

分析代码

```go
package data

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Post struct {
	XMLName xml.Name `xml:"post"`
	Id      string   `xml:"id,attr"`
	Content string   `xml:"content"`
	Author  Author   `xml:"author"`
	Xml     string   `xml:",innerxml"`
}

type Author struct {
	Id   string `xml:"id,attr"`
	Name string `xml:",chardata"`
}

func main() {
	xmlFile, err := os.Open("post.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return
	}
	defer xmlFile.Close()
	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML data:", err)
		return
	}
	var post Post
	xml.Unmarshal(xmlData, &post)
	fmt.Println(post)
}

```

使用Decoder处理流方式传输的xml和体积较大的xml

```go
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Post struct {
	XMLName  xml.Name  `xml:"post"`
	Id       string    `xml:"id,attr"`
	Content  string    `xml:"content"`
	Author   Author    `xml:"author"`
	Xml      string    `xml:",innerxml"`
	Comments []Comment `xml:"comments>comment"`
}
type Author struct {
	Id   string `xml:"id,attr"`
	Name string `xml:",chardata"`
}
type Comment struct {
	Id      string `xml:"id,attr"`
	Content string `xml:"content"`
	Author  Author `xml:"author"`
}

func main() {
	xmlFile, err := os.Open("post.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	for {
		t, err := decoder.Token() 
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error decoding XML into tokens:", err)
			return
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "comment" {
				var comment Comment
				decoder.DecodeElement(&comment, &se)
			}
		}
	}
}

```





## 创建xml



```go
package main
import (
　"encoding/xml"
　"fmt"
　"io/ioutil"
)
type Post struct {
　XMLName xml.Name `xml:"post"`
　Id　　　string　 `xml:"id,attr"`
　Content string　 `xml:"content"`
　Author　Author　 `xml:"author"`
}
type Author struct {
　Id　 string `xml:"id,attr"`
　Name string `xml:",chardata"`
}
func main() {
　post := Post{
　Id:　　　"1",
　　Content: " Hello World!", ❶
　　Author: Author{
　　　Id:　 "2",
　　　Name: "Sau Sheong",
　　},
　}
output, err := xml.Marshal(&post)
if err != nil { ❷
　　fmt.Println("Error marshalling to XML:", err)
　　return
　}
　err = ioutil.WriteFile("post.xml", output, 0644)
　if err != nil {
　　fmt.Println("Error writing XML to file:", err)
　　return
　}
}
```







# toml

https://pkg.go.dev/github.com/burntsushi/toml#section-readme

```
go get github.com/BurntSushi/toml
```

## 基础

编码和解码

```go
type Config struct {
  Age int
  Cats []string
  Pi float64
  Perfection []int
  DOB time.Time // requires `import time`
}

var conf Config
if _, err := toml.Decode(tomlData, &conf); err != nil {
  // handle error
}
```

注意：嵌套结构在go代码中要有对应

```toml
[app]
addr = "127.0.0.1"
port = 3306
```

对应

```go
type AppConfig struct {
	App TestConfig
}

type TestConfig struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}
```





# yaml

https://pkg.go.dev/gopkg.in/yaml.v3

## 基础

注意：嵌套结构在go代码中要有对应

```yml
app:
  addr: 127.0.0.1
  port: 3306
```

对应

```go
type AppConfig struct {
	App TestConfig
}

type TestConfig struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}
```







