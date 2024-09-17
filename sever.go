package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// online user list
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// message channel
	Message chan string
}

// NewServer 创建一个 server 的接口
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// ListenMessage 监听 Message channel 的方法，一旦有消息，就发送给全部在线的用户
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		// send msg to all online users
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// BroadCast 监听 Message channel 的方法，一旦有消息就发送给全部在线的用户
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + ": " + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	// current connection handler
	user := NewUser(conn)

	// save user to online map
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	// broadcast user online
	s.BroadCast(user, user.Name+" is online")
	//fmt.Println("user online:", user.Name)

	// accept messages sent by the client
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read err:", err)
				return
			} else if n == 0 {
				s.BroadCast(user, user.Name+" is offline")
				return
			}

			// get user message
			msg := string(buf[:n-1])

			// broadcast user message
			s.BroadCast(user, msg)
		}
	}()

	// 阻塞当前 handler
	select {}
}

// Start 启动 server
func (s *Server) Start() {
	// socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			fmt.Println("listen.Close err:", err)
			return
		}
	}(listen)

	// 启动监听 Message 的 goroutine
	go s.ListenMessage()

	for {
		// accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err:", err)
			return
		}

		// handle
		go s.Handler(conn)
	}

}
