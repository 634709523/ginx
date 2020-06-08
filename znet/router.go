package znet

import (
	"zinx/ziface"
)

type BaseRouter struct {}


func (br *BaseRouter)PreHandle(request IRequest){}
func (br *BaseRouter)Handle(request IRequest) {}
func (br *BaseRouter)PostHandle(request IRequest){}