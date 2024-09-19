package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

var severIp string
var severPort int

// func init() 是 Go 语言中的特殊函数，它会自动执行。
// 具体来说，init() 函数的执行时间是程序启动时，在 main() 函数运行之前执行。你不需要显式调用 init()，Go 运行时会自动调用它。
func init() {
	flag.StringVar(&severIp, "ip", "127.0.0.1", "设置服务器IP地址（默认是127.0.0.1)")
	flag.IntVar(&severPort, "port", 8888, "设置服务器端口（默认是8888）")
}

func NewClient(severIp string, severPort int) *Client {

	// create a client
	client := &Client{
		ServerIp:   severIp,
		ServerPort: severPort,
	}

	// connect to server
	conn, err := net.Dial("TCP", fmt.Sprintf("%s:%d", severIp, severPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	return client

}

func main() {

	// command line parsing
	flag.Parse()

	client := NewClient(severIp, severPort)

	if client == nil {
		fmt.Println("Failed to connect to server")
		return
	}

	fmt.Println("Successfully connected to server")

	// 启动业务
	select {}
}
