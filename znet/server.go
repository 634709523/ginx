package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

// iServer的接口实现 定义一个Server服务
type Server struct {
	Name      string // 服务器名称
	IP        string // IP地址
	IPVersion string // tcp4 or other
	Port      int    // 端口
    // 当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务
	msgHandler ziface.IMsgHandler
	//当前Server的连接管理器
	connManager ziface.IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

/*
  创建一个服务器句柄
*/
func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.Port,
		msgHandler: NewMsgHandler(),
		connManager:NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	var (
		addr     *net.TCPAddr
		err      error
		listener *net.TCPListener
		conn     *net.TCPConn
	)
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)	// 开启一个goroutine去做服务端Listener服务
	go func() {
		// 0. 启动worker工作池机制
		s.msgHandler.StartWorkerPool()
		// 1. 获取一个TCP Addr
		if addr, err = net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port)); err != nil {
			fmt.Printf("resolve tcp addr err:%s\n", err.Error())
			return
		}
		// 2. 监听服务器地址
		if listener, err = net.ListenTCP(s.IPVersion, addr); err != nil {
			fmt.Printf("listen %s,err:%s", s.IPVersion, err.Error())
		}
		// 3. 已经监听成功
		fmt.Printf("start Zinx Server %s success,now listen %s:%d\n", s.Name, s.IP, s.Port)

		// TODO server.go 这里应该有一个自动生成ID的方法
		var cid uint32
		cid = 0
		// 启动Server网络连接业务
		for {
			// 3.1 阻塞等待客户端的连接
			if conn, err = listener.AcceptTCP(); err != nil {
				fmt.Printf("Accept Failed,err:%s\n", err)
				continue
			}
			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.connManager.Len() >= utils.GlobalObject.MaxConn{
				conn.Close()
				continue
			}
			//3.3 处理该新连接请求的业务方法，此时应该有handler和conn是绑定
			dealConn := NewConnection(s,conn, cid, s.msgHandler)
			cid++
			// 3.4 启动当前连接的处理业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	// 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.connManager.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}



//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server)AddRouter(msgId uint32,router ziface.IRouter) {
	s.msgHandler.AddRoute(msgId,router)
    fmt.Println("Add Router succ! " )
}

func (s *Server)GetConnMgr() ziface.IConnManager {
	return s.connManager
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接停止时Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil{
		fmt.Println("---> CallOnConnStart")
		s.OnConnStart(conn)
	}
}
//调用连接CallOnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil{
		fmt.Println("---> CallOnConnStop")
		s.OnConnStop(conn)
	}
}