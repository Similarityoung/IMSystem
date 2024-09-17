package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动监听当前 user channel 消息的 goroutine
	go user.ListenMessage()

	return user
}

// ListenMessage 监听当前 user channel 的方法，一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		// 	n (int)：表示写入的字节数，即函数成功将多少字节写入目标。
		write, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			return
		}

		fmt.Println("write:", write)
	}
}
