package test

import (
	"net"
	"util"
	"io/ioutil"
	"time"
)

func main() {
	server()
	time.Sleep(5* time.Second)
	client()

	time.Sleep(time.Hour)
}

/*服务端*/
func server() {
	addr := "0.0.0.0:8081"
	var listener, err = net.Listen("tcp", addr)
	util.LogError(err, "服务端启动失败", err)
	//defer listener.Close()

	go serverAccept(listener)

}

/*服务端等待连接*/
func serverAccept(listener net.Listener) {
	for {
		// 接受连接
		var conn, err = listener.Accept()
		util.LogError(err, "服务端接收连接失败:", err)
		println("主机:", conn.RemoteAddr(), "，连接成功")
		// 处理该连接
		go serverProcessConn(conn)
	}
}

/*处理服务端连接*/
func serverProcessConn(conn net.Conn) {
	for {
		// 读取消息
		var byteMsg, err = ioutil.ReadAll(conn)
		if !util.LogError(err, "服务端读取消息失败,当前主机:", conn.RemoteAddr(), "，异常:", err) {
			// 失败退出
			break
		}
		println("接收到客户端:", conn.RemoteAddr(), "的消息：", string(byteMsg))
		conn.Write([]byte("1"))
	}
}

/*客户端*/
func client() {
	var addr = "127.0.0.1:8081"
	var conn, err = net.Dial("tcp", addr)
	util.LogError(err, "连接到主机:", addr, ",发生异常:", err)

	// 处理读取
	go clientProcessConn(conn)

	for i := 0; i < 20; i++ {
		conn.Write([]byte("你好"))
		time.Sleep(time.Second)
	}
	conn.Close()
}

/*处理客户端连接*/
func clientProcessConn(conn net.Conn) {
	for {
		// 读取消息
		var byteMsg, err = ioutil.ReadAll(conn)
		if !util.LogError(err, "客户端读取消息异常:", err) {
			// 失败退出
			break
		}
		println("接收到服务端:", conn.RemoteAddr(), "的消息：", string(byteMsg))
	}

}
