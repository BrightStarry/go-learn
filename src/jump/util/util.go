package util

import (
	"encoding/binary"
	"bytes"
	"errors"
	"net"
	"time"
	"strconv"
	"strings"
)

/*工具类*/


/*读取消息,指定默认读取长度*/
func readMessage(conn *net.TCPConn,defaultLen int) ( []byte, error) {
	// 读取字节
	buf := make([]byte,defaultLen)
	n,err := conn.Read(buf)
	if err != nil {
		return nil,err
	}
	if n <= 0 {
		err = errors.New("读取到的字节数小于等于0")
		return nil,err
	}
	return buf[:n],nil
}

/*读取消息,指定默认读取长度*/
func readMessageByReader(reader *bytes.Reader,defaultLen int) ( []byte, int, error) {
	// 读取字节
	buf := make([]byte,defaultLen)
	n,err := reader.Read(buf)
	if err != nil {
		return nil,n,err
	}
	if n <= 0 {
		err = errors.New("读取到的字节数小于等于0")
		return nil,n,err
	}
	return buf[:n],n,nil
}

/*发送响应*/
func sendResponse(conn *net.TCPConn,data interface{})(err error) {
	// 设置超时时间
	conn.SetWriteDeadline(time.Now().Add(WriteTimeout))
	// 发送对象
	err = binary.Write(conn,binary.LittleEndian,data)
	return
}

/*读取1个字节的长度字节，并读取对应长度的后续字节*/
func readByLenField(reader *bytes.Reader) (length byte,data []byte,err error) {
	length,err = reader.ReadByte()
	if err != nil {
		return
	}
	var n int
	data,n,err = readMessageByReader(reader,int(length))
	if n != int(length) {
		err = errors.New("长度不正确")
		return
	}
	return
}

/*增加用户到map*/
func addUser(user User) {
	userLock.Lock()
	AuthenticationUser[user.Ip] = user
	userLock.Unlock()
}

/*增加连接到map*/
func addTargetConn(key string,targetConn *net.TCPConn) {
	userTargetMap[key] = targetConn
}

/*[]byte转ip 字符串*/
func Bytes2Ip(data []byte) string {
	var b  bytes.Buffer
	for _,v := range data{
		b.WriteString("." + strconv.Itoa(int(v)))
	}
	return b.String()[1:]
}

/*
	地址字符串转ip([]byte) 和 port(uint16)
	类似 192.168.1.111:54698, 如果没有端口，返回0
*/
func Ip2Bytes(addr string) (ip []byte,port uint16) {
	// 转为 ["192.168.1.111","54698"] 或 （如果没有端口）["192.168.1.111"]
	arr := strings.Split(addr,":")

	ip = make([]byte,4)
	// ip
	// 转为  ["192","168","1","111"]
	ips := strings.Split(arr[0],".")
	for i,v := range ips{
		v2,_ := strconv.Atoi(v)
		ip[i] = byte(v2)
	}

	// port
	if len(arr) == 1{
		port = 0
		return
	}
	v3,_ := strconv.Atoi(arr[1])
	port = uint16(v3)
	return
}
