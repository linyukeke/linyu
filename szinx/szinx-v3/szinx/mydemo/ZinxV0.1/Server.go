package main
import (
	"fmt"
	"szinx/znet"
	"szinx/ziface"
)

//基于zinx框架来开发的服务器端应用程序
//v0.1，基本封装，只有回显功能
//v0.2,简单的连接封装和事务绑定
//链接的方法;连接的启动和停止，获取当前连接的conn对象套接字，得到链接，得到客户端链接的地址和端口，发送数据,链接绑定的业务处理函数
//链接的属性socket 套接字，链接的id，链接的状态，与当前链接绑定的处理业务方法，等待链接被动退出的channel
//v0.3 request封装，将连接和数据绑定在一起
//属性：链接，请求数据，
//方法：得到当前链接，得到链接数据，新建request请求
//route模块
//方法：处理业务的方法，处理业务之前的钩子，处理业务之后的钩子
//基础的baseroute：它实现了抽象层的所有方法，可以根据需求重写该方法
//IServer需要添加route方法，Server类增添route成员，connection类添加route成员，在connection调用，已经注册的route处理业务

type PingRouter struct {
	znet.BaseRouter
}
//业务之前
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call back PreHandle.....")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("开始。。。。。。\n"))
	if err != nil {
		fmt.Println("Call back ping err",err)
		return
	}
}
//处理业务
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call back Handle.....")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("ping。。。。。。\n"))
	if err != nil {
		fmt.Println("ping err",err)
		return
	}
}
//业务之后
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call back PostHandle.....")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("最后。。。。。。\n"))
	if err != nil {
		fmt.Println("post ping err",err)
		return
	}
}


func main() {
	//创建一个server句柄
	s := znet.NewServer("[zinx V0.2]")
	//此处打印s.Name失败，原因是返回一个实现ziface.IServer接口的对象，不是返回一个结构体
	fmt.Printf("s = %#v\n",s)
	//添加自定义路由
	s.AddRouter(&PingRouter{})

	//启动
	s.Serve()

}
