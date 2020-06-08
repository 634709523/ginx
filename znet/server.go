package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/ziface"
)

// iServer的接口实现 定义一个Server服务
type Server struct {
	Name      string // 服务器名称
	IP        string // IP地址
	IPVersion string // tcp4 or other
	Port      int    // 端口
}

func (s *Server) Start() {
	var (
		addr     *net.TCPAddr
		err      error
		listener *net.TCPListener
		conn     *net.TCPConn
		cnt int 
	)
	fmt.Printf("[START] Server listener at IP %s,Port %d,is starting\n", s.IP, s.Port)
	// 开启一个goroutine去做服务端Listener服务
	go func() {
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
		// 启动Server网络连接业务
		for {
			// 3.1 阻塞等待客户端的连接
			if conn, err = listener.AcceptTCP(); err != nil {
				fmt.Printf("Accept Failed,err:%s\n", err)
				continue
			}
			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			//3.3 TODO Server.Start() 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的

			// 做一个回复512字节的回显服务
			go func() {
				// 不断的循环从客户端获取数据
				for {
					buf := make([]byte, 512)
					if cnt, err = conn.Read(buf); err != nil {
						fmt.Printf("recv buf err:%s\n", err.Error())
						continue
					}
					// 回显
					if _, err = conn.Write(buf[:cnt]); err != nil {
						fmt.Printf("send back buf err:%s\n", err.Error())
						continue
					}
				}
			}()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

/*
  创建一个服务器句柄
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}
	return s
}
