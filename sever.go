package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer 创建一个 server 的接口
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

func (s *Server) Handler(conn net.Conn) {
	// current connection handler
	fmt.Println("connect success")
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
