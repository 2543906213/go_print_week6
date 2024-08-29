package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Server struct {
	Port int
}

// 创建一个 HTTP 服务器实例，监听端口 9898
func (s *Server) HttpServer() {

	//这行代码将根路径 (”/”) 的所有请求映射到s.handle方法进行处理。
	//http.NewServeMux()：创建一个新的路由多路复用器（ServeMux）。
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handler)
	fmt.Printf("serving at port: %d \n", s.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), mux)
	if err != nil {
		fmt.Printf("Error starting server: %v \n", err)
	}

}

// 处理传入的 HTTP 请求
func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// 忽略请求的favicon
	if path == "/favicon.ico" {
		return
	}

	//提取路径中的参数
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	if parts != nil {

		batchId := parts[1] // 提取 batch_id
		if batchId == "" {
			fmt.Printf("打印ID不能为空 \n")
			return
			//w.Write([]byte("打印ID不能为空"))
		}
		a := Action{}
		a.PrintPath(batchId)

		result := a.PrintCode(batchId)

		// 设置响应状态码和头部
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 将数据作为 JSON 发送给客户端，并在发生错误时返回一个 500 错误。
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}

}
