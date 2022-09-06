package ziface

import "net"

//定义链接模块的抽象层
type IConnection interface {
	//启动链接，让当前的链接准备开始工作
	Start()
	//停止链接，结束当前连接的工作
	Stop()
	//获取当前链接所绑定的socket
	GetTCPConnection() *net.TCPConn
	//获取当前连接模块的链接id
	GetConnID() uint32
	//获取远程客户端的tcp状态，ip，port
	RemoteAddr() net.Addr
	//发送数据，将数据发送给远程的客户端
	Send(msgid uint32,data []byte) error
}


