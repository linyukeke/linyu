package znet

import (
	"szinx/ziface"
)

type Request struct {
	//已经和客户端建立好的链接
	conn ziface.IConnection
	//客户端请求的数据
	msg ziface.IMessage
}


func (this *Request)GetConnection() ziface.IConnection{
	return this.conn
}
//得到请求的消息数据
func (this *Request)GetData() []byte{
	return this.msg.GetMsgData()
}

func (this *Request)GetMsgId() uint32{
	return this.msg.GetMsgId()
}


