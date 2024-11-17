



在 Go 中，与 I/O（输入/输出）相关的模块有以下几个常用的包：

1. `fmt` 包：`fmt` 包提供了格式化的输入和输出功能。它包含了 `Print`、`Scan`、`Printf`、`Scanf` 等函数，可用于从标准输入读取数据，或将数据格式化输出到标准输出。

2. `io` 包：`io` 包提供了基本的 I/O 操作功能。它定义了 `Reader` 和 `Writer` 接口，以及一些实现这些接口的类型。通过 `io` 包，你可以进行文件读写、流处理、缓冲、复制等操作。

3. `os` 包：`os` 包提供了与操作系统交互的功能，包括文件和目录操作等。你可以使用 `os` 包来打开、创建、删除文件，读取目录内容，获取环境变量等。

4. `bufio` 包：`bufio` 包提供了缓冲读写功能，用于提高 I/O 操作的效率。它包装了 `io.Reader` 和 `io.Writer` 接口，提供了缓冲读取和写入数据的方法。

5. `ioutil` 包：`ioutil` 包提供了一些方便的 I/O 实用函数。例如，你可以使用 `ioutil.ReadFile` 来读取整个文件的内容到内存中，或使用 `ioutil.WriteFile` 将数据写入文件。

这些是 Go 中常用的与 I/O 相关的模块和包。当然，还有其他一些包可以用于特定的 I/O 操作，具体取决于你的需求。

你可以通过查阅 Go 官方文档来获取更详细的信息和使用示例。

希望这个回答能够满足你的需求。如果你还有其他问题，请随时提问。









# fmt

https://pkg.go.dev/fmt

