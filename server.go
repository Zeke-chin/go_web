package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

// Handler 业务
func (s *Server) Handler(conn net.Conn) {
	//业务...
	fmt.Println("链接成功")
}

// Start 实现服务器接口的方法
func (s *Server) Start() {
	//socket listen接口监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port)) //"127.0.0.1:8888"
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()
	for {
		//accept 接收请求
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handle 做业务
		go s.Handler(conn)
	}
}
