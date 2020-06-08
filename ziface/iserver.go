package ziface

// 定义服务端接口
type IServer interface {
	Start() // 启动服务端方法
	Stop()  // 停止服务端方法
	Serve() // 开启业务服务方法
	AddRouter() // 路由功能: 给当前服务注册一个路由业务方法,供客户端连接处理使用
}
