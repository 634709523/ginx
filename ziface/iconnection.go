package ziface

import (
	"net"
)

type IConnection interface {
	Start()                         // 启动连接
	Stop()                          // 停止连接
	GetTCPConnection() *net.TCPConn // 从当前连接获取原始的socket TCPConn
	GetConnID() uint32              // 获取当前连接ID
	RemoteAddr() net.Addr           // 获取远程地址
}

//定义一个统一处理链接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
