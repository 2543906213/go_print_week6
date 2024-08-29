package main

import (
	"bufio"
	"fmt"
	"go_print_week6/PrintFile/server"
	"os"
	"strings"
)

func main() {
	//创建一个 HTTP 服务器实例，监听端口 9898
	var serverport server.Server
	serverport.Port = 9898

	// 使用 Goroutine 来启动 HTTP 服务器
	go serverport.HttpServer()
	//Go 不直接提供守护线程的概念，因为所有的 Goroutines 在主程序退出时都会自动终止。

	// 无限循环，等待用户输入。如果用户输入 q，循环会终止，程序退出。
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("输入q键退出: ")
		ipt, _ := reader.ReadString('\n')
		// 去掉换行符
		ipt = strings.TrimSpace(ipt)
		if ipt == "q" {
			break
		}
	}
}
