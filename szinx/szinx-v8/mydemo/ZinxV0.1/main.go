package main

import (
	"fmt"
	"strconv"
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

	a := make(map[int]string)
	if _,ok := a[0];!ok{
		fmt.Println("request api,msgID="+strconv.Itoa(int(1))) //msgid是一个uint32，将它转换为int，然后使用itoa转换为字符串
	}

}