* [func Append(b [\]byte, a ...any) []byte](https://pkg.go.dev/fmt#Append)
* [func Appendf(b [\]byte, format string, a ...any) []byte](https://pkg.go.dev/fmt#Appendf)
* [func Appendln(b [\]byte, a ...any) []byte](https://pkg.go.dev/fmt#Appendln)
* [func Errorf(format string, a ...any) error](https://pkg.go.dev/fmt#Errorf)
* [func FormatString(state State, verb rune) string](https://pkg.go.dev/fmt#FormatString)
* [func Fprint(w io.Writer, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Fprint)
* [func Fprintf(w io.Writer, format string, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Fprintf)
* [func Fprintln(w io.Writer, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Fprintln)
* [func Fscan(r io.Reader, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Fscan)
* [func Fscanf(r io.Reader, format string, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Fscanf)
* [func Fscanln(r io.Reader, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Fscanln)
* [func Print(a ...any) (n int, err error)](https://pkg.go.dev/fmt#Print)
* [func Printf(format string, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Printf)
* [func Println(a ...any) (n int, err error)](https://pkg.go.dev/fmt#Println)
* [func Scan(a ...any) (n int, err error)](https://pkg.go.dev/fmt#Scan)
* [func Scanf(format string, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Scanf)
* [func Scanln(a ...any) (n int, err error)](https://pkg.go.dev/fmt#Scanln)
* [func Sprint(a ...any) string](https://pkg.go.dev/fmt#Sprint)
* [func Sprintf(format string, a ...any) string](https://pkg.go.dev/fmt#Sprintf)
* [func Sprintln(a ...any) string](https://pkg.go.dev/fmt#Sprintln)
* [func Sscan(str string, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Sscan)
* [func Sscanf(str string, format string, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Sscanf)
* [func Sscanln(str string, a ...any) (n int, err error)](https://pkg.go.dev/fmt#Sscanln)

# io

https://pkg.go.dev/io

* func Copy(dst Writer, src Reader) (written int64, err error) 从src复制到dst，直到EOF或发生错误
* func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) 和Copy相同，只是指定了缓冲
* func CopyN(dst Writer, src Reader, n int64) (written int64, err error) 复制n个字节
* func Pipe() (*PipeReader, *PipeWriter) 创造一个内存中的管道
* func ReadAll(r Reader) ([]byte, error) 从r读取直到EOF或出现错误
* func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error) 读取n个字节，少于n个字节发生错误
* func ReadFull(r Reader, buf []byte) (n int, err error) 从r中读取buf长度
* func WriteString(w Writer, s string) (n int, err error) 将字符串写入w

为何没有提供写入字节的功能，一般由对应文件或网络包提供

writer接口

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

reader接口

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

文件读取

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	filePath := "path/to/file.txt"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	content := make([]byte, 1024)
	n, err := file.Read(content)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("File content:", string(content[:n]))
}
```

文件写入

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	filePath := "path/to/file.txt"

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	content := []byte("Hello, World!")

	n, err := file.Write(content)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Bytes written:", n)
}
```

数据流处理

```go
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	sourcePath := "path/to/source.txt"
	destinationPath := "path/to/destination.txt"

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return
	}
	defer destinationFile.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := sourceFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading data:", err)
			return
		}

		_, err = destinationFile.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error writing data:", err)
			return
		}
	}

	fmt.Println("Data copied successfully.")
}
```



# os

https://pkg.go.dev/os

类似于linux的相关命令，需要时查

文件和目录操作

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	// 创建目录
	err := os.Mkdir("path/to/directory", 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// 创建文件
	file, err := os.Create("path/to/file.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// 重命名文件
	err = os.Rename("path/to/file.txt", "path/to/newfile.txt")
	if err != nil {
		fmt.Println("Error renaming file:", err)
		return
	}

	// 删除文件
	err = os.Remove("path/to/newfile.txt")
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}

	// 删除目录
	err = os.Remove("path/to/directory")
	if err != nil {
		fmt.Println("Error removing directory:", err)
		return
	}

	fmt.Println("File and directory operations completed successfully.")
}
```

进程管理

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// 执行外部命令
	cmd := exec.Command("ls", "-l")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	fmt.Println("Command output:")
	fmt.Println(string(output))
}
```

环境变量

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	// 获取单个环境变量
	value := os.Getenv("HOME")
	fmt.Println("HOME:", value)

	// 设置环境变量
	err := os.Setenv("MY_VAR", "my_value")
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return
	}

	// 获取所有环境变量
	env := os.Environ()
	fmt.Println("All environment variables:")
	for _, v := range env {
		fmt.Println(v)
	}
}
```



# bufio

https://pkg.go.dev/bufio

* func ScanBytes(data []byte, atEOF bool) (advance int, token []byte, err error)
* func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error)
* func ScanRunes(data []byte, atEOF bool) (advance int, token []byte, err error)
* func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error)

## Reader

- type Reader
  - func NewReader(rd io.Reader) *Reader
  - func NewReaderSize(rd io.Reader, size int) *Reader
  - func (b *Reader) Buffered() int
  - func (b *Reader) Discard(n int) (discarded int, err error)
  - func (b *Reader) Peek(n int) ([]byte, error)
  - func (b *Reader) Read(p []byte) (n int, err error)
  - func (b *Reader) ReadByte() (byte, error)
  - func (b *Reader) ReadBytes(delim byte) ([]byte, error)
  - func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
  - func (b *Reader) ReadRune() (r rune, size int, err error)
  - func (b *Reader) ReadSlice(delim byte) (line []byte, err error)
  - func (b *Reader) ReadString(delim byte) (string, error)
  - func (b *Reader) Reset(r io.Reader)
  - func (b *Reader) Size() int
  - func (b *Reader) UnreadByte() error
  - func (b *Reader) UnreadRune() error
  - func (b *Reader) WriteTo(w io.Writer) (n int64, err error)





## Writer

- type Writer
  - func NewWriter(w io.Writer) *Writer
  - func NewWriterSize(w io.Writer, size int) *Writer
  - func (b *Writer) Available() int
  - func (b *Writer) AvailableBuffer() []byte
  - func (b *Writer) Buffered() int
  - func (b *Writer) Flush() error
  - func (b *Writer) ReadFrom(r io.Reader) (n int64, err error)
  - func (b *Writer) Reset(w io.Writer)
  - func (b *Writer) Size() int
  - func (b *Writer) Write(p []byte) (nn int, err error)
  - func (b *Writer) WriteByte(c byte) error
  - func (b *Writer) WriteRune(r rune) (size int, err error)
  - func (b *Writer) WriteString(s string) (int, error)





## Scanner

- type Scanner
  - func NewScanner(r io.Reader) *Scanner
  - func (s *Scanner) Buffer(buf []byte, max int)
  - func (s *Scanner) Bytes() []byte
  - func (s *Scanner) Err() error
  - func (s *Scanner) Scan() bool
  - func (s *Scanner) Split(split SplitFunc)
  - func (s *Scanner) Text() string



# ioutil

https://pkg.go.dev/io/ioutil

**从go1.6开始此包已弃用，应使用io和os包实现对应功能**



# fs

https://pkg.go.dev/io/fs#pkg-index



文件系统



























