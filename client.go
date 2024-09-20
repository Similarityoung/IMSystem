package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", severIp, severPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	return client

}

func (client *Client) DealResponse() {

	// io.Copy 是 Go 语言中用于从一个 Reader 复制数据到一个 Writer 的常用函数。
	// 它的作用是高效地将数据从输入流复制到输出流，直到遇到 EOF（文件结束）或发生错误。
	// 永久阻塞监听服务器广播消息
	_, err := io.Copy(os.Stdout, client.conn)
	if err != nil {
		return
	}
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

func (client *Client) UpdateName() bool {

	fmt.Println("Please enter your new name:")
	_, err := fmt.Scanln(&client.Name)
	if err != nil {
		return false
	}

	sendMsg := "rename|" + client.Name + "\n"
	_, err = client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (client *Client) PublicChat() {

	// send message to server
	var chatMsg string

	fmt.Println("Please enter your message, type 'exit' to exit")
	_, err := fmt.Scanln(&chatMsg)
	if err != nil {
		return
	}

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("Please enter your message, type 'exit' to exit")
		_, err := fmt.Scanln(&chatMsg)
		if err != nil {
			return
		}
	}
}

// SelectUsers search online users
func (client *Client) SelectUsers() {

	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {

	// send message to server
	var chatMsg string
	var remoteName string

	client.SelectUsers()
	fmt.Println("Please enter the username you want to chat with, type 'exit' to exit")
	_, err1 := fmt.Scanln(&remoteName)
	if err1 != nil {
		return
	}

	fmt.Println("Please enter your message, type 'exit' to exit")
	_, err := fmt.Scanln(&chatMsg)
	if err != nil {
		return
	}

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("Please enter your message, type 'exit' to exit")
		_, err := fmt.Scanln(&chatMsg)
		if err != nil {
			return
		}
	}
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
			client.PublicChat()
			break
		case 2:
			fmt.Println("Private chat")
			client.PrivateChat()
		case 3:
			fmt.Println("Rename...")
			client.UpdateName()
			break
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
	// 启动监听服务器返回消息的 goroutine
	go client.DealResponse()

	fmt.Println("Successfully connected to server")

	// 启动业务
	client.Run()
}
