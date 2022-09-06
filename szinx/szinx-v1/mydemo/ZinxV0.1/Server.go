package main

import (
	"fmt"
	"szinx/znet"
)

//基于zinx框架来开发的服务器端应用程序
//v0.1，基本封装，只有回显功能


func main() {
	//创建一个server句柄
	s := znet.NewServer("[zinx V0.1]")
	//此处打印s.Name失败，原因是返回一个实现ziface.IServer接口的对象，不是返回一个结构体
	fmt.Printf("s = %#v",s)

	s.Serve()
}
