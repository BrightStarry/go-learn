package util

import (
	"net"
	"time"
	"encoding/binary"
	"errors"
	"bytes"
	"sync"
	"strconv"
	"log"
)

/*socket相关*/

// 用于增加用户到map的锁
var userLock = new(sync.Mutex)



/*
	启动tcp服务
*/
func StartTcpServer(port string) {
	listener,err  := net.Listen("tcp",":" + port)
	if err != nil {
		log.Fatalln("服务启动失败:",err)
	}
	log.Println("服务启动成功,地址:",listener.Addr().String())
	defer listener.Close()

	for{
		conn,err  := listener.Accept()
		if err != nil {
			log.Fatalln("接受客户端连接失败:",err)
		}
		// 处理请求
		go HandlerTcpConn(conn.(*net.TCPConn))
	}
}

/*
	处理tcp请求
*/
func HandlerTcpConn(conn *net.TCPConn) {
	//// 如果是密码认证
	//if util.Config.Method == util.PwdMethod {
	//	// 获取该用户信息
	//	_,ok := util.AuthenticationUser[conn.RemoteAddr().String()]
	//	// 不存在
	//	if !ok{
	//		// 进行预处理（握手等）
	//		util.HandlerPre(conn)
	//	}
	//}else {
	//	// 如果是无需验证
	//	util.HandlerPre(conn)
	//}

	HandlerPre(conn)
	HandlerConnect(conn)
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
	if Config.Method == PwdMethod{
		var user User
		user,err = handlerPwdAuthentication(conn)
		if err != nil {
			return
		}
		/**
			增加用户到map
		*/
		addUser(user)
	}

	return
}

/*处理握手*/
func handlerHandshake(conn *net.TCPConn) error {
	// 读取客户端握手请求对象,生成响应对象
	handshakeResponse,err := readHandshakeRequestGenerateResponse(conn)
	// 响应不为空，则表示认证不支持或成功，发送响应
	if handshakeResponse != nil {
		// 发送响应
		err = sendResponse(conn,handshakeResponse)
		if err != nil {
			return err
		}
	}
	// 有异常，直接退出
	if err != nil {
		return err
	}
	return nil
}

/*处理密码认证*/
func handlerPwdAuthentication(conn *net.TCPConn) (user User,err error){
	// 读取
	var request *PwdAuthenticationRequest
	request,err = readPwdAuthenticationRequest(conn)
	if err != nil {
		return
	}
	// 读取成功才对客户端进行响应
	response := PwdAuthenticationResponse{One,Zero}

	/**
		校验密码
	 */
	 if Config.Username != string(request.Username) || Config.Password != string(request.Password) {
	 	// 失败响应
	 	response.Result = One
	 	sendResponse(conn,response)
	 	err = errors.New("用户密码错误")
	 	return
	 }

	 // 成功响应
	 err = sendResponse(conn,response)
	 if err != nil {
	 	return
	 }

	 user = User{conn.RemoteAddr().String()}
	 return
}

/*读取客户端握手请求对象,生成响应对象*/
func readHandshakeRequestGenerateResponse(conn *net.TCPConn) (response *HandshakeResponse,err error){
	// 设置读取超时时间
	conn.SetReadDeadline(time.Now().Add(ReadTimeout))

	// 读取为对象
	request := new(HandshakeRequest)
	err = binary.Read(conn,binary.LittleEndian,request)
	if err != nil {
		return
	}
	// 校验版本
	if request.Version != Version {
		err = errors.New("只支持SOCKET5")
		return
	}
	// Methods长度字段
	if request.NMethods != One {
		err = errors.New("methods长度字段不为1")
		return
	}

	// 响应对象
	response = &HandshakeResponse{Version,Config.Method}

	// methods字段
	if request.Methods != Config.Method{
		response.Method = NotSupport
		err = errors.New("不支持客户端的认证方式")
		return
	}
	return
}



