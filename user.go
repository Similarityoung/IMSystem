package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	sever *Server
}

func NewUser(conn net.Conn, sever *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:  userAddr,
		Addr:  userAddr,
		C:     make(chan string),
		conn:  conn,
		sever: sever,
	}

	// 启动监听当前 user channel 消息的 goroutine
	go user.ListenMessage()

	return user
}

// Online 用户的上线业务
func (u *User) Online() {
	// save user to online map
	s := u.sever

	s.mapLock.Lock()
	s.OnlineMap[u.Name] = u
	s.mapLock.Unlock()

	// broadcast user online
	s.BroadCast(u, u.Name+" is online")
}

// Offline 用户的下线业务
func (u *User) Offline() {

	s := u.sever
	s.mapLock.Lock()
	delete(s.OnlineMap, u.Name)
	s.mapLock.Unlock()

	s.BroadCast(u, u.Name+" is offline")
}

// SendMsg 给当前用户的客户端发送消息
func (u *User) SendMsg(msg string) {
	_, _ = u.conn.Write([]byte(msg))
}

// DoMessage 用户处理消息的业务
func (u *User) DoMessage(msg string) {
	s := u.sever
	if msg == "who" {
		// 查询当前在线用户
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": online...\n"
			u.SendMsg(onlineMsg)
		}
		s.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := msg[7:]
		// 判断 newName 是否存在
		_, ok := s.OnlineMap[newName]
		if ok {
			u.SendMsg("The username is already in use\n")
		} else {
			s.mapLock.Lock()
			delete(s.OnlineMap, u.Name)
			u.Name = newName
			s.OnlineMap[newName] = u
			s.mapLock.Unlock()

			u.SendMsg("Rename successful\n")
		} // rename|nickname

	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式 to|username|message
		remoteName := strings.Split(msg, "|")[1]

		if s.OnlineMap[remoteName] == nil {
			u.SendMsg("The user does not exist\n")
			return
		}

		remoteUser := s.OnlineMap[remoteName]
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("Message cannot be empty\n")
			return
		}

		remoteUser.SendMsg(u.Name + " to you: " + content)

	} else {
		s.BroadCast(u, msg)
	}
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
