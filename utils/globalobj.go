package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"zinx/ziface"
)

/*
	存储一切与zinx框架的全局参数，供其他模块使用
	参数可以通过zinx.json来配置
 */

type GlobalObj  struct {
	TcpServer ziface.IServer // 当前Zinx的全局server对象
	Host string `json:"Host"`// 当前服务器主机IP
	Port int `json:"Port"`//当前服务器主机监听端口号
	Name string `json:"Name"` // 当前服务器名称
	Version string // 当前Zinx版本号

	MaxPacketSize uint32 // 数据包的最大值
	MaxConn int `json:"MaxConn"`// 当前服务器主机允许的最大连接个数
}

// 定义一个全局的对象
var (
	GlobalObject  *GlobalObj
)

// 读取配置文件
func (g *GlobalObj)Reload(){
	var (
		data []byte
		err error
	)
	if data, err = ioutil.ReadFile("conf/zinx.json");err != nil{
		goto ERR
	}
	if err = json.Unmarshal(data,&GlobalObject);err != nil{
		goto  ERR
	}
	return
ERR:
	log.Println(err)
	panic(err)
}

// 加载全局对象
func init(){
	GlobalObject  = &GlobalObj{
		Name: "ZinxServerApp",
		Version:"v0.6",
		Host: "0.0.0.0",
		Port:7777,
		MaxPacketSize: 4000,
		MaxConn: 12000,
	}
	// 从配置文件中加载用户配置参数
	GlobalObject.Reload()
}