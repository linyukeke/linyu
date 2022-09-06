package main

import (
	"fmt"
	"io"
	"net"
	"szinx/znet"
	"time"
	//"szinx/ziface"
)


func main() {
	fmt.Println("client1 start...")
	time.Sleep(1 *time.Second)
	//直接连接远程服务器，得到一个conn连接
	conn,err := net.Dial("tcp","127.0.0.1:7777")
	if err != nil {
		fmt.Println("client connect err",err)
		return
	}
	//调用write函数写数据
	for {
		//发送封包的message消息
		dp := znet.NewDataPack()
		m := znet.NewMsgPackage(1,[]byte("你好"))
		//fmt.Println(m)
		//fmt.Printf("s = %#v\n",m)
		binaryMsg,err := dp.Pack(m)
		if err != nil{
			fmt.Println("pack create err",err)
			return
		}
		if _,err := conn.Write(binaryMsg);err != nil{    //err是已经申明了的变量，此处使用if，使得err仅在该if的作用域中
			fmt.Println("write err",err)
			return
		}

		//服务器回复
		//先读取流中的head部分
		binaryHead := make([]byte,dp.GetHeadLen())
		if _,err := io.ReadFull(conn,binaryHead);err != nil{
			fmt.Println("head red err",err)
			break
		}
		//将二进制的head拆包到msg结构体中
		msgHead,err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack err",err)
			break
		}
		//读取流中的data部分
		if msgHead.GetMsgLen() >0 {
			msg,_ := msgHead.(*znet.Message)  //类型断言，将接口的指针转换为结构体
			//msg是有数据的，需要进行第二次读取
			msg.Data = make([]byte,msg.GetMsgLen())
			//第二次从conn读，根据包的datalen，读取data内容
			if _,err := io.ReadFull(conn,msg.Data);err != nil{
				fmt.Println("client read data err",err)
				return
			}
			fmt.Println("---->recv Server msgid=",msg.Id,"len=",msg.DataLen,"data=",string(msg.Data))
		}
		//阻塞cpu，防止cpu跑死,每隔一秒发送一次
		time.Sleep(1*time.Second)
	}
}
