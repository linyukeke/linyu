package ziface


//抽象层的server，定义了各种接口


type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	//	给当前的服务提供注册一个路由的方法，供客户端的连接处理
	AddRouter(msgID uint32,router IRouter)

}
