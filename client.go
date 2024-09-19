package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
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
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("Failed to connect to server")
		return
	}

	fmt.Println("Successfully connected to server")

	// 启动业务
	select {}
}
