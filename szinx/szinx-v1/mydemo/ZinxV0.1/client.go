package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(1 *time.Second)
	//直接连接远程服务器，得到一个conn连接
	conn,err := net.Dial("tcp","127.0.0.1:8099")
	if err != nil {
		fmt.Println("client connect err",err)
		return
	}
	//调用write函数写数据
	for {
		_,err := conn.Write([]byte("hello Zinx V0.1"))
		if err != nil {
			fmt.Println("write err",err)
			return
		}
		buf := make([]byte,512)
		cnt,err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err",err)
			return
		}
		fmt.Printf("server call back：%s,cnt=%d\n",buf,cnt)
		//阻塞cpu，防止cpu跑死,每隔一秒发送一次
		time.Sleep(1*time.Second)
	}
}
