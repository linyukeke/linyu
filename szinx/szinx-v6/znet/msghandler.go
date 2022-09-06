package znet

import (
	"fmt"
	"strconv"
	"szinx/ziface"
	//"szinx/utils"
)


type MsgHandle struct {
	apis map[uint32]ziface.IRouter
}

func NewMsgHandle()*MsgHandle{
	return &MsgHandle{
		apis: make(map[uint32]ziface.IRouter),
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







