package ziface

/*
	消息管理抽象层
 */

type IMsgHandler interface {
	DoMsgHandler(request IRequest) // 以非阻塞方式处理消息
	AddRoute(msgId uint32,router IRouter) // 为消息添加具体处理方法
}