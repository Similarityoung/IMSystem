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
	flag       int
}

func NewClient(severIp string, severPort int) *Client {

	// create a client
	client := &Client{
		ServerIp:   severIp,
		ServerPort: severPort,
		flag:       999,
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

var severIp string
var severPort int

// func init() 是 Go 语言中的特殊函数，它会自动执行。
// 具体来说，init() 函数的执行时间是程序启动时，在 main() 函数运行之前执行。你不需要显式调用 init()，Go 运行时会自动调用它。
// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&severIp, "ip", "127.0.0.1", "设置服务器IP地址（默认是127.0.0.1)")
	flag.IntVar(&severPort, "port", 8888, "设置服务器端口（默认是8888）")
}

func (client *Client) menu() bool {
	var f int

	fmt.Println("1. Public chat")
	fmt.Println("2. Private chat")
	fmt.Println("3. Rename")
	fmt.Println("0. Quit")

	_, err := fmt.Scanln(&f)
	if err != nil {
		return false
	}

	if f >= 0 && f <= 3 {
		client.flag = f
		return true
	} else {
		fmt.Println("Please enter a valid option")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			fmt.Println("Public chat")
		case 2:
			fmt.Println("Private chat")
		case 3:
			fmt.Println("Rename")
		}
	}
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
	client.Run()
}
