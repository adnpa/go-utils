

# 时间表示和获取时间

操作系统提供 wall clock 和 monotonic clock

* wall clock
  * 表示现实世界时间，可能收到同步等因素发生变化
  * 用于报时
* monotonic clock 
  * 不会收到其他因素影响，单调递增
  * 用于测量时间

使用 time.Now 会同时返回wall clock和monotonic clock ，但之后的报时操作会使用wall clock，而测量操作（加减法）会使用monotonic clock 



常见时区名

```go
UTC: 协调世界时
America/New_York: 美国东部时间（EST/EDT）
America/Los_Angeles: 美国太平洋时间（PST/PDT）
Europe/London: 英国时间（GMT/BST）
Europe/Berlin: 德国时间（CET/CEST）
Asia/Tokyo: 日本时间（JST）
Asia/Shanghai: 中国标准时间（CST）
Australia/Sydney: 澳大利亚东部时间（AEDT/AEST）
Africa/Johannesburg: 南非标准时间（SAST）
America/Sao_Paulo: 巴西时间（BRT/BRST）
```





# 定时器

## 时间间隔





## time.Timer

 用于在指**定时间后触发一次**事件。你可以使用 `time.After` 或 `time.NewTimer` 来创建一个定时器。

注：1.23版本前，需要使用Stop方法帮助gc

```go
//实现和Ticker相同的效果

go func() {
    timer := time.NewTimer(c.Option.AutoReloadInterval)
    for range timer.C {

        timer.Reset(c.Option.AutoReloadInterval)
    }
}()
```



## time.Ticker

在**指定时间间隔内重复触发**事件。它适合需要定期执行某个操作的场景。

简便写法

```go
c := time.Tick(5 * time.Second)
for next := range c {
    fmt.Printf("%v %s\n", next, statusUpdate())
}
```



#### 



### 2. 使用 `time.Ticker`

`time.Ticker` 用于在指定时间间隔内重复触发事件。它适合需要定期执行某个操作的场景。

#### 示例：使用 `time.NewTicker`

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 创建一个 Ticker，每 1 秒触发一次
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop() // 确保在程序结束时停止 Ticker

    // 使用 goroutine 来处理 Ticker
    go func() {
        for t := range ticker.C {
            fmt.Println("Tick at", t)
        }
    }()

    // 让主程序运行 5 秒
    time.Sleep(5 * time.Second)
    fmt.Println("Stopping the ticker")
}
```

### 3. 停止定时器和 Ticker

- 对于 `Timer`，可以使用 `timer.Stop()` 来停止定时器。
- 对于 `Ticker`，使用 `ticker.Stop()` 以停止触发。

### 示例：停止 Timer 和 Ticker

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 定义一个定时器
    timer := time.NewTimer(5 * time.Second)

    go func() {
        // 在 2 秒后停止定时器
        time.Sleep(2 * time.Second)
        timer.Stop()
        fmt.Println("Timer stopped")
    }()

    // 等待定时器到期
    <-timer.C
    fmt.Println("Timer expired after 5 seconds")
}
```

### 总结

- **`time.Timer`**：用于在指定延时后触发一次事件。
- **`time.Ticker`**：用于在指定时间间隔内重复触发事件。
- **停止定时器**：可以使用相应的 `Stop` 方法停止定时器或 Ticker，避免资源泄漏。

这些工具在实现定时任务、定期检查或超时控制时非常有用。根据你的需求选择合适的定时器。







- [Duration.Hours](https://pkg.go.dev/time#example-Duration.Hours)
- [Duration.Microseconds](https://pkg.go.dev/time#example-Duration.Microseconds)
- [Duration.Milliseconds](https://pkg.go.dev/time#example-Duration.Milliseconds)
- [Duration.Minutes](https://pkg.go.dev/time#example-Duration.Minutes)
- [Duration.Nanoseconds](https://pkg.go.dev/time#example-Duration.Nanoseconds)
- [Duration.Round](https://pkg.go.dev/time#example-Duration.Round)
- [Duration.Seconds](https://pkg.go.dev/time#example-Duration.Seconds)
- [Duration.String](https://pkg.go.dev/time#example-Duration.String)
- [Duration.Truncate](https://pkg.go.dev/time#example-Duration.Truncate)
- [FixedZone](https://pkg.go.dev/time#example-FixedZone)
- [LoadLocation](https://pkg.go.dev/time#example-LoadLocation)
- [Location](https://pkg.go.dev/time#example-Location)
- [Month](https://pkg.go.dev/time#example-Month)
- [NewTicker](https://pkg.go.dev/time#example-NewTicker)
- [Parse](https://pkg.go.dev/time#example-Parse)
- [ParseDuration](https://pkg.go.dev/time#example-ParseDuration)
- [ParseInLocation](https://pkg.go.dev/time#example-ParseInLocation)
- [Sleep](https://pkg.go.dev/time#example-Sleep)
- [Tick](https://pkg.go.dev/time#example-Tick)
- [Time.Add](https://pkg.go.dev/time#example-Time.Add)
- [Time.AddDate](https://pkg.go.dev/time#example-Time.AddDate)
- [Time.After](https://pkg.go.dev/time#example-Time.After)
- [Time.AppendFormat](https://pkg.go.dev/time#example-Time.AppendFormat)
- [Time.Before](https://pkg.go.dev/time#example-Time.Before)
- [Time.Date](https://pkg.go.dev/time#example-Time.Date)
- [Time.Day](https://pkg.go.dev/time#example-Time.Day)
- [Time.Equal](https://pkg.go.dev/time#example-Time.Equal)
- [Time.Format](https://pkg.go.dev/time#example-Time.Format)
- [Time.Format (Pad)](https://pkg.go.dev/time#example-Time.Format-Pad)
- [Time.GoString](https://pkg.go.dev/time#example-Time.GoString)
- [Time.Round](https://pkg.go.dev/time#example-Time.Round)
- [Time.String](https://pkg.go.dev/time#example-Time.String)
- [Time.Sub](https://pkg.go.dev/time#example-Time.Sub)
- [Time.Truncate](https://pkg.go.dev/time#example-Time.Truncate)
- [Time.Unix](https://pkg.go.dev/time#example-Time.Unix)
- [Unix](https://pkg.go.dev/time#example-Unix)
- [UnixMicro](https://pkg.go.dev/time#example-UnixMicro)
- [UnixMilli](https://pkg.go.dev/time#example-UnixMilli)









