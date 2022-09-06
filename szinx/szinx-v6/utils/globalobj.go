package utils

import (
	"encoding/json"
	"io/ioutil"
	"szinx/ziface"
)
/*
存储一切有关zinx框架的全局参数，供其他模块使用
 */

type GlobalObj struct {
	TcpServer ziface.IServer    //当前全局的server对象
	Host string
	TcpPort int
	Name string
	Version string
	Maxconn int         //当前允许的最大连接数
	MaxPackageSize uint32       //当前框架数据包的最大值
}
var G *GlobalObj

func (this *GlobalObj) Reload(){
	data,err := ioutil.ReadFile("D:\\goproject\\szinx\\szinx\\conf\\zinx.json")
	if err != nil {
		panic(err)
	}

	//将json文件数据解析到struct中,G是一个结构体变量的指针
	err = json.Unmarshal(data,&G)  //将数据写入G中，data是一个字节类型
	if err != nil {
		panic(err)
	}
}



func init(){
	G = &GlobalObj{
		Name: "App",
		Version: "V0.5",
		TcpPort: 7777,
		Host: "127.0.0.1",
		Maxconn: 1000,
		MaxPackageSize: 4096,
	}
	G.Reload()
}





















