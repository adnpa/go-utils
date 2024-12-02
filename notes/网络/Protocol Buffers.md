

https://github.com/protocolbuffers/protobuf

[https://protobuf.dev/](https://protobuf.dev/)
[https://hexdocs.pm/gpb/](https://hexdocs.pm/gpb/)
erlang: [https://github.com/tomas-abrahamsson/gpb](https://github.com/tomas-abrahamsson/gpb)


Protocol Buffers（简称 Protobuf）是一种轻量级、高效的数据序列化格式，由 Google 在内部开发中使用，并于 2008 年对外开放。与 XML 和 JSON 等传统的数据交换格式相比，Protobuf 具有更小的数据体积、更快的编解码速度和更强的可扩展性。

使用 Protobuf，开发者可以定义数据结构和消息格式，并使用 Protocol Compiler 生成对应的代码，以便在不同的平台和编程语言之间进行数据交换。Protobuf 支持多种编程语言，包括 C++、Java、Python、Go、C# 等主流语言。

Protobuf 的数据格式是二进制的，因此比文本格式更加紧凑和高效。同时，Protobuf 支持消息的版本控制、向后和向前兼容等功能，使得数据结构的演化更加灵活和可控。

总的来说，Protobuf 是一种优秀的数据序列化格式，适用于需要高效、可扩展、跨平台的数据传输和存储场景。





# 安装

下载release文件，复制到bin

```
cp protoc/bin/protoc /usr/local/bin
```

安装go插件

```插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

vscode插件

```
vscode-proto3
```



| Language                             | Source                                                       |
| ------------------------------------ | ------------------------------------------------------------ |
| C++ (include C++ runtime and protoc) | [src](https://github.com/protocolbuffers/protobuf/blob/main/src) |
| Java                                 | [java](https://github.com/protocolbuffers/protobuf/blob/main/java) |
| Python                               | [python](https://github.com/protocolbuffers/protobuf/blob/main/python) |
| Objective-C                          | [objectivec](https://github.com/protocolbuffers/protobuf/blob/main/objectivec) |
| C#                                   | [csharp](https://github.com/protocolbuffers/protobuf/blob/main/csharp) |
| Ruby                                 | [ruby](https://github.com/protocolbuffers/protobuf/blob/main/ruby) |
| Go                                   | [protocolbuffers/protobuf-go](https://github.com/protocolbuffers/protobuf-go) |
| PHP                                  | [php](https://github.com/protocolbuffers/protobuf/blob/main/php) |
| Dart                                 | [dart-lang/protobuf](https://github.com/dart-lang/protobuf)  |
| JavaScript                           | [protocolbuffers/protobuf-javascript](https://github.com/protocolbuffers/protobuf-javascript) |

## go相关库

* https://github.com/golang/protobuf
* https://github.com/protocolbuffers/protobuf-go

主要功能

* 代码生成器 protoc-gen-go
  protoc的一个编译插件，用于生成go代码
* 运行时库
  各种格式（json, text） -> 序列化格式

### 生成代码

 [generate Go specific code for a given `.proto` file](https://protobuf.dev/reference/go/go-generated).



### 运行时库













































# 编译

```
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/addressbook.proto

protoc --go_out=. *proto
```





# protobuf语法
数据类型及其使用

- int64
- string
- repeated 
- bytes 

在proto文件中定义多个message
```erlang
syntax = "proto3";

message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 results_per_page = 3;
}
```

支持类型

Protocol Buffers 支持以下基本数据类型：

1. 数值类型：包括 int32、int64、uint32、uint64、sint32、sint64、fixed32、fixed64、sfixed32、sfixed64、float、double 等。

2. 布尔类型：包括 bool 类型。

3. 字符串类型：包括 string 类型。

4. 字节类型：包括 bytes 类型。

此外，Protobuf 还支持以下复合数据类型：

1. 枚举类型：包括 enum 类型，用于表示一组离散的值。

2. 嵌套类型：包括 message 类型，用于表示一个复合数据类型。

3. 列表类型：包括 repeated 字段，用于表示一个列表或数组。

4. Map 类型：包括 map 字段，用于表示一个键值对的映射关系。

在实际应用中，可以根据具体的需求选择不同的数据类型来表示数据。例如，数值类型适用于表示数字，字符串类型适用于表示文本数据，枚举类型适用于表示离散的取值范围等。同时，使用嵌套类型和列表类型可以更好地表示复杂的数据结构，使用 Map 类型可以更好地表示键值对的映射关系。








# PB例子
```erlang
syntax = "proto3";

//-----------------------------------------------------
// 战斗请求
message fight_request_msg {
  bytes  orders = 1;
  string  version = 2;
  int32  is_log = 3; // 是(1/0)否打开写日志入
  fight_scene_msg  scene_data = 4;
}

// 战斗结果
message fight_result{
  string error = 1;       //战斗错误，正常返回 ok
  int64  winner = 2;      //进攻方胜利1/0失败
  int64  frames = 3;      //帧数
  string  md5 = 4;     //战斗数据 md5
}

//
message fight_scene_msg {
  int64  seed = 1;  //随机种子
  int64  fight_type = 2;  //战斗类型
  int64  scene_id = 3;  //战斗场景id
  fight_camp attacker = 4;  //进攻者
  fight_camp  defender = 5;  //防守者
}

//
message fight_camp{
  int64  uid = 1;  //Uid
  string name = 2;  //名字
  int32  level = 3;  //level
  int32  is_npc = 4;  //is_npc
  repeated fighter fighters = 5;  //fighters
  repeated fighter  help_fighters = 6;  //help_fighters
}

//
message fighter{
  int64  uid = 1; //武将uid
  int64  sid = 2; //武将sid
  int64  pos = 3; //站位
  repeated fight_kv_int skill = 4;  //skill
  int64  level = 5;   //武将等级
  int64  star = 6;    //武将星级
  int64  soldier_sid = 7; //带兵Sid
  int64  soldier_num = 8; //带兵数量

  int64  hp = 9;    //生命值
  int64  hp_max = 10;    //最大生命
  int64  attack = 11;    //攻击
  int64  armor = 12;    //防御
}



message fight_kv_int {
  int64 k = 1;
  int64 v = 2;
}

```







# 参考：

- [Protobuf通信协议详解：代码演示、详细原理介绍等](https://zhuanlan.zhihu.com/p/141415216)

