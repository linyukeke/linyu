package main

import (
	"fmt"
	"szinx/znet"
)

//基于zinx框架来开发的服务器端应用程序
//v0.1，基本封装，只有回显功能
//v0.2,简单的连接封装和事务绑定
//链接的方法;连接的启动和停止，获取当前连接的conn对象套接字，得到链接，得到客户端链接的地址和端口，发送数据,链接绑定的业务处理函数
//链接的属性socket 套接字，链接的id，链接的状态，与当前链接绑定的处理业务方法，等待链接被动退出的channel



func main() {
	//创建一个server句柄
	s := znet.NewServer("[zinx V0.2]")
	//此处打印s.Name失败，原因是返回一个实现ziface.IServer接口的对象，不是返回一个结构体
	fmt.Printf("s = %#v\n",s)

	s.Serve()

}
