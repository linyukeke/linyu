package ziface

/*
路由的抽象接口
 */

type IRouter interface {
	//业务之前
	PreHandle(request IRequest)
	//处理业务
	Handle(request IRequest)
	//业务之后
	PostHandle(request IRequest)
}

