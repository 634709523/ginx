package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter //一定要先基础BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Printf("recv from client : msgId=%d data=%s\n", request.GetMsgID(),request.GetData())
	//回写数据
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}


//HelloZinxRouter Handle
type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router V0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer("[zinx V0.1]")
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloZinxRouter{})

	s.Serve()
}
