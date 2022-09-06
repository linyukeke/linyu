package main

import (
	"fmt"
	"szinx/znet"
)

func main(){
	dp := znet.NewDataPack()
	//模拟粘包过程,封装两个包发送
	//封装第一个包
	msg1 := &znet.Message{
		Id:1,
		Data: []byte{'a','d','d'},
	}
	fmt.Println(msg1)
	fmt.Printf("s = %#v\n",msg1)
	m,_ := dp.Pack(msg1)
	fmt.Println(m)
	fmt.Printf("s = %#v\n",m)
}
