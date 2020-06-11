package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	TcpServer ziface.IServer // 当前conn属于哪个server
	Conn     *net.TCPConn // 当前连接的socket TCP套接字
	ConnID   uint32       // 当前连接的ID
	isClosed bool         // 当前连接的关闭状态
	ExitBuffChan chan bool          // 告知该连接已经退出/ 停止的Channel
	MsgHandler   ziface.IMsgHandler // 消息管理MsgId和对应处理方法的消息管理模块
	msgChan      chan []byte        // 无缓冲通道，用于读、写两个goroutine之间的消息通信
	msgBufChan chan []byte // 有缓冲通道，用于读、写两个goroutine之间的消息通信
	property map[string]interface{} // 连接属性
	propertyLock sync.RWMutex
}



func NewConnection(server ziface.IServer,conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		MsgHandler:   msgHandler,
		msgChan: make(chan []byte),
		msgBufChan: make(chan []byte,utils.GlobalObject.MaxMsgChanLen),
		property: make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}



// 处理conn读数据的goroutine
func (c *Connection) StartReader() {
	var (
		err      error
		dp       *DataPack
		headData []byte
		msg      ziface.IMessage
		data     []byte
	)
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit!")
	defer c.Stop()
	for {
		// 创建解包对象
		dp = NewDataPack()
		// 读取客户端的msg head
		headData = make([]byte, dp.GetHeadLen())
		if _, err = io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read head failed,err:", err)
			break
		}
		// 拆包，得到msg id和data
		if msg, err = dp.UnPack(headData); err != nil {
			fmt.Println("unpack failed,err:", err)
			break
		}
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err = io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("get data failed,err:", err)
				break
			}
		}
		msg.SetData(data)
		// 得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0{
			// 已经启动工作池机制，将消息交给worker处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		}else{
			// 从路由Routers中找到注册绑定Conn的对应Handle
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

func (c *Connection) Start() {
	// 开启处理该连接到客户端数据之后的请求业务
	go c.StartReader()
	// 开启用于写回客户端数据流程的goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
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
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)


	// 释放连接
	c.Conn.Close()
	// 通知从缓冲队列读取数据的业务，该连接已经关闭
	c.ExitBuffChan <- true
	// 删除连接信息
	c.TcpServer.GetConnMgr().Remove(c)
	// 关闭该连接的通道
	close(c.ExitBuffChan)
	close(c.msgBufChan)
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

// 将消息发送到有缓冲的消息通道
func(c *Connection)SendBuffMsg(msgId uint32, data []byte) error{
	var (
		msg []byte
		err error
	)
	if msg, err = c.ConvertDataToMsg(msgId,data);err != nil {
		return err
	}
	// 返回客户端
	c.msgBufChan <- msg
	return nil
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	var (
		msg []byte
		err error
	)
	if msg, err = c.ConvertDataToMsg(msgId,data);err != nil {
		return err
	}
	// 返回客户端
	c.msgChan <- msg
	return nil
}

func (c *Connection) ConvertDataToMsg(msgId uint32, data []byte)([]byte,error){
	var (
		dp  *DataPack
		msg []byte
		err error
	)
	// 判断连接是否关闭
	if c.isClosed {
		return nil,errors.New("Connection Closed when send msg")
	}
	// 将data封包，并且发送
	dp = NewDataPack()
	if msg, err = dp.Pack(NewMsgPackage(msgId, data)); err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return nil,err
	}
	return msg,nil
}

/*
   写消息Goroutine， 用户将数据发送给客户端
*/
func (c *Connection)StartWriter(){
	defer fmt.Println(c.RemoteAddr().String(),"conn Writer exit")
	for {
		select {
		case data :=<- c.msgChan:
			// 将数据写回客户端
			if _,err := c.Conn.Write(data);err != nil{
				fmt.Println("Send Data error, ",err,"Conn Writer error")
				return
			}
		// 针对有缓冲channel的数据处理
		case data,ok := <- c.msgBufChan:
			if ok{
				if _,err := c.Conn.Write(data);err != nil{
					fmt.Println("Send Buff Data error, ",err,"Conn Writer error")
					return
				}
			}else{
				break
				fmt.Println("msgBuffChan is Closed")
			}
		case <- c.ExitBuffChan:
			// 连接已经关闭
			return
		}
	}
}


func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (value interface{},err error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	var ok bool
	if value,ok = c.property[key];!ok{
		// 没有这个属性
		return nil,errors.New("Property not have the key,ker:"+key)
	}
	return value,nil
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	var ok bool
	if _,ok = c.property[key];ok{
		delete(c.property, key)
	}
	return
}