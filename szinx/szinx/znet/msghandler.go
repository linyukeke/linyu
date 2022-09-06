package znet

import (
	"fmt"
	"strconv"
	"szinx/ziface"
	"szinx/utils"
)


type MsgHandle struct {
	//存放每个msgid对应的方法
	apis map[uint32]ziface.IRouter
	//消息队列的数目，默认和工作池数目一致
	Taskqueue []chan ziface.IRequest  //一个存储ziface.IRequest类型的chan的切片
	//负责工作worker池的worker数量
	WorkerPoolSize uint32

}

func NewMsgHandle()*MsgHandle{
	return &MsgHandle{
		apis: make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.G.WorkerPoolSize,
		Taskqueue: make([]chan ziface.IRequest,utils.G.WorkerPoolSize),   //初始化切片
	}
}

//调度执行对应的Router消息处理方法
func (this *MsgHandle)DoMsgHandler(request ziface.IRequest){
	//从request中找到MsgID
	handler,ok := this.apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID=",request.GetMsgId(),"is not found")
	}
	//根据MsgID，调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
//为消息添加具体的处理逻辑
func (this *MsgHandle)AddRouter(msgId uint32,route ziface.IRouter){
	//判断当前msg绑定的api处理方法是否已经存在
	if _,ok := this.apis[msgId];ok{       //判断map是否有值的方法
		panic("request api,msgID="+strconv.Itoa(int(msgId))) //msgid是一个uint32，将它转换为int，然后使用itoa转换为字符串
	}
	//添加
	this.apis[msgId] = route
	fmt.Println("api msgid=",msgId,"succ!")

}
//将消息交给taskqueue，由work进行处理
func (this *MsgHandle)SendMsgTotaskqueue(request ziface.IRequest){
	//1将消息平均分配给不通的worker
	//根据连接id进行分配，轮询
	workid := request.GetConnection().GetConnID() % this.WorkerPoolSize
	fmt.Println("add connID=",request.GetConnection().GetConnID(),"request MsgID = ",request.GetMsgId(),"to workerID=",workid,".....\n")
	//2将消息发送给对应的worker的taskqueue即可。
	this.Taskqueue[workid] <- request
}

//启动一个worker工作池(开启工作池的动作只能发生一次，一个zinx框架只能只能有一个worker工作池)
func (this *MsgHandle)StartWorkPool(){
	//根据workerpoolsize,分别开启worker，每个worker用一个go来承载
	for i:=0; i<int(this.WorkerPoolSize); i++{
		//启动一个worker
		//给当前的worker对应的channel消息队列开辟空间，即每个队列的最大长度
		this.Taskqueue[i] = make(chan ziface.IRequest,utils.G.MaxWorkerTask)
		//启动当前的worker，阻塞等待消息从channel传递进来
		go this.StartOneWorker(i,this.Taskqueue[i])
	}
}

//启动一个worker工作流程
func (this *MsgHandle)StartOneWorker(WorkerID int,taskqueue chan ziface.IRequest){
	fmt.Println("workid=",WorkerID,"is started.....")
	//不断阻塞等待对应消息队列的消息
	for{
		select {
		//如果有消息过来，出列的就是一个客户端的request，执行当前request所绑定的业务。
		case request:=<- taskqueue:
			this.DoMsgHandler(request)
		}
	}

}





