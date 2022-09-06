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
//集成：server，connection
//zinx0.7，新增读写分离模型
//添加reader和writer之间通信的channel
//添加write goroutine
//reader由之间发送给客户端，改成发送给通信的channel
//启动reader，与writer一起工作
//zinx0.8消息队列和任务池
//创建消息队列，msghandle模块
//创建多任务work工作池并启动,根据workerpoolsize的数量去创建worker，每个worker都应该阻塞等待与当前worker对应的channel的消息，一旦有消息到来，worker应该处理当前消息对应的业务，调用domsghandler
//将之前发送消息全部改成发送消息给消息队列和work池处理,定义方法，将消息发送给队列queue的方法：1保证每个worker所受到的req任务是均衡的(轮询)，让那个worker去处理，只需要将request请求发送给对应的taskqueue，2将消息发送到对应额channel
//将消息队列机制集成到zinx框架中，1开启并调用消息队列及worker工作池(保证workerpool只有一个，应该在创建server时创建工作池)，2将客户端处理的消息，发送给当前的worker工作池来处理(将已经处理完拆包，得到req请求，交给工作池处理)


type PingRouter struct {
	znet.BaseRouter
}
type HelloRouter struct {
	znet.BaseRouter
}


//处理业务
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call back Handle.....")
	//先读取客户端的数据，在会写ping
	fmt.Println("recv from client: msgid=",request.GetMsgId(),
		            ",msg=",string(request.GetData()))
	err := request.GetConnection().Send(0,[]byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}
func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call back Handle.....")
	//先读取客户端的数据，在会写ping
	fmt.Println("\"recv from client: msgid=",request.GetMsgId(),
		",msg=",string(request.GetData()))
	err := request.GetConnection().Send(1,[]byte("hello zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//创建一个server句柄
	s := znet.NewServer()
	//此处打印s.Name失败，原因是返回一个实现ziface.IServer接口的对象，不是返回一个结构体
	fmt.Printf("s = %#v\n",s)
	//添加自定义路由
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloRouter{})

	//启动
	s.Serve()

}
