package main

import (
	"net"
	"fmt"
	"net/url"
	"strings"
	"io"
	"bytes"
	"time"
)

func main() {
	startServer1("10000")
}

/*启动服务器*/
func startServer1(port string) {
	listener,err := net.Listen("tcp",":" + port)
	if err != nil {
		println("服务端启动失败：",err)
	}

	println("1服务端启动成功，当前端口:",port)

	for{
		client,err := listener.Accept()
		if err != nil {
			println("接收客户端失败:",err)
		}

		println("客户端连接成功：",client.RemoteAddr().String())

		go handlerClientRequest(client)
	}

}

/*处理客户端请求*/
func handlerClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	//defer client.Close()

	client.SetDeadline(time.Now().Add(time.Duration(30) * time.Second))

	// 尝试读取1024个字节
	var buf [1024]byte
	// n是真正读取的个数
	n,err :=client.Read(buf[:])
	if err != nil {
		println("处理客户端请求失败:",err)
		return
	}

	var method, host, address string
	// 取出请求报文中的第一行，将第一个作为 http请求方法(GET)， 第二个作为 主机名(http://www.flysnow.org/)
	// GET http://www.flysnow.org/ HTTP/1.1

	// 如果是https请求，这样的报文 CONNECT 106.11.61.114:443 HTTP/1.1
	// 在进行url.Parse(host)，此时host是106.11.61.114:443会报异常

	// 此处需要追加判断，因为可能请求报文第一行长度超过1024字节，则无法找到 '\r' ，则会下标越界
	//test1 := string(buf[:])
	//log.Println("测试数据:",test1)

	i := bytes.IndexByte(buf[:],'\r')
	if i== -1{
		println("请求报文首行字节数过长")
		return
	}
	fmt.Sscanf(string(buf[:i]),"%s%s",&method,&host)



	// 拼接出 请求的主机名
	// https
	if host[len(host)-4:] == ":443" {
		//address = hostPortUrl.Scheme + ":443"
		address = host
	}else {
		//http

		hostPortUrl,err :=url.Parse(host)
		if err != nil{
			println("解析host异常:",err)
			return
		}

		if strings.Index(hostPortUrl.Host,":") == -1{
			// 如果不带端口号，则默认80
			address = hostPortUrl.Host + ":80"
		}else{
			address = hostPortUrl.Host
		}
	}

	// 请求主机
	server,err :=net.Dial("tcp",address)
	if err != nil {
		println("连接到对应主机:",address,",异常:",err)
		return
	}
	server.SetDeadline(time.Now().Add(time.Duration(30) * time.Second))

	// https
	if method == "CONNECT" {
		// https需要先响应客户端，表示已经建立连接，然后在进行数据传输
		fmt.Fprint(client,"HTTP/1.1 200 Connection Established\r\n\r\n")
	} else{
		// http
		// 向目标主机发送  buf数组中读取到的所有字节
		server.Write(buf[:n])
	}
	// 将从client读取到的数据 写入 到 目标服务器
	go io.Copy(server,client)
	// 将从目标服务器读取的数据 写入到 客户端
	 io.Copy(client,server)

}


