package main

import (
	"net"
	"strings"
)

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
	}else if len(msg) >7 && msg[:7] == "rename|"{
		//消息格式 rename|张三
		newName := strings.Split(msg,"|")[1]
		//判断张三是否存在
		_, ok := this.server.OnliMap[newName]
		if ok {
			this.SendMsg("当前用户名被使用\n")
		} else {
			this.server.maplock.Lock()
			delete(this.server.OnliMap,this.Name)
			this.server.OnliMap[newName]=this
			this.server.maplock.Unlock()

			this.Name = newName
			this.SendMsg("您已更新用户名:" + this.Name + "\n")

		}

	} else if len(msg) >4 && msg[:3] == "to|"{
		//消息格式：to|张三|消息内容

		//1 获取用户名
		remoteName := strings.Split(msg,"|")[1]
		if remoteName == ""{
			this.SendMsg("消息格式不正确，请使用 \"to|张三|你好\"格式\n")
			return
		}
		//2根据用户名得到user对象
		remoteUser,ok := this.server.OnliMap[remoteName]
		if !ok {
			this.SendMsg("该用户不存在\n")
			return
		}
		//3获取消息内容，发送消息
		content := strings.Split(msg,"|")[2]
		if content == ""{
			this.SendMsg("无消息内容，请重发\n")
			return
		}
		remoteUser.SendMsg(this.Name + "对你说" + content)

	} else {
		this.server.BroadCast(this, msg)
	}
}

//



