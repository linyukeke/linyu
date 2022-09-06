package ziface

/*
IRequest接口，将客户端请求的连接信息和请求的数据封装到一个request中
 */

type IRequest interface {
	//得到当前连接数据
	GetConnection() IConnection
	//得到请求的消息数据
	GetData() []byte
	//得到请求消息的id
	GetMsgId() uint32
}


