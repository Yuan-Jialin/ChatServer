package main

/*
author:袁佳林
*/
import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("1.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("输入范围不合法")
		return false
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}
func (client *Client) SelectUsers() {
	a := "who\n"
	_, err := client.conn.Write([]byte(a))
	if err != nil {
		fmt.Println("conn Write err", err)
		return
	}
}

func (client *Client) PrivateChat() {
	client.SelectUsers()
	fmt.Println("请输入聊天对象的用户名")
	var name string
	var Msg string
	fmt.Scanln(name)
	for name != "exit" {
		fmt.Println("请输入消息内容，exit退出")
		fmt.Scanln(&Msg)
		for Msg != "exit" {
			for len(Msg) != 0 {
				sendMsg := "to|" + name + "|" + Msg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}
			Msg = ""
			fmt.Println("请输入消息内容，exit退出")
			fmt.Scanln(&Msg)
		}
		client.SelectUsers()
	}

}

func (client *Client) PublicChat() {

	var chatMsg string
	fmt.Println("请输入要发送的内容，输入exit退出")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}

		}
		chatMsg = ""
		fmt.Println("请输入要发送的内容，输入exit退出")
		fmt.Scanln(&chatMsg)

	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.Name)
	seedMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(seedMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() == false {
			switch client.flag {
			case 1:
				//fmt.Println("进入公聊模式")
				client.PublicChat()
				break

			case 2:
				//fmt.Println("进入私聊模式")
				client.PrivateChat()
				break
			case 3:
				//fmt.Println("进入更新用户名")
				client.UpdateName()
				break

			}
		}
	}
}

func main() {

	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("链接服务器失败！！！")
		return
	}
	go client.DealResponse()
	fmt.Println("链接服务器成功！！！")
	client.Run()

}
