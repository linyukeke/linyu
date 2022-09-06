package znet

import (
	"szinx/ziface"
)

type BaseRouter struct {
}
//此处Baseroute的方法为空，是因此此处继承了抽象层的所有方法，业务可以根据所用，灵活重写某个方法，不需要的方法不用关心
//业务之前
func (this BaseRouter) PreHandle(request ziface.IRequest) {}
//处理业务
func (this BaseRouter) Handle(request ziface.IRequest) {}
//业务之后
func (this BaseRouter) PostHandle(request ziface.IRequest) {}


