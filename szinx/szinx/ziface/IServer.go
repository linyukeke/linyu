package ziface


//抽象层的server，定义了各种接口，需要对外暴漏的接口，定义在抽象层


type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	//	给当前的服务提供注册一个路由的方法，供客户端的连接处理
	AddRouter(msgID uint32,router IRouter)
	GetConmgr()IConnManager
	//注册OnConnStart()
	SetOnConnStart(func(conn IConnection))
	//注册OnConnStop()
	SetOnConnStop(func(conn IConnection))
	//调用OnConnStart() //存在意义：我们不能直接调用，server是唯一的对外接口
	CallOnConnStart(conn IConnection)
	//调用OnConnStop()
	CallOnConnStop(conn IConnection)
}
