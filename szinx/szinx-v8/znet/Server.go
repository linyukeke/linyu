package znet

import (
	"fmt"
	"net"
	"szinx/ziface"
	"szinx/utils"
)

//对server实列化
//继承IServer接口
type Server struct {
//服务器的名称，ip版本，ip，port
	Name string
	IPVersion string
	IP string
	Port int
	//当前server的消息管理模块，用来绑定msgID和对应业务api关系
	MsgHandle ziface.IMsgHandle
	//server的连接管理器
	Connmgr ziface.IConnManager
	//该server创建连接之后自动调用hook函数，OnConnStart
	OnConnStart func(conn ziface.IConnection)
	//该server创建连接之后自动调用hook函数，OnConnStop
	OnConnStop func(conn ziface.IConnection)
}


func (s *Server) Start(){
	fmt.Printf("[Start] Server Listenner at Name:%s,IP:%s,Port:%d,is starting\n",utils.G.Name,utils.G.Host,utils.G.TcpPort)
	fmt.Printf("[zinx]Version:%s,Maxconn:%d,MaxpacketSize:%d\n",utils.G.Version,utils.G.MaxConn,utils.G.MaxPackageSize)
	//将逻辑放在一个go中，防止start阻塞
	go func() {
		//开启消息队列及worker工作池
		s.MsgHandle.StartWorkPool()
		//获取一个tcp的句柄
		addr,err := net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:",err)
			return
		}
		//监听服务器的地址
		listenner,err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen",s.IPVersion,"err:",err)
			return
		}
		fmt.Println("start zinx server success",s.Name,"listening...")
		var cid uint32
		cid = 0
		//阻塞等待客户端连接，处理客户端连接业务(读写)
		for {
			//如果有客户端链接过来，阻塞会返回
			conn ,err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err",err)
				continue
			}
			//判断最大连接个数的判断，超出，则关闭新连接
			if s.Connmgr.Len() >= utils.G.MaxConn{
				fmt.Println("too many connection ;curr conn=",s.Connmgr.Len())
				//给客户端相应一个超出最大连接的错误包
				conn.Close()
				continue
			}

			//将处理新连接的业务方法和conn绑定，得到链接模块
			dealConn := NewConnection(s,conn,cid,s.MsgHandle)
			cid ++
			//启动当前链接业务
			go dealConn.Start()
		}
	}()
}
func (s *Server) GetConmgr()ziface.IConnManager{
	return s.Connmgr
}

func (s *Server) Stop(){
//将服务器的资源状态或者一些已经开辟的连接信息运行停止或者回收
	fmt.Println("[STOP] zinx server name",s.Name)
	s.Connmgr.ClearAll()
}

func (s *Server) Serve() {
	s.Start()
	//可以做一些额外工作，钩子之类的
	//阻塞
	select {
	}
}

func (s *Server) AddRouter(msgID uint32,router ziface.IRouter){
	s.MsgHandle.AddRouter(msgID,router)
	fmt.Println("add Router success...")
}

//注册OnConnStart()
func (s *Server) SetOnConnStart(hookfun func(conn ziface.IConnection)){
	s.OnConnStart = hookfun
}
//注册OnConnStop()
func (s *Server) SetOnConnStop(hookfun func(conn ziface.IConnection)){
	s.OnConnStop = hookfun
}
//调用OnConnStart() //存在意义：我们不能直接调用，server是唯一的对外接口
func (s *Server) CallOnConnStart(conn ziface.IConnection){
	if s.OnConnStart != nil {
		fmt.Println("---->>>Call ONConnStart----")
		s.OnConnStart(conn)
	}
}
//调用OnConnStop()
func (s *Server) CallOnConnStop(conn ziface.IConnection){
	if s.OnConnStop != nil {
		fmt.Println("---->>>Call ONConnStop----")
		s.OnConnStop(conn)
	}
}


func NewServer() ziface.IServer{
	s := &Server{
		Name: utils.G.Name,
		IPVersion: "tcp4",
		IP: utils.G.Host,
		Port: utils.G.TcpPort,
		MsgHandle: NewMsgHandle(),
		Connmgr: NewConnManager(),
	}
	return s

}
