package znet

import (
	"fmt"
	"net"
	"time"
	"testing"
)

// 模拟客户端
func ClientTest(){
	fmt.Println("Client Test ... Start")
	var (
		conn net.Conn
		err error
		cnt int 
	)
	// 休息3秒，给建立服务端一些缓冲时间
	time.Sleep(time.Second * 3)
	if conn,err = net.Dial("tcp","127.0.0.1:7777");err != nil{
		fmt.Printf("dial %s failed,err:%s\n","127.0.0.1:7777",err)
		return 
	}
	for {
		if _, err = conn.Write([]byte("hello ZINX"));err != nil{
			fmt.Println("send message failed,err:",err)
			continue
		}
		recv := make([]byte,512)
		if cnt, err = conn.Read(recv);err != nil{
			fmt.Println("recv message failed,err:",err)
			continue
		}
		fmt.Printf(" server call back : %s, cnt = %d\n", recv,  cnt)
		time.Sleep(1 *time.Second)
	}
}

// 测试服务
func ServerTest(t *testing.T){
    s := NewServer("[zinx V0.1]")
	go ClientTest()
	s.Serve()
}