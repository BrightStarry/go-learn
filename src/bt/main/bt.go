package main

import (
	"net"
	"errors"
	"fmt"
	"encoding/hex"
	"crypto/rand"
	"time"
	"bytes"
	"io"
)

/*仅测试 从peer获取metadata信息*/

//获取种子元信息时,第一条握手信息的前缀, 28位byte. 第2-20位,是ASCII码的BitTorrent protocol,
// 第一位19,是固定的,表示这个字符串的长度.后面八位是BT协议的版本.可以全为0,某些软件对协议进行了扩展,协议号不全为0,不必理会.
var handshakePrefix = []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114,
	111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 16, 0, 1}

const(
	// 表示是一个扩展消息
	EXTENDED = 20
	// 表示握手的一位
	HANDSHAKE = 0
)

func main() {
	//addr := ""
	//conn := connect(addr)

	//bytes := hexInfoHashToByte("eb8abb5d2b4711b4d545b9d0ebb05f22b63f5ca3")
	//fmt.Println(bytes)



}

/*连接到peer*/
func connect(addr string) net.Conn{
	conn,err :=net.Dial("tcp",addr)
	if err != nil {
		panic(errors.New(fmt.Sprintln("连接到",addr ,",异常:" + err.Error())))
	}
	return conn
}

/*发送握手消息*/
func sendHandshakeMessage(conn *net.TCPConn, infoHash string,peerId []byte) error {
	data := make([]byte,68)

	copy(data[:28],handshakePrefix)
	copy(data[28:48],hexInfoHashToByte(infoHash))
	copy(data[48:],peerId)

	// 设置超时时间10s
	conn.SetWriteDeadline(time.Now().Add(10*time.Second))

	_,err := conn.Write(data)
	return err
}

/*处理握手响应*/
func onHandshake(data []byte) (err error) {
	if !bytes.Equal(handshakePrefix[:20],data[:20]) && data[25]&0x10 != 0 {
		err = errors.New("无效握手响应")
	}
	return
}

/*发送扩展握手协议,获取ut_metadata  和 metadata_size*/
func sendExtHandshake(conn *net.TCPConn) {
	//data := append(
	//	[]byte{EXTENDED, HANDSHAKE},
	//
	//
	//)
}

/*从连接中读取指定长度的消息*/
func read(conn *net.TCPConn,len int, data *bytes.Buffer) error {
	conn.SetReadDeadline(time.Now().Add(15 * time.Second))

	n,err := io.CopyN(data,conn,int64(len))
	// 如果读取有异常，或者字节数不等
	if err != nil ||  n != int64(len) {
		return errors.New("读取异常" + err.Error())
	}
	return nil
}

/*从连接中读取消息*/
func readMessage(conn *net.TCPConn,data *bytes.Buffer) (len int,err error){


}

/*16进制infoHash转[]byte*/
func hexInfoHashToByte(infoHash string) []byte {
	result,_ := hex.DecodeString(infoHash)
	return result
}

/*生成指定长度的随机字符*/
func randomString(size int) string {
	return string(randomBytes(size))
}

/*生成指定长度的随机字节*/
func randomBytes(size int) []byte {
	buff := make([]byte,size)
	rand.Read(buff)
	return buff
}
