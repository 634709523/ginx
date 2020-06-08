package ziface

// 定义服务端接口
type IServer interface {
	Start() // 启动服务端方法
	Stop()  // 停止服务端方法
	Serve() // 开启业务服务方法
}
