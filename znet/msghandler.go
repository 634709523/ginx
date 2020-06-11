package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

type MsgHandler struct {
	Apis map[uint32] ziface.IRouter
}

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

// 创建消息处理器
func NewMsgHandler()*MsgHandler{
	return &MsgHandler{Apis:make(map[uint32]ziface.IRouter)}
}


