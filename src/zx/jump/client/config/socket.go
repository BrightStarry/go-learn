package config

import (
	"net"
	"time"
	"encoding/binary"
	"errors"
	"bytes"
	"strconv"
	"log"
	"io"
	"zx/jump/util"
	"fmt"
)

/*socket相关*/

/*
	处理tcp请求
*/
func Handler(conn *net.TCPConn) {
	defer conn.Close()
	err := HandlerPre(conn)
	if err != nil {
		log.Println("握手失败：", err)
		return
	}
	if err = HandlerConnect(conn);err != nil {
		log.Println(err)
	}
}

/*预先处理，包括 握手，认证等*/
func HandlerPre(conn *net.TCPConn) (err error) {
	/**
		处理握手
	*/
	err = handlerHandshake(conn)
	if err != nil {
		return
	}
	/**
		处理密码验证
	*/
	if Config.Method == PwdMethod {
		err = handlerPwdAuthentication(conn)
	}
	return
}

/*处理握手*/
func handlerHandshake(conn *net.TCPConn) (err error) {
	response, err := readHandshakeRequest(conn) // 读取客户端握手请求对象,生成响应对象
	if response != nil { // 响应不为空，则表示认证不支持或成功，发送响应
		err = util.SendMessage(conn, response, WriteTimeout) // 发送响应
	}
	return // 返回err
}

/*处理密码认证*/
func handlerPwdAuthentication(conn *net.TCPConn) (err error) {
	var request *PwdAuthenticationRequest // 读取
	request, err = readPwdAuthenticationRequest(conn)
	if err != nil {
		return
	}
	response := &PwdAuthenticationResponse{One, Zero} // 读取成功才对客户端进行响应

	if Config.Username != string(request.Username) || Config.Password != string(request.Password) { //校验密码
		// 失败响应
		response.Result = One
		util.SendMessage(conn, response, WriteTimeout)
		err = errors.New("用户密码错误")
		return
	}
	err = util.SendMessage(conn, response, WriteTimeout) // 成功响应
	return
}

/*读取客户端握手请求对象,生成响应对象*/
func readHandshakeRequest(conn *net.TCPConn) (response *HandshakeResponse, err error) {
	conn.SetReadDeadline(time.Now().Add(ReadTimeout)) // 设置读取超时时间

	request := new(HandshakeRequest) // 读取为对象
	err = binary.Read(conn, binary.LittleEndian, request)
	if err != nil {
		return
	}

	if request.Version != Version { // 校验版本
		err = errors.New("只支持SOCKET5")
		return
	}

	if request.MethodsLen != One { // Methods长度字段
		err = errors.New("methods长度字段不为1")
		return
	}

	response = &HandshakeResponse{Version, Config.Method} // 响应对象

	if request.Methods != Config.Method { // methods字段
		response.Method = NotSupport
		err = errors.New("不支持客户端的认证方式")
	}
	return
}

/*读取密码认证请求对象*/
func readPwdAuthenticationRequest(conn *net.TCPConn) (result *PwdAuthenticationRequest, err error) {
	var buf []byte // 读取字节
	buf, err = util.ReadMessage(conn, 512, ReadTimeout)
	if err != nil {
		return
	}

	result = new(PwdAuthenticationRequest) // 返回对象

	reader := bytes.NewReader(buf)

	if result.Pointless, err = reader.ReadByte(); err != nil { // 读取无意义标识
		return
	}

	if result.UsernameLength, result.Username, err = util.ReadByLen(reader); err != nil { // 读取用户名
		return
	}

	if result.PasswordLength, result.Password, err = util.ReadByLen(reader); err != nil { // 读取密码
		return
	}
	return
}

