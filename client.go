package main

import (
	"flag"
	"fmt"
	"net"
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
			break
		case 2:
			//私聊模式
			break
		case 3:
			//修改用户名
			break
		case 0:
			//退出
			break
		}
	}
}

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
