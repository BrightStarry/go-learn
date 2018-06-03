package test

import (
	"net"
	"util"
	"time"
	"io"
	"os"
)

func main() {
	startServer()
}

/*启动Socket服务端*/
func startServer() {
	addr := "0.0.0.0:8081" // 监听本地端口，也可以 ":8081" 这么写
	listener,err := net.Listen("tcp", addr)
	util.LogError(err,"启动服务异常:",err)
	defer listener.Close() //关闭监听的端口

	for{
		// 接受连接
		conn,err := listener.Accept()
		util.LogError(err,"服务端接受连接异常:",err)
		conn.Write([]byte("服务端收到消息\n"))
		// 读取客户端
		io.Copy(os.Stdout,conn)
		// 等待20s，再关闭连接
		time.Sleep(20*time.Second)
		conn.Close()
	}

	println("服务端退出")
}
