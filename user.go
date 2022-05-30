package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	coon net.Conn
}

// NewUser 创建一个用户的API
func NewUser(coon net.Conn) *User {
	userAddr := coon.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		coon: coon,
	}
	//启动一个go 用来监听当前User的消息
	go user.ListenMessage()
	return user
}

// ListenMessage 监听User chan的方法，一旦有消息就发送给客户端
func (user *User) ListenMessage() {
	for {
		meg := <-user.C
		user.coon.Write([]byte(meg + "\n"))
	}
}
