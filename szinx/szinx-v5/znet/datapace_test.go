package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//封包拆包的单元测试,单元测试的必须以Test开头
func TestDataPack(t *testing.T){
	/*
	模拟服务器
	 */
	//创建sockettcp
	listenner, err := net.Listen("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("net.listen err",err)
		return
	}
	//创建一个go，负责从客户端处理业务
	go func(){
		//从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err !=nil {
				fmt.Println("server accept err",err)
				return
			}
			go func(conn net.Conn){
				//处理客户端请求
				//拆包过程
				//定义拆包对象
				dp := NewDataPack()
				for {
					//第一次从conn读，将包的head读出来
					headData := make([]byte,dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData) //将全部数据读取
					if err != nil{
						fmt.Println("full读取err,",err)
						return
					}
					//将headData字节流 拆包到msg中
					msgHead,err := dp.Unpack(headData)
					if err!=nil{
						fmt.Println("Unpack err",err)
						return
					}

					if msgHead.GetMsgLen() >0 {
						msg,_ := msgHead.(*Message)  //类型断言，将接口的指针转换为结构体
						//msg是有数据的，需要进行第二次读取
						msg.Data = make([]byte,msg.GetMsgLen())
						//第二次从conn读，根据包的datalen，读取data内容
						_,err := io.ReadFull(conn,msg.Data)
						if err != nil{
							fmt.Println("二次读取err",err)
							return
						}
						fmt.Println("打印结果-----")
						fmt.Println(msg)
						fmt.Println(string(msg.Data))

					}

				}

			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("dial,",err)
		return
	}
	//创建封包对象
	dp :=NewDataPack()
	//模拟粘包过程,封装两个包发送
	//封装第一个包
	msg1 := &Message{
		Id:1,
		DataLen: 3,
		Data: []byte{'a','d','d'},
	}
	sendData1,_ := dp.Pack(msg1)
	//封装第二个包
	msg2 := &Message{
		Id:2,
		DataLen: 5,
		Data: []byte{'s','w','d','f','!'},
	}
	sendData2,_ := dp.Pack(msg2)
	//将两个包黏在一起
	sendData1 = append(sendData1,sendData2...)  //打散合并为一个
	//一次性发给服务器
	conn.Write(sendData1)
	//客户端阻塞
	select {}
}