/*读取密码认证请求对象*/
func readPwdAuthenticationRequest(conn *net.TCPConn) (result *PwdAuthenticationRequest,err error){
	// 设置超时时间
	conn.SetReadDeadline(time.Now().Add(ReadTimeout))
	// 读取字节
	var buf []byte
	buf,err = readMessage(conn,512)
	if err != nil {
		return
	}

	// 返回对象
	result = new(PwdAuthenticationRequest)

	// 转换为读取器
	reader := bytes.NewReader(buf)
	// 读取无意义标识
	result.Pointless,err = reader.ReadByte()
	if err != nil {
		return
	}
	// 读取用户名
	result.UsernameLength,result.Username,err = readByLenField(reader)
	if err != nil {
		return
	}
	// 读取密码
	result.PasswordLength,result.Password,err = readByLenField(reader)
	if err != nil {
		return
	}
	return
}



/*处理连接请求*/
func HandlerConnect(conn *net.TCPConn) (error) {
	// 构造一个默认的失败响应
	response := ConnectResponse{
		Version:Version,
		Response:One,
		Reserve:Zero,
		AddressType:Ipv4,
		Address: []byte{0,0,0,0},
		Port:0,
	}

	// 读取请求
	request,err  := readConnectRequest(conn)
	if err != nil {
		sendResponse(conn,response)
		return err
	}

	// 连接到目标服务器
	targetConn,err := connectToTarget(request.Target)
	if err != nil {
		sendResponse(conn,response)
		return err
	}

	// 连接成功
	response.Response = Zero
	response.Address,request.Port = Ip2Bytes(targetConn.LocalAddr().String())
	err = sendResponse(conn,response)
	if err != nil {
		return nil
	}
	// TODO
	addTargetConn(conn.RemoteAddr().String())

	return nil
}

/*连接到目标服务器*/
func connectToTarget(target string) (conn *net.TCPConn, err error) {
	var dial net.Conn
	dial,err = net.DialTimeout("tcp",target,ReadTimeout)
	if err != nil {
		return
	}
	conn = dial.(*net.TCPConn)
	return
}

/*读取连接请求*/
func readConnectRequest(conn *net.TCPConn) (request *ConnectRequest,err error) {
	// 设置超时时间
	conn.SetReadDeadline(time.Now().Add(ReadTimeout))
	// 读取字节
	var buf []byte
	buf,err = readMessage(conn,512)
	if err != nil {
		return
	}
	// 转换为读取器
	reader := bytes.NewReader(buf)

	// 必须先创建该对象
	request = new(ConnectRequest)

	// 读取版本
	request.Version,err = reader.ReadByte()
	if err != nil {
		return
	}
	// 读取命令
	request.CMD,err = reader.ReadByte()
	if err != nil {
		return
	}
	// 保留字段
	request.Reserve,err = reader.ReadByte()
	if err != nil {
		return
	}
	// 地址类型
	request.AddressType,err = reader.ReadByte()
	if err != nil {
		return
	}

	var n int
	// 目标地址
	switch request.AddressType {
	// ipv4,4个字节
	case Ipv4:
		request.AddressBytes,n,err = readMessageByReader(reader,4)
		if err != nil {
			return
		}
		if n != 4 {
			err = errors.New("长度不正确")
			return
		}
	// 域名，下个字节表示域名长度
	case DomainName:
		_,request.AddressBytes,err  = readByLenField(reader)
		if err != nil {
			return
		}
	// ipv6
	case Ipv6:
		request.AddressBytes,n,err = readMessageByReader(reader,16)
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
	portBytes,n,err = readMessageByReader(reader,2)
	if err != nil {
		return
	}
	if n != 2 {
		err = errors.New("长度不正确")
		return
	}
	binary.BigEndian.PutUint16(portBytes,request.Port)

	// 将目标地址[]byte转string
	if request.AddressType == DomainName {
		request.Address = string(request.AddressBytes)
	}else {
		request.Address = Bytes2Ip(request.AddressBytes)
	}
	// 拼接完成目标地址
	if request.Port != 0 {
		request.Target = request.Address + ":" + strconv.Itoa(int(request.Port))
	}else {
		request.Target = request.Address
	}

	return
}





