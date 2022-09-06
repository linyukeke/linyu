package main

import (
	"fmt"
	"io"
	"net"
	"sync"
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
	user := Newuser(conn)

	this.maplock.Lock()
	this.OnliMap[user.Name] = user
	this.maplock.Unlock()
	//广播当前用户上线
	this.BroadCast(user,"已上线")
	//接受客户端发送的消息
	go func() {
		buf := make([]byte,4096)
		for {
			n,err := conn.Read(buf) //n是读取的字节数
			if n == 0{
				this.BroadCast(user,"下线")
				return
			}
			if err !=err && err != io.EOF{
				fmt.Println("conn Read err:",err)
				return
			}
			//提取用户消息(去除\n)
			msg := string(buf[:n-1])
			//得到的消息进行广播
			this.BroadCast(user,msg)
		}
	}()

	//阻塞handle
	select {

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


