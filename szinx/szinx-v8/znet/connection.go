package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"szinx/ziface"
	"szinx/utils"
)

//链接模块
type Connection struct {
	//当前conn创建的时候属于那个server
	tcpServer ziface.IServer
	//socket tcp套接字
	Conn *net.TCPConn
	//链接id
	ConnID uint32
	//链接状态
	isClosed bool
	//等待退出的channel
	ExitChan chan bool
	//无缓冲管道，用于读写goroutine之间的通信
	MsgChan chan []byte
	//消息的管理msgid和对应的处理业务api关系
	MsgHandle ziface.IMsgHandle
}


func NewConnection(Tcpserver ziface.IServer,conn *net.TCPConn,connID uint32,msghandle ziface.IMsgHandle) *Connection{
	s := &Connection{
		tcpServer: Tcpserver,
		Conn: conn,
		ConnID: connID,
		MsgHandle: msghandle,
		isClosed: false,
		ExitChan: make(chan bool,1),
		MsgChan: make(chan []byte),
	}
	//将conn加入到connmanager中
	s.tcpServer.GetConmgr().Add(s)
	return s
}

//读数据的业务
func (this *Connection) StartReader() {
	fmt.Println("Reader Groutine is runing....")
	defer fmt.Println("connID=", this.ConnID)
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
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(this.GetTCPConnection(), headData); err != nil {
			fmt.Println("red msg head err:", err)
			break //如果客户端关闭，读取失败，会跳出循环，执行stop函数退出
		}
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err：", err)
			break
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(this.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data:", err)
				break
			}
		}
		msg.SetMsgData(data)

		//拆包，得到msgid和msg datalen，放在msg消息中
		//根据datalen二次读取数据，放在msg.data中
		//得到当前conn数据的request请求数据
		req := Request{
			conn: this,
			msg:  msg,
		}
		if utils.G.WorkerPoolSize > 0 {
			//已经开启工作池，将消息发送给worker工作池处理即可
			this.MsgHandle.SendMsgTotaskqueue(&req)
		} else {
			//从路由中找到注册绑定的Conn对应的route调用
			//根据绑定好的msgid，找到对应处理api业务，执行
			go this.MsgHandle.DoMsgHandler(&req)
		}
	}
}
//用来写消息，专门发送给客户端消息的模块，客户端发消息给reader，需要写，会通过管道发给write
func (this *Connection) StartWriter() {
	fmt.Println("writer goroutine is runing ....\n")
	defer fmt.Println(this.RemoteAddr().String(),"[conn writer exit!]")
	//不断阻塞等待channel的消息，进行写给客户端
	for {
		select {
		case data:= <-this.MsgChan:
			//如果有数据，写给客户端
			if _,err := this.Conn.Write(data);err != nil{
				fmt.Println("chan send data err...",err)
				return
		}
		case <-this.ExitChan:
			return
		}
	}
}

//启动链接，让当前的链接准备开始工作
func (this *Connection) Start(){
	fmt.Println("conn start()...connID=",this.ConnID)
	//启动当前链接读数据的业务
	go this.StartReader()
	//启动当前链接写数据的业务
	go this.StartWriter()

	//调用创建连接之后的hook函数
	this.tcpServer.CallOnConnStart(this)


}
//停止链接，结束当前连接的工作,资源回收
func (this *Connection) Stop(){
	fmt.Println("conn stop()...connID=",this.ConnID)
	if this.isClosed == true {
		return
	}
	this.isClosed = true

	//调用销毁连接之前的hook函数
	this.tcpServer.CallOnConnStop(this)
	// Close socket connection
	this.Conn.Close()


	//告知writer关闭
	this.ExitChan <- true
	//将当前连接从连接管理中删除
	this.tcpServer.GetConmgr().Remove(this)

	//关闭channel
	close(this.ExitChan)
	close(this.MsgChan)

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
func (this *Connection) RemoteAddr() net.Addr{
	return this.Conn.RemoteAddr()
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
	////将数据发送会客户端
	//if _, err := this.Conn.Write(binaryMsg);err !=nil{
	//	fmt.Println("write msg id,",msgid,"error:",err)
	//	errors.New("conn write msg error")
	//}
	//把消息发送给管道
	this.MsgChan <- binaryMsg
	return nil
}


