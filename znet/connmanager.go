package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接信息
	connLock sync.RWMutex // 读写连接的读写锁
}

func NewConnManager()*ConnManager{
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
		connLock:    sync.RWMutex{},
	}
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	// 判断连接信息是否已经存在
	connId := conn.GetConnID()
	if _,ok := c.connections[connId];ok{
		fmt.Println("Add Connection Failed,connID %d is already existed",connId)
		return
	}
	// 添加连接信息
	c.connections[connId] = conn
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	// 判断连接信息是否已经存在
	connId := conn.GetConnID()
	if _,ok := c.connections[connId];!ok{
		fmt.Println("Remove Connection Failed,CoonID is not %d exsited",connId)
		return
	}
	// 删除连接信息
	delete(c.connections,connId)
}

func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	// 获取连接信息
	var (
		conn ziface.IConnection
		ok bool
	)
	if conn, ok = c.connections[connID]; !ok {
		return nil, errors.New(fmt.Sprintf("get connection failed,connection is not exsits,connID:%d", connID))
	}
	return conn, nil
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	for connID,conn := range c.connections{
		// 停止
		conn.Stop()
		// 删除
		delete(c.connections, connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", c.Len())
}
