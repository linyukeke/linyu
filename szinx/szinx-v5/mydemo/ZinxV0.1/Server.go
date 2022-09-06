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
//v0.4全局配置
//创建一个全局配置模块utils/globalobj.go,init方法读取conf下的配置文件到globalobj对象中，//替换原代码中的参数
//v0.5message消息模块配置
//定义一个消息的结构，属性：消息id，消息长度，消息内容；方法：get，set
//定义datapack，针对message进行tlv格式的封包和解包(先得到长度和类型，再根据长度读取内容)，单元测试(未成功)
//将消息封装机制集成到zinx框架中：
//将message添加到request属性中，
//修改connection读取的机制，将之前的单纯的读取byte改成拆包形式的读取按照tlv形式读取，给链接一个发包机制，将消息打包发送
//server的测试用例
//zinx0.6，消息管理模块属性，支持多路由业务api调度管理
//属性：消息id，对应的route对应关系
//方法：根据msgid来索引调度路由方法，添加路由方法到map集合中



type PingRouter struct {
	znet.BaseRouter
}

//处理业务
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call back Handle.....")
	//先读取客户端的数据，在会写ping
	fmt.Println("msgid=",request.GetMsgId(),
		            ",msg=",request.GetData())
	err := request.GetConnection().Send(1,[]byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("ping err",err)
	}
}


func main() {
	//创建一个server句柄
	s := znet.NewServer()
	//此处打印s.Name失败，原因是返回一个实现ziface.IServer接口的对象，不是返回一个结构体
	fmt.Printf("s = %#v\n",s)
	//添加自定义路由
	s.AddRouter(&PingRouter{})

	//启动
	s.Serve()

}
