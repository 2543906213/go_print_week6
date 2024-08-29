package server

import (
	"io"
	"net/http"
	"time"
)

type HttpConnection struct{}

func (hc *HttpConnection) Fetch(url string) (string, error) {
	//创建一个 http.Client 对象并设置超时时间为 10 秒。
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	//使用 client.Get 方法发送 GET 请求。
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	// 确保在函数返回之前关闭响应的主体
	defer resp.Body.Close()

	//读取响应体内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//将响应体内容转换为字符串并返回。
	return string(body), nil
}