/*处理连接请求*/
func HandlerConnect(conn *net.TCPConn) (err error) {
	// 构造一个默认的失败响应
	response := ConnectResponse{
		Version:     Version,
		Response:    One,
		Reserve:     Zero,
		AddressType: Ipv4,
		Address:     []byte{0, 0, 0, 0},
		Port:        0,
	}

	// 读取请求
	request, err := readConnectRequest(conn)
	if err != nil {
		util.SendMessage(conn, response, WriteTimeout)
		return
	}

	var targetConn *net.TCPConn
	// 连接到目标服务器
	switch Config.Pattern {
	case Common:
		targetConn, err = util.ConnectToTarget(request.Target, ReadTimeout)
		if err != nil {
			util.SendMessage(conn, response, WriteTimeout)
			return
		}
		defer targetConn.Close()
	case CS:
		targetConn, err = util.ConnectToTarget(Config.Server, ReadTimeout)
		if err != nil {
			util.SendMessage(conn, response, WriteTimeout)
			return errors.New(fmt.Sprintln("连接到服务器失败:", err))
		}

		bTar := []byte(request.Target)
		jumpRequest := util.JumpRequest{
			PwdLen:    Config.ServerPwdLen,
			Pwd:       Config.ServerPwdByte,
			TargetLen: byte(len(bTar)),
			Target:    bTar,
		}
		util.SendMessage(targetConn, jumpRequest, WriteTimeout)
		var message []byte
		message, err = util.ReadMessage(targetConn, 1, ReadTimeout)
		if err != nil || message[0] != util.Success {
			return errors.New(fmt.Sprintln("连接服务器异常:",message))
		}
	default:
		log.Fatalln("模式错误,当前模式:", Config.Pattern)
	}

	// 连接成功
	response.Response = Zero
	// 实际ip和端口应该是可以不返回的
	//response.Address, response.Port = util.Ip2Bytes(targetConn.LocalAddr().String())
	// 发送响应
	if err = util.SendMessage(conn, response, WriteTimeout); err != nil {
		return
	}

	go io.Copy(targetConn, conn)
	io.Copy(conn, targetConn)

	return
}

/*读取连接请求*/
func readConnectRequest(conn *net.TCPConn) (request *ConnectRequest, err error) {
	// 设置超时时间
	conn.SetReadDeadline(time.Now().Add(ReadTimeout))
	// 读取字节
	var buf []byte
	buf, err = util.ReadMessage(conn, 512, ReadTimeout)
	if err != nil {
		return
	}
	// 转换为读取器
	reader := bytes.NewReader(buf)

	// 必须先创建该对象
	request = new(ConnectRequest)

	// 读取版本

	if request.Version, err = reader.ReadByte(); err != nil {
		return
	}
	// 读取命令
	if request.CMD, err = reader.ReadByte(); err != nil {
		return
	}

	// 保留字段
	if request.Reserve, err = reader.ReadByte(); err != nil {
		return
	}

	// 地址类型
	if request.AddressType, err = reader.ReadByte(); err != nil {
		return
	}

	var n int
	// 目标地址
	switch request.AddressType {
	// ipv4,4个字节
	case Ipv4:
		request.AddressBytes, n, err = util.ReadMessageByReader(reader, 4)
		if err != nil {
			return
		}
		if n != 4 {
			err = errors.New("长度不正确")
			return
		}
		// 域名，下个字节表示域名长度
	case DomainName:
		_, request.AddressBytes, err = util.ReadByLen(reader)
		if err != nil {
			return
		}
		// ipv6
	case Ipv6:
		request.AddressBytes, n, err = util.ReadMessageByReader(reader, 16)
		if err != nil {
			return
		}
		if n != 16 {
			err = errors.New("长度不正确")
			return
		}
	default:
		err = errors.New("不支持该地址类型")
	}

	// 读取端口
	var portBytes []byte
	portBytes, n, err = util.ReadMessageByReader(reader, 2)
	if err != nil {
		return
	}
	if n != 2 {
		err = errors.New("长度不正确")
		return
	}
	binary.Read(bytes.NewReader(portBytes), binary.BigEndian, &request.Port)

	if request.AddressType == DomainName { // 将目标地址[]byte转string
		request.Address = string(request.AddressBytes)
	} else {
		request.Address = util.Bytes2Ip(request.AddressBytes)
	}

	request.Target = net.JoinHostPort(request.Address, strconv.Itoa(int(request.Port))) // 拼接完成目标地址

	return
}
