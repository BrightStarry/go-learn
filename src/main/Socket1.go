package main

import (
	"net"
	"fmt"
	"reflect"
	"io"
	"bufio"
	"util"
)

func main() {
	addr := "www.baidu.com:80" //主机名
	conn, err := net.Dial("tcp", addr)
	// 如果发生异常
	util.LogError(err,"连接到主机:", addr,",发生异常:",err)

	fmt.Println("访问公网ip：",conn.RemoteAddr().String())
	fmt.Println("客户端的地址和端口是:",conn.LocalAddr())
	fmt.Println("“conn.LocalAddr()”所对应的数据类型是：",reflect.TypeOf(conn.LocalAddr()))
	fmt.Println("“conn.RemoteAddr().String()”所对应的数据类型是：",reflect.TypeOf(conn.RemoteAddr().String()))

	// 向服务端发送数据, http get请求报文，\r\n\r\n是 请求头和请求主体的 分割
	n,err := conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
	util.LogError(err,"发送http报文发生异常:",err)
	println("发送的数据大小是",n)


	// 这种读取方式也可以 循环读取
	// 定义切片，长度为1024
	//buf := make([]byte,10240)
	//
	//// 读取消息到buf
	//n,err = conn.Read(buf)
	//if err != nil && err != io.EOF {//io.EOF在网络编程中表示对端把链接关闭了。
	//	log.Fatalln("读取http报文发生异常:",err)
	//}
	//println("读取的数据大小是",n)

	// 或者直接全部读取 os.Stdout表示标准输出，也就是打印到控制台
	//io.Copy(os.Stdout,conn)

	// 按行读取,最后，全部读取完成后，会被阻塞在ReadString方法中
	reader := bufio.NewReader(conn)
	for {
		line,err := reader.ReadString('\n')
		if err == io.EOF {
			conn.Close()
			break
		}
		fmt.Print(line)
	}

	// 打印出读取到的内容
	//println(string(buf[:n]))


	// 关闭连接
	defer conn.Close()




}