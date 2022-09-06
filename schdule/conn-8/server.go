package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	Ip string
	Port int
	//在线用户列表
	OnliMap map[string]*User
	maplock sync.RWMutex
	//消息广播channel
	Message chan string
}

func Newserver(ip string,port int) *Server{
	server := &Server{
		Ip: ip,
		Port: port,
		OnliMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}
//监听message广播消息channel的goroutine，一旦有消息就发送给全部的在线user
func (this *Server) ListenMessage(){
	for {
		msg := <-this.Message
		this.maplock.Lock()
		for _,cli := range this.OnliMap{
			cli.C <-msg  //用户的chan
		}
		this.maplock.Unlock()

	}
}
//广播消息的方法
func (this *Server) BroadCast(user *User, msg string)  {
	sendMsg := "[" + user.Addr + "]" +user.Name + ":" +msg
	this.Message <- sendMsg
}

func (this *Server) Handle(conn net.Conn){
	//当前连接的业务
	//fmt.Println("连接建立成功\n")
	//用户上线,将用户加入onliemap
	user := Newuser(conn,this)

	//用户上线
	user.Online()
	//接受客户端发送的消息

	//监听用户是否活跃channel
	isLive := make(chan bool)
	go func() {
		buf := make([]byte,4096)
		for {
			n,err := conn.Read(buf) //n是读取的字节数
			if n == 0{
				user.Offline()
				return
			}
			if err !=err && err != io.EOF{
				fmt.Println("conn Read err:",err)
				return
			}
			//提取用户消息(去除\n)
			msg := string(buf[:n-1])
			//用户对msg进行消息处理
			user.Domessage(msg)
			//用户的任意操作，代表当前用户活跃的
			isLive <-true
		}
	}()

	//阻塞handle
	for {
		select {
		case <-isLive:
			//当前用户是活跃的，应该重置定时器
			//不做任何事情，为了激活select，更新下面的定时器
		

		case <-time.After(time.Second * 300):
			//10秒后认为超时，将当前的user强制关闭
			//重新使用time.Second方法会重置计数器
			user.SendMsg("你被踢了")
			conn.Close()
			//推出当前的handle
			runtime.Goexit()
		
		}
	}


}


func (this *Server) Start(){
	//socket listen
	listener,err := net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err != nil {
		fmt.Println("net listen err",err)
		return
	}
	//close listen socket
	defer listener.Close()

	//用一个goroutine单独监听message
	go this.ListenMessage()

	for {
		//accept
		conn,err := listener.Accept()  //conn //套接字
		if err != nil {
			fmt.Println("listener accept err:",err)
			continue
		}
		
		//do hendler
		go this.Handle(conn)
	}

}


