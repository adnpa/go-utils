## 简单使用

1. 定义配置结构体
2. 创建配置文件
3. 使用相应的解析库解析





## 复杂功能

需要配置改变触发回调等复杂功能时使用

* viper https://github.com/spf13/viper
* confd https://github.com/kelseyhightower/confd
  使用来自 etcd 或 consul 的模板和数据管理本地应用程序配置文件
* configor https://github.com/jinzhu/configor





## viper

### feature

- setting defaults
- reading from JSON, TOML, YAML, HCL, envfile and Java properties config files
- live watching and re-reading of config files (optional)
- reading from environment variables
- reading from remote config systems (etcd or Consul), and watching changes
- reading from command line flags
- reading from buffer
- setting explicit values

### 加载优先级

- explicit call to `Set`
- flag
- env
- config
- key/value store
- default



### 加载或设置配置

set方法

```go
```



explicit call to `Set`





flag

```go
serverCmd.Flags().Int("port", 1138, "Port to run Application server on")
viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))

pflag.Int("flagname", 1234, "help message for flagname")

pflag.Parse()
viper.BindPFlags(pflag.CommandLine)

i := viper.GetInt("flagname") // retrieve values from viper instead of pflag
```

env

- `AutomaticEnv()` 使用Get方法获取时会从这里获取
- `BindEnv(string...) : error` 第一个参数,key的名字,其余的是绑定到此键的环境变量的名称
- `SetEnvPrefix(string)`  AutomaticEnv和BindEnv都会使用这个前缀
- `SetEnvKeyReplacer(string...) *strings.Replacer`  自动替换,在Get中允许用 `-` ,而环境变量使用 `_` 形式
- `AllowEmptyEnv(bool)`

config

key/value store

### default

```go
viper.SetDefault("ContentDir", "content")
viper.SetDefault("LayoutDir", "layouts")
viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
```

### 读配置文件



### 写配置文件



### 其他功能































