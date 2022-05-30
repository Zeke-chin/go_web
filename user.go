package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	coon net.Conn

	server *Server
}

// NewUser 创建一个用户的API
func NewUser(coon net.Conn, server *Server) *User {
	userAddr := coon.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		coon:   coon,
		server: server,
	}
	//启动一个go 用来监听当前User的消息
	go user.ListenMessage()
	return user
}

// OnLine 用户上线业务
func (user *User) OnLine() {
	// 把用户添加进OnlineMap
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	//广播用户上线信息
	user.DoMessage("已上线")

}

// OffLine 用户下线业务
func (user *User) OffLine() {
	// 把用户从OnlineMap删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	//广播用户下线信息
	user.DoMessage("已下线")
}

// DoMessage 用户处理信息业务
func (user *User) DoMessage(msg string) {
	user.server.BroadCast(user, msg)
}

// ListenMessage 监听User chan的方法，一旦有消息就发送给客户端
func (user *User) ListenMessage() {
	for {
		meg := <-user.C
		user.coon.Write([]byte(meg + "\n"))
	}
}
