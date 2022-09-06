package ziface




type IMsgHandle interface {
	//调度执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加具体的处理逻辑
	AddRouter(msgId uint32,route IRouter)
	StartWorkPool()
	//将消息发送给任务队列处理
	SendMsgTotaskqueue(request IRequest)
}