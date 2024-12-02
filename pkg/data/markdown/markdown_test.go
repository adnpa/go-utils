package markdown

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestMark(t *testing.T) {
	file, err := os.Open("test.md")
	defer file.Close()
	bytes, err := io.ReadAll(file)

	str, err := MD2Html(bytes)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(str)

	err = os.WriteFile("res.html", []byte(str), 0644)
	if err != nil {
		t.Error(err)
	}
}
