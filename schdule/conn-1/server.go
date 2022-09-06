package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

func Newserver(ip string,port int) *Server{
	server := &Server{
		Ip: ip,
		Port: port,
	}
	return server
}

func (this *Server) Handle(conn net.Conn){
	//当前连接的业务
	fmt.Println("连接建立成功\n")
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


