package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn    // 当前连接的socket TCP套接字
	ConnID       uint32          // 当前连接的ID
	isClosed     bool            // 当前连接的关闭状态
	// handAPI      ziface.HandFunc // 该连接的处理方法api
	ExitBuffChan chan bool       // 告知该连接已经退出/ 停止的Channel
	Router       ziface.IRouter  // 该连接的处理方法router
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter)*Connection{
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		Router: router,
	}
	return c
}

// 处理conn读数据的goroutine
func (c *Connection) StartReader() {
	var (
		err error
		dp *DataPack
		headData []byte
		msg ziface.IMessage
		data []byte
	)
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit!")
	defer c.Stop()
	for {
		// 创建解包对象
		dp = NewDataPack()
		// 读取客户端的msg head
		headData = make([]byte, dp.GetHeadLen())
		if _,err = io.ReadFull(c.GetTCPConnection(),headData);err != nil{
			fmt.Println("read head failed,err:",err)
			break
		}
		// 拆包，得到msg id和data
		if msg, err = dp.UnPack(headData);err != nil{
			fmt.Println("unpack failed,err:",err)
			break
		}
		if msg.GetDataLen() > 0 {
			data = make([]byte,msg.GetDataLen())
			if _,err = io.ReadFull(c.GetTCPConnection(),data);err != nil{
				fmt.Println("get data failed,err:",err)
				break
			}
		}
		msg.SetData(data)
		// 得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		// 从路由Routers中找到注册绑定Conn的对应Handle
		go func(request ziface.IRequest){
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
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

// 发送消息
func (c *Connection) SendMsg(msgId uint32, data []byte) error{
	var (
		dp *DataPack
		msg []byte
		err error
	)
	// 判断连接是否关闭
	if c.isClosed{
		return errors.New("Connection Closed when send msg")
	}
	// 将data封包，并且发送
	dp = NewDataPack()
	if msg,err = dp.Pack(NewMsgPackage(msgId,data));err != nil{
		fmt.Println("Pack error msg id = ",msgId)
		return err
	}
	// 返回客户端
	if _,err = c.Conn.Write(msg);err != nil{
		fmt.Println("write msg id ",msgId,"error")
		return err
	}

	return nil

}