package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var serverIp string
var serverPort int

// client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	Flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Flag:       -1,
	}
	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.Conn = conn
	//返回对象
	return client
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>连接服务器失败")
	}

	//创建一个goroutine 来处理server回执消息
	go client.DealResponse()

	fmt.Println(">>>>>>>连接服务器成功")

	//启动客户端业务
	client.Run()
}

func (c *Client) Run() {
	for c.Flag != 0 {
		for c.menu() != true {
		}
		//根据不同的模式 处理不同的业务
		switch c.Flag {
		case 1:
			//公聊模式
			c.PublicChat()
			break
		case 2:
			//私聊模式
			c.PrivateChat()
			break
		case 3:
			//修改用户名
			c.UpdateName()
			break
		case 0:
			//退出
			break
		}
	}
}

// menu 选择菜单
func (c *Client) menu() bool {
	var f int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.修改用户名")
	fmt.Println("0.退出")

	_, err := fmt.Scanln(&f)
	if err != nil {
		fmt.Println(">>>>>>请输入合法范围内的数字<<<<<<")
		return false
	} else if f >= 0 && f <= 3 {
		c.Flag = f
		return true
	} else {
		fmt.Println(">>>>>>请输入合法范围内的数字<<<<<<")
		return false
	}

}

// UpdateName 更新用户名方法
func (c *Client) UpdateName() bool {
	fmt.Println(">>>>>>请输入你要更改的姓名")
	_, err := fmt.Scanln(&c.Name)
	if err != nil {
		fmt.Println("Scan err:", err)
		return false
	}
	sendMSg := "rename|" + c.Name + "\n"

	_, err = c.Conn.Write([]byte(sendMSg))
	if err != nil {
		fmt.Println("Conn.Write err:", err)
		return false
	}
	return true
}

// DealResponse 处理server回应的消息 直接显示在client
func (c *Client) DealResponse() {
	//一旦client。conn有数据,就copy到stdout,并永久阻塞监听
	_, err := io.Copy(os.Stdout, c.Conn)
	if err != nil {
		fmt.Println("io.Copy err", err)
		return
	}
	/*与一下代码相等
	for{
		buf := make([]byte, 4096)
		client.conn.Read(buf)
		fmt.Println(buf)
	}
	*/
}

// PublicChat 公聊模式
func (c *Client) PublicChat() {
	//用户输入
	var chatMsg string

	fmt.Println(">>>>>>请输入聊天内容,exit推出.")
	_, err := fmt.Scanln(&chatMsg)
	if err != nil {
		fmt.Println("Scan err:", err)
		return
	}

	for chatMsg != "exit" {
		//信息发给服务器

		//消息不为空
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("Conn.Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>>请输入聊天内容,exit推出.")
		_, err := fmt.Scanln(&chatMsg)
		if err != nil {
			fmt.Println("Scan err:", err)
			return
		}
	}
}

// PrivateChat 用户私聊
func (c *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	c.SelectUsers()
	fmt.Println(">>>>>>请输入聊天对象[用户名],exit退出.")
	fmt.Scanln(&remoteName)

	for {
		fmt.Println(">>>>>>请输入消息内容,exit退出.")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//消息不为空
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"

				_, err := c.Conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("Conn.Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>>请输入消息内容,exit退出.")
			fmt.Scanln(&chatMsg)
		}
		c.SelectUsers()
		fmt.Println(">>>>>>请输入聊天对象[用户名],exit退出.")
		fmt.Scanln(&remoteName)
	}
}

// SelectUsers 查询用户列表
func (c *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := c.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("Conn.Write err:", err)
		return
	}
}
