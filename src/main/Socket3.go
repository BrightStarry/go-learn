package main

import (
	"net"
	"util"
	"io"
	"os"
)

func main() {
	startClient()
}

/*启动客户端*/
func startClient() {
	addr := "127.0.0.1:8081"
	conn, err := net.Dial("tcp",addr)
	util.LogError(err,"连接到主机:", addr,",发生异常:",err)

	conn.Write([]byte("主机主机，我是东东八"))

	// 读取消息
	io.Copy(os.Stdout,conn)

	defer conn.Close()

}
