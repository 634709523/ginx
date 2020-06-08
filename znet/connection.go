package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn    // 当前连接的socket TCP套接字
	ConnID       uint32          // 当前连接的ID
	isClosed     bool            // 当前连接的关闭状态
	handAPI      ziface.HandFunc // 该连接的处理方法api
	ExitBuffChan chan bool       // 告知该连接已经退出/ 停止的Channel
}

func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		handAPI:      callback_api,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

// 处理conn读数据的goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit!")
	defer c.Stop()
	for {
		// 读取最大的数据到buf中
		var cnt int 
		var err error
		buf := make([]byte, 512)
		if cnt, err = c.Conn.Read(buf); err != nil {
			fmt.Printf("recv data failed,err:%s\n", err)
			continue
		}
		// 调用当前连接业务
		if err = c.handAPI(c.Conn, buf, cnt); err != nil {
			fmt.Printf("connID %d handle is error\n", c.ConnID)
			c.ExitBuffChan <- true
			return
		}
	}
}

func (c *Connection) Start() {
	// 开启处理该连接到客户端数据之后的请求业务
	go c.StartReader()
	for {
		select {
		// 得到退出消息，不再阻塞
		case <-c.ExitBuffChan:
			return
		}
	}
}

// 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	if c.isClosed != true {
		c.isClosed = true
	}
	// 释放连接
	c.Conn.Close()
	// 通知从缓冲队列读取数据的业务，该连接已经关闭
	c.ExitBuffChan <- true
	// 关闭该连接的通道
	close(c.ExitBuffChan)
}

// 获取当前tcp连接
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
