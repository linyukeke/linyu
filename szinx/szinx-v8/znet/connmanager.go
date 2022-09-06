package znet


import (
	"errors"
	"fmt"
	"sync"
	"szinx/ziface"
	//"szinx/utils"
)

type ConnManager struct {
	conntions map[uint32]ziface.IConnection
	connLock sync.RWMutex  //保护连接集合的读写锁
}

//创建当前连接的方法
func NewConnManager()  *ConnManager{
	return &ConnManager{
		conntions: make(map[uint32] ziface.IConnection),
	}
}








//添加连接
func (this *ConnManager)Add(Conn ziface.IConnection){
	//保护共享资源map，上写锁
	this.connLock.Lock()
	defer this.connLock.Unlock()
	//将conn加入到connmanager中
	this.conntions[Conn.GetConnID()] = Conn
	fmt.Println("connecion ID:",Conn.GetConnID(),"add to ConnManager succ;now conn num=",this.Len())
}
//删除连接
func (this *ConnManager)Remove(Conn ziface.IConnection){
	//保护共享资源map，上写锁
	this.connLock.Lock()
	defer this.connLock.Unlock()
	delete(this.conntions,Conn.GetConnID())
	fmt.Println("connecion ID:",Conn.GetConnID()," delete to ConnManager succ;now conn num=",this.Len())
}
//获取连接
func (this *ConnManager)Get(connID uint32)(ziface.IConnection,error){
	//保护共享资源map，上读锁
	this.connLock.RLock()
	defer this.connLock.RUnlock()
	if conn,ok := this.conntions[connID];ok{
		return conn,nil
	}else {
		return nil,errors.New("connection not FOUNT!!!")
	}
}
//得到连接总数
func (this *ConnManager)Len()int{
	return len(this.conntions)
}
//清除所有链接
func (this *ConnManager)ClearAll(){
	//保护共享资源map，上写锁
	this.connLock.Lock()
	defer this.connLock.Unlock()
	//删除conn，并停止conn的工作
	for connID,conn := range this.conntions{
		//停止
		conn.Stop()
		//删除
		delete(this.conntions,connID)
	}
}
