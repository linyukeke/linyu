package main

import "net"

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn

	server *Server
}
//创建一个用户的api
func Newuser(conn net.Conn,server *Server) *User{
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn:conn,

		server: server,
	}
	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

//监听当前user channel的方法，有消息，发送给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

//用户上线
func (this *User) Online(){
	this.server.maplock.Lock()
	this.server.OnliMap[this.Name] = this
	this.server.maplock.Unlock()
	this.server.BroadCast(this,"已上线")

}
//用户下线
func (this *User) Offline(){
	this.server.maplock.Lock()
	delete(this.server.OnliMap,this.Name)
	this.server.maplock.Unlock()
	this.server.BroadCast(this,"下线")

}

//给当前user对应的客户端发送消息
func (this *User) SendMsg(msg string){
	this.conn.Write([]byte(msg))
}

//用户处理消息的业务
func (this *User) Domessage(msg string){
	if msg == "who"{
		//查询当前在线用户有哪些
		this.server.maplock.Lock()
		for _,user := range this.server.OnliMap{
			onlimsg := "[" + user.Addr + "]" + user.Name +": 在线...\n"
			this.SendMsg(onlimsg)
		}
		this.server.maplock.Unlock()
	} else {
		this.server.BroadCast(this, msg)
	}
}

//



