package main

func main() {
	server := Newserver("127.0.0.1",8888)
	server.Start()
}
//v2
//新增用户结构体，构造函数，监听user对应的消息
//server新增onliemap和message，处理客户端上线的Handler创建并添加用户，新增广播消息的方法，新增监听广播消息的方法
//用一个goroutine单独监听message


//v3
//handler新增用户消息广播

//v4
//对用户功能进行封装，和server分离
//修改user结构体，和server关联，新增online，offline，domessage方法，精简server

//v5
//新增查询在线用户
//修改user，添加senmsg函数

//v6
//修改用户名
//domessage中添加新的逻辑rename

//v7
//新增用户10秒不活跃强行踢
//每个用户的handle goroutine中，添加活跃channel，在handle中添加定时器，超时强踢 
