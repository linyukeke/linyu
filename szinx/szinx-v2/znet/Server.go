package znet

import (
	"errors"
	"fmt"
	"net"
	"szinx/ziface"
)

//对server实列化
//继承IServer接口
type Server struct {
//服务器的名称，ip版本，ip，port
	Name string
	IPVersion string
	IP string
	Port int
}
//写死回调函数
func CallBackToClient(conn *net.TCPConn,data []byte,cnt int) error{
	//回显业务
	fmt.Println("Conn Handle CallBackClient...")
	if _,err := conn.Write(data[:cnt]); err != nil{
		fmt.Println("write back buf err",err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

func (s *Server) Start(){
	fmt.Printf("[Start] Server Listenner at IP:%s,Port:%d,is starting\n",s.IP,s.Port)
	//将逻辑放在一个go中，防止start阻塞
	go func() {
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
			//将处理新连接的业务方法和conn绑定，得到链接模块
			dealConn := NewConnection(conn,cid,CallBackToClient)
			cid ++
			//启动当前链接业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop(){
//将服务器的资源状态或者一些已经开辟的连接信息运行停止或者回收
}

func (s *Server) Serve() {
	s.Start()
	//可以做一些额外工作，钩子之类的
	//阻塞
	select {
	}
}

func NewServer(name string) ziface.IServer{
	s := &Server{
		Name: name,
		IPVersion: "tcp4",
		IP: "0.0.0.0",
		Port: 8099,
	}
	return s

}
