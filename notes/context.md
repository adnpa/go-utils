# Context

https://pkg.go.dev/context

Go语言的`context`模块提供了一种在应用程序中传递请求范围数据、控制并发和取消操作的机制。它允许你在一个请求的整个处理过程中跟踪上下文信息，并在需要时传递和使用该信息。

`context`模块的主要作用如下：

1. 传递请求范围的数据：使用`context.Context`类型，你可以将请求范围的数据附加到上下文中，并在整个调用链中传递该上下文。这对于在请求处理过程中共享**请求标识、认证信息、语言首选项**等非常有用。

2. 控制并发：`context`模块提供了管理并发操作的功能。你可以使用`context.Context`来创建带有取消和超时功能的上下文。当一个上下文被取消时，所有基于该上下文的操作都会被取消，这对于避免不必要的资源消耗和提高应用程序的可靠性非常重要。

3. 取消操作：通过使用`context.Context`，你可以方便地取消正在进行的操作。当你需要提前中止一个请求或取消一个长时间运行的操作时，可以调用上下文的`cancel`函数来取消相关的操作。

4. 超时控制：`context`模块还提供了超时控制的功能。你可以使用`context.WithTimeout`或`context.WithDeadline`函数创建一个带有超时限制的上下文。当超过指定的时间限制时，上下文将自动取消。

通过使用`context`模块，你可以更好地管理并发操作、处理请求范围的数据，并在需要时方便地取消操作。这对于构建可靠的、高效的并发应用程序非常有帮助。



## Context 应用场景:

- 处理超时和取消操作。
- 在 goroutine 之间传递请求范围内的数据。
- 在 API 服务器中管理请求生命周期。gin包装了类似的机制，handler的参数就是context，从中用c.Param等方法获取请求、连接等信息，回复用ctx.JSON方法。
- 在数据库操作中处理超时。gorm框架，context.WithTimeout(context.Background(), 2*time.Second)管理超时。
- 在 RPC 系统中控制请求的超时和取消。

## 模块方法

