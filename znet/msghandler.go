package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	Apis map[uint32] ziface.IRouter
	WorkerPoolSize uint64 // 业务工作worker池的数量
	TaskQueue []chan ziface.IRequest // worker负责取任务的消息队列
}


// 创建消息处理器
func NewMsgHandler()*MsgHandler{
	return &MsgHandler{
		Apis:make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		// 一个worker对应一个queue
		TaskQueue:make([]chan ziface.IRequest,utils.GlobalObject.WorkerPoolSize),
	}
}

// 处理消息处理器
func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	var (
		router ziface.IRouter
		ok bool
	)
	// 从Api中获取路由
	if router,ok = m.Apis[request.GetMsgID()];!ok{
		fmt.Println("api msgId = ",request.GetMsgID(),"is not Found")
		return
	}
	// 执行处理函数
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

// 添加路由
func (m *MsgHandler) AddRoute(msgId uint32, router ziface.IRouter) {
	var (
		ok bool
	)
	// 已经注册过处理器则报错
	if _,ok = m.Apis[msgId];ok{
		panic("repeated api, msgId = " + strconv.Itoa(int(msgId)))
	}
	m.Apis[msgId] = router
	fmt.Println("Add api msgId =",msgId)
}

// 启动一个worker处理消息
func (m *MsgHandler)StartOneWorker(workerID int, taskQueue chan ziface.IRequest){
	fmt.Println("worker ID = ",workerID,"is started")
	// 不断的等待消息队列中的消息
	for {
		select {
		case request := <- taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// 启动worker工作池
func (m *MsgHandler) StartWorkerPool() {
	for i :=1 ;i<=int(m.WorkerPoolSize);i++{
		// 给当前worker对应的消息队列开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest,utils.GlobalObject.WorkerPoolSize)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go m.StartOneWorker(i,m.TaskQueue[i])
	}
}

// 将消息交给TaskQueue,由worker进行处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 根据ConnID 来分配当前的连接应该由哪个worker负责处理
	// 轮询的平均分配去轮询
	// 得到需要处理此条连接的workerid
	workerID := int(request.GetConnection().GetConnID()) % int(m.WorkerPoolSize)
	fmt.Println("Add ConnID=",request.GetConnection().GetConnID(),"request msgID=",workerID)
	m.TaskQueue[workerID] <- request
}
