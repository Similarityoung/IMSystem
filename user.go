package main

import "net"

type user struct {
	Name string
	Addr int
	C    chan string
	conn net.Conn
}
