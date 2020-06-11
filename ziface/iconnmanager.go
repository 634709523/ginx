package ziface

/*
	连接管理抽象
*/
type IConnManager interface {
	Add(conn IConnection)                   // 添加连接
	Remove(coon IConnection)                // 删除连接
	Get(connID uint32) (IConnection, error) // 根据连接ID获取连接
	Len() int                               // 获取当前连接
	ClearConn()                             // 删除并停止所有连接
}
