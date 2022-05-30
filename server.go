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

	//在线用户的map
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//广播消息的channel
	Message chan string
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// Handler 业务
func (s *Server) Handler(conn net.Conn) {
	//业务...
	//fmt.Println("链接成功")

	//创建连接服务器的用户
	user := NewUser(conn)

	//用户上线,将用户加入OnlineMap
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//给所以用户广播新用户上线消息
	s.BroadCast(user, "已上线")

	//接收客户端发来的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			// n是int型 用户下线之0
			if n == 0 {
				s.BroadCast(user, "已下线")
				return
			}
			//error
			if err != nil {
				fmt.Println(err)
				return
			} else if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			// 读取用户信息(除去'\n') 并广播
			msg := string(buf[:n-1])
			s.BroadCast(user, msg)
		}
	}()
	//当前Handler阻塞
	select {}
}

// Start 实现服务器接口的方法
func (s *Server) Start() {
	//socket listen 接口监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port)) //"127.0.0.1:8888"
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//启动监听message的goroutine
	go s.ListenMessage()

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

// BroadCast 广播消息的方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

// ListenMessage 监听广播消息的方法
// 用于把消息从server.Message中取出 传到给所有用户 监听消息的channel中
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}
