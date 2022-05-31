package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// NewUser 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
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
	// 查询当前在线用户都有哪些
	if msg == "who" {
		user.server.mapLock.Lock()
		// 大坑！————User不能写成user不然会导致"user.SendMsg(onlineMsg)"发送给每个用户
		for _, User := range user.server.OnlineMap {
			onlineMsg := "[" + User.Addr + "]" + User.Name + ":" + "在线...\n"
			user.SendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()

		// 更改用户名 -> rename|张三
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		NewName := msg[7:]

		_, ok := user.server.OnlineMap[NewName]
		if ok {
			user.SendMsg("该用户名已被占用\n")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[NewName] = user
			user.server.mapLock.Unlock()

			user.Name = NewName
			user.SendMsg("您已更改用户名:" + user.Name + "\n")
		}

	} else {
		user.server.BroadCast(user, msg)
	}
}

// SendMsg 给当前User对应的客户端 发送消息
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// ListenMessage 监听User chan的方法，一旦有消息就发送给客户端
func (user *User) ListenMessage() {
	for {
		meg := <-user.C
		user.conn.Write([]byte(meg + "\n"))
	}
}