* [func AfterFunc(ctx Context, f func()) (stop func() bool)](https://pkg.go.dev/context#AfterFunc) 在ctx完成（取消、超时）后执行`f`函数，调用返回的 stop 函数会停止 ctx 与 f 的关联
* [func Cause(c Context) error](https://pkg.go.dev/context#Cause)  返回导致上下文取消的错误
* [func WithCancel(parent Context) (ctx Context, cancel CancelFunc)](https://pkg.go.dev/context#WithCancel) 创建一个带有取消机制的上下文
* [func WithCancelCause(parent Context) (ctx Context, cancel CancelCauseFunc)](https://pkg.go.dev/context#WithCancelCause)
* [func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)](https://pkg.go.dev/context#WithDeadline) 创建一个带有截止时间的上下文
* [func WithDeadlineCause(parent Context, d time.Time, cause error) (Context, CancelFunc)](https://pkg.go.dev/context#WithDeadlineCause)
* [func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)](https://pkg.go.dev/context#WithTimeout) 创建一个带有超时时间的上下文
* [func WithTimeoutCause(parent Context, timeout time.Duration, cause error) (Context, CancelFunc)](https://pkg.go.dev/context#WithTimeoutCause)

## 示例一：管理超时

定义一个方法等待条件变量，在ctx取消时停止等待

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	waitOnCond := func(ctx context.Context, cond *sync.Cond, conditionMet func() bool) error {
		// AfterFunc方法，在ctx取消唤醒所有进程，停止等待
		stopf := context.AfterFunc(ctx, func() {
			// We need to acquire cond.L here to be sure that the Broadcast
			// below won't occur before the call to Wait, which would result
			// in a missed signal (and deadlock).
			cond.L.Lock()
			defer cond.L.Unlock()

			// If multiple goroutines are waiting on cond simultaneously,
			// we need to make sure we wake up exactly this one.
			// That means that we need to Broadcast to all of the goroutines,
			// which will wake them all up.
			//
			// If there are N concurrent calls to waitOnCond, each of the goroutines
			// will spuriously wake up O(N) other goroutines that aren't ready yet,
			// so this will cause the overall CPU cost to be O(N²).
			cond.Broadcast()
		})
		defer stopf()

		// Since the wakeups are using Broadcast instead of Signal, this call to
		// Wait may unblock due to some other goroutine's context becoming done,
		// so to be sure that ctx is actually done we need to check it in a loop.
		for !conditionMet() {
			cond.Wait()
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}

		return nil
	}

	cond := sync.NewCond(new(sync.Mutex))

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
			defer cancel()

			cond.L.Lock()
			defer cond.L.Unlock()

			err := waitOnCond(ctx, cond, func() bool { return i%2 == 0 })
			// err := waitOnCond(ctx, cond, func() bool { return false })
			fmt.Println(err)
		}()
	}
	wg.Wait()

}

```

## 示例二：管理连接

此示例使用 AfterFunc 定义一个从 net.Conn 读取的函数，并在取消上下文时停止读取。

```go
package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func main() {
	readFromConn := func(ctx context.Context, conn net.Conn, b []byte) (n int, err error) {
		stopc := make(chan struct{})
		stop := context.AfterFunc(ctx, func() {
			conn.SetReadDeadline(time.Now())
			close(stopc)
		})
		n, err = conn.Read(b)
		if !stop() {
			// The AfterFunc was started.
			// Wait for it to complete, and reset the Conn's deadline.
			<-stopc
			conn.SetReadDeadline(time.Time{})
			return n, ctx.Err()
		}
		return n, err
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	conn, err := net.Dial(listener.Addr().Network(), listener.Addr().String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	b := make([]byte, 1024)
	_, err = readFromConn(ctx, conn, b)
	fmt.Println(err)

}

```

## 示例三：管理错误

此示例使用 AfterFunc 定义一个组合两个 Context 的取消信号的函数。

```go
package main

import (
	"context"
	"errors"
	"fmt"
)

func main() {
	// mergeCancel returns a context that contains the values of ctx,
	// and which is canceled when either ctx or cancelCtx is canceled.
	mergeCancel := func(ctx, cancelCtx context.Context) (context.Context, context.CancelFunc) {
		ctx, cancel := context.WithCancelCause(ctx)
		stop := context.AfterFunc(cancelCtx, func() {
			cancel(context.Cause(cancelCtx))
		})
		return ctx, func() {
			stop()
			cancel(context.Canceled)
		}
	}

	ctx1, cancel1 := context.WithCancelCause(context.Background())
	defer cancel1(errors.New("ctx1 canceled"))

	ctx2, cancel2 := context.WithCancelCause(context.Background())

	mergedCtx, mergedCancel := mergeCancel(ctx1, ctx2)
	defer mergedCancel()

	cancel2(errors.New("ctx2 canceled"))
	<-mergedCtx.Done()
	fmt.Println(context.Cause(mergedCtx))

}
```

## 底层结构

Go 语言中的 `context` 包的底层实现比较简单,主要由以下几个部分组成:

1. **context 结构体**:
   - `context` 结构体是 `context.Context` 接口的核心实现。
   - 它包含了以下几个关键字段:
     - `done chan struct{}`: 一个用于通知 context 被取消的通道。
     - `err error`: 表示 context 被取消的原因。
     - `deadline time.Time`: context 的截止时间。
     - `values map[interface{}]interface{}`: 存储 context 相关的键值对数据。

2. **emptyCtx 结构体**:
   - `emptyCtx` 是 `context.Background()` 和 `context.TODO()` 返回的 context 实现。
   - 它是一个单例模式的结构体,没有任何状态字段。

3. **cancelCtx 结构体**:
   - `cancelCtx` 是 `context.WithCancel()` 返回的 context 实现。
   - 它在 `context` 结构体的基础上增加了以下字段:
     - `cancel func()`: 取消当前 context 的函数。
     - `children []*cancelCtx`: 保存该 context 的所有子 context。

4. **timerCtx 结构体**:
   - `timerCtx` 是 `context.WithDeadline()` 和 `context.WithTimeout()` 返回的 context 实现。
   - 它在 `cancelCtx` 结构体的基础上增加了以下字段:
     - `deadline time.Time`: context 的截止时间。
     - `timer *time.Timer`: 用于在截止时间到达时取消 context 的定时器。

5. **valueCtx 结构体**:
   - `valueCtx` 是 `context.WithValue()` 返回的 context 实现。
   - 它在 `context` 结构体的基础上增加了以下字段:
     - `key, val interface{}`: 存储 context 的键值对。

这些结构体共同构成了 `context` 包的底层实现。当使用不同的 context 创建函数时,会返回对应的结构体实例。