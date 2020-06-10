package znet

import (
	"zinx/ziface"
)

type Request struct {
	conn ziface.IConnection // 已经和客户端建立的连接
	msg ziface.IMessage // 客户端请求的数据
}

func (r *Request)GetConnection()ziface.IConnection{
	return r.conn
}

func (r *Request)GetMsgID()uint32{
	return r.msg.GetMsgId()
}

func (r *Request)GetData()[]byte{
	return r.msg.GetData()
}