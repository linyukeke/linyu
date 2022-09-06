package main

func main() {
	server := Newserver("127.0.0.1",8888)
	server.Start()
}

//新增用户结构体，构造函数，监听user对应的消息
//server新增onliemap和message，处理客户端上线的Handler创建并添加用户，新增广播消息的方法，新增监听广播消息的方法
//用一个goroutine单独监听message