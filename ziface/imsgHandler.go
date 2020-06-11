package ziface

/*
	消息管理抽象层
 */

type IMsgHandler interface {
	DoMsgHandler(request IRequest) // 以非阻塞方式处理消息
	AddRoute(msgId uint32,router IRouter) // 为消息添加具体处理方法
	StartWorkerPool()                       //启动worker工作池
	SendMsgToTaskQueue(request IRequest)    //将消息交给TaskQueue,由worker进行处理
	StartOneWorker(workerID int, taskQueue chan IRequest)
}