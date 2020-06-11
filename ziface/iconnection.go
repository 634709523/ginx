package ziface

import (
	"net"
)

type IConnection interface {
	Start()                                  // 启动连接
	Stop()                                   // 停止连接
	GetTCPConnection() *net.TCPConn          // 从当前连接获取原始的socket TCPConn
	GetConnID() uint32                       // 获取当前连接ID
	RemoteAddr() net.Addr                    // 获取远程地址
	SendMsg(msgId uint32, data []byte) error //直接将Message数据发送数据给远程的TCP客户端(无缓冲)
	//直接将Message数据发送给远程的TCP客户端(有缓冲)
	SendBuffMsg(msgId uint32, data []byte) error   //添加带缓冲发送消息接口
	ConvertDataToMsg(msgId uint32, data []byte)([]byte,error) // 将数据转换成消息
	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string)(interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}

//定义一个统一处理链接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
