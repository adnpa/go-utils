package markdown

import (
	"bytes"
	"github.com/yuin/goldmark"
)

func MD2Html(source []byte) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
