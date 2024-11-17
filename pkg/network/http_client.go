package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type HttpClient struct {
	*http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{&http.Client{}}
}

func (c *HttpClient) Get(url string, header map[string]string) (*http.Response, error) {
	return c.Req("GET", url, header, nil)
}

func (c *HttpClient) PostJson(url string, header map[string]string, body map[string]any) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New("Error on convert body params ")
	}
	header["Content-Type"] = "application/json"
	return c.Req("POST", url, header, bytes.NewReader(data))
}

// PostFile 一般是服务器发送响应才发送文件
func (c *HttpClient) PostFile(url string, header map[string]string, filePath string) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "file.txt")
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %v", err)
	}

	// 将文件内容写入表单数据
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("error copying file content: %v", err)

	}

	// 关闭编写器以结束表单
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing writer %v", err)
	}

	header["Content-Type"] = writer.FormDataContentType()
	return c.Req("POST", url, header, body)
}

func (c *HttpClient) Req(method, url string, header map[string]string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		request.Header.Set(k, v)
	}
	return c.Do(request)
}
