package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"szinx/ziface"
	//"szinx/utils"
)

//链接模块
type Connection struct {
	//socket tcp套接字
	Conn *net.TCPConn
	//链接id
	ConnID uint32
	//链接状态
	isClosed bool
	//等待退出的channel
	ExitChan chan bool
	//该链接处理的方法
	Router ziface.IRouter
}


func NewConnection(conn *net.TCPConn,connID uint32,router ziface.IRouter) *Connection{
	s := &Connection{
		Conn: conn,
		ConnID: connID,
		Router: router,
		isClosed: false,
		ExitChan: make(chan bool,1),
	}
	return s
}

//读数据的业务
func (this *Connection) StartReader(){
	fmt.Println("Reader Groutine is runing....")
	defer fmt.Println("connID=",this.ConnID)
	defer this.Stop()
	for {
		//读取客户数据到buf中
		//buf := make([]byte,utils.G.MaxPackageSize)
		//_,err := this.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("revc buf err",err)
		//	//break //跳出循环
		//	continue //跳出本次循环
		//}
		//创建一个拆包解包的对象
		dp := NewDataPack()
		//读取客户端的msg head,二进制流8个字节
		headData := make([]byte,dp.GetHeadLen())
		if  _, err := io.ReadFull(this.GetTCPConnection(), headData);err != nil{
			fmt.Println("red msg head err:",err)
			break
		}
		msg,err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err：",err)
			break
		}
		var data []byte
		if msg.GetMsgLen() >0 {
			data = make([]byte,msg.GetMsgLen())
			if _, err := io.ReadFull(this.GetTCPConnection(), data);err != nil{
				fmt.Println("read msg data:",err)
				break
			}
		}
		msg.SetMsgData(data)

		//拆包，得到msgid和msg datalen，放在msg消息中
		//根据datalen二次读取数据，放在msg.data中
		//得到当前conn数据的request请求数据
		req := Request{
			conn: this,
			msg: msg,
		}

		//执行注册的路由方法，关注
		go func(request ziface.IRequest) {
			//从路由中找到当前对应的方法
			this.Router.PreHandle(request)
			this.Router.Handle(request)
			this.Router.PostHandle(request)
		}(&req)

	}
}

//启动链接，让当前的链接准备开始工作
func (this *Connection) Start(){
	fmt.Println("conn start()...connID=",this.ConnID)
	//启动当前链接读数据的业务
	this.StartReader()
	//启动当前链接写数据的业务

}
//停止链接，结束当前连接的工作
func (this *Connection) Stop(){
	fmt.Println("conn stop()...connID=",this.ConnID)
	if this.isClosed == true {
		return
	}
	this.isClosed = true
	// Close socket connection
	this.Conn.Close()

	//关闭channel
	close(this.ExitChan)

}
//获取当前链接所绑定的socket
func (this *Connection) GetTCPConnection() *net.TCPConn{
	return this.Conn
}

//获取当前连接模块的链接id
func (this *Connection) GetConnID() uint32{
	return this.ConnID
}
//获取远程客户端的tcp状态，ip，port
func (this *Connection) RemoreAddr() net.Addr{
	return this.RemoreAddr()
}
//发送数据，将数据发送给远程的客户端
func (this *Connection) Send(msgid uint32,data []byte) error{
	if this.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	//将data进行封包,MsgDatalen,msgId,data
	dp := NewDataPack()
	binaryMsg,err := dp.Pack(NewMsgPackage(msgid,data))
	if err != nil {
		fmt.Println("pack error:",err)
		return errors.New("pack error msg")
	}
	//将数据发送会客户端
	if _, err := this.Conn.Write(binaryMsg);err !=nil{
		fmt.Println("write msg id,",msgid,"error:",err)
		errors.New("conn write msg error")
	}
	return nil
}


