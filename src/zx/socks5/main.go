package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"zx/socks5/config"
)



func main() {
	startTcpServer(config.Config.Port, Handler)
}

/*
	启动tcp服务
*/
func startTcpServer(port string,handler func (*net.TCPConn)) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln("服务启动失败:", err)
	}
	log.Println("服务启动成功!")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("接受客户端连接失败:", err)
		}
		log.Println("accept connect:", conn.RemoteAddr().String())
		// 处理请求
		go handler(conn.(*net.TCPConn))
	}
}

/**
初始化
*/
func init() {
	err := readExternalParam()
	if err != nil {
		log.Fatalln(errors.New("参数读取失败:" +err.Error()))
	}
	log.Println("准备启动，当前参数:", config.Config)

	// 设置超时时间
	config.ReadTimeout = time.Duration(config.Config.SocketTimeout) * time.Second
	config.WriteTimeout = time.Duration(config.Config.SocketTimeout) * time.Second
	config.ConnectTimeout = time.Duration(config.Config.ConnectTimeout) * time.Second
}

/*读取外部参数*/
func readExternalParam() (err error) {
	// 根据默认配置config.Config的值，设置默认方法名
	var defaultMethodName string
	switch config.Config.Method {
	case config.UnMethod:
		defaultMethodName = config.UnMethodName
	case config.PwdMethod:
		defaultMethodName = config.PwdMethodName
	}

	flag.StringVar(&config.Config.Port, "port", config.Config.Port, "服务器端口")
	flag.StringVar(&config.Config.Pattern, "pattern", config.Config.Pattern, "代理模式 common:socks5; cs:client-server")
	method := flag.String( "method",  defaultMethodName, "认证方式 no:无需认证; pwd:用户名/密码认证")
	flag.StringVar(&config.Config.Username, "username", config.Config.Username, "用户名,用于pwd方式")
	flag.StringVar(&config.Config.Password, "pwd", config.Config.Password, "密码,用于pwd方式")
	flag.StringVar(&config.Config.Server, "server", config.Config.Server, "服务器地址,用于CS代理模式， ip:port")
	flag.StringVar(&config.Config.ServerPwd, "spwd", config.Config.ServerPwd, "服务器密码,用于CS代理模式")
	flag.IntVar(&config.Config.SocketTimeout,"socketTimeout",config.Config.SocketTimeout,"socket超时时间，秒")
	flag.IntVar(&config.Config.ConnectTimeout,"connectTimeout",config.Config.ConnectTimeout,"socket连接建立超时时间，秒")
	flag.Parse()

	switch *method {
	case config.UnMethodName:
		config.Config.Method = config.UnMethod
	case config.PwdMethodName:
		config.Config.Method = config.PwdMethod
	default:
		return errors.New("无法识别该认证方式")
	}

	config.Config.ServerPwdLen = byte(len(config.Config.ServerPwd))
	config.Config.ServerPwdByte = []byte(config.Config.ServerPwd)

	return
}


/*socket相关*/

/*
	处理tcp请求
*/
func Handler(conn *net.TCPConn) {
	defer conn.Close()
	// 预处理，握手，密码认证等
	err := HandlerPre(conn)
	if err != nil {
		log.Println("握手失败：", err)
		return
	}
	// 处理连接
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
	if config.Config.Method == config.PwdMethod {
		err = handlerPwdAuthentication(conn)
	}
	return
}

/*处理握手*/
func handlerHandshake(conn *net.TCPConn) (err error) {
	// 读取客户端握手请求对象,生成响应对象
	request := new(config.HandshakeRequest) // 创建握手请求对象

	var buf []byte // 读取字节
	if buf,err = ReadMessage(conn, 512);err != nil {
		return
	}
	reader := bytes.NewReader(buf)

	if request.Version, err = reader.ReadByte(); err != nil { // 读取版本号 默认值为0x05
		return
	}

	if request.Version != config.Version { // 校验版本
		err = errors.New("只支持SOCKET5")
		return
	}

	if request.MethodsLen, request.Methods, err = ReadByLen(reader); err != nil { // 读取用户名
		return
	}

	if request.MethodsLen == config.Zero { // Methods长度字段,如果为0，表示支持的方法数为0
		err = errors.New("客户端methods长度字段为0，不支持任何方法")
		return
	}

	// 根据方法长度，读取客户端支持的对应方法数组
	methodLen := int(request.MethodsLen)
	isSupport := false
	for i:=0;i<methodLen ; i++ {
		if request.Methods[i] == config.Config.Method {
			isSupport = true
			break
		}
	}
	if !isSupport {
		err = errors.New("不支持客户端的认证方式")
		return
	}
	// 发送响应
	response := &config.HandshakeResponse{config.Version, config.Config.Method} // 响应对象
	err = SendMessage(conn, response)
	return // 返回err
}

/*处理密码认证*/
func handlerPwdAuthentication(conn *net.TCPConn) (err error) {
	var buf []byte // 读取字节
	if buf,err= ReadMessage(conn, 512);err != nil {
		return
	}
	reader := bytes.NewReader(buf)

	request := new(config.PwdAuthenticationRequest) // 返回对象

	if request.Pointless, err = reader.ReadByte(); err != nil { // 读取无意义标识
		return
	}

	if request.UsernameLength, request.Username, err = ReadByLen(reader); err != nil { // 读取用户名
		return
	}

	if request.PasswordLength, request.Password, err = ReadByLen(reader); err != nil { // 读取密码
		return
	}

	response := &config.PwdAuthenticationResponse{config.One, config.Zero} // 读取成功才对客户端进行响应

	if config.Config.Username != string(request.Username) || config.Config.Password != string(request.Password) { //校验密码
		// 失败响应
		response.Result = config.One
		SendMessage(conn, response)
		err = errors.New("用户密码错误")
		return
	}
	err = SendMessage(conn, response) // 成功响应
	return
}

// 连接请求 失败响应
var readConnectFailedResponse = config.ConnectResponse{
	Version:     config.Version,
	Response:    config.One,
	Reserve:     config.Zero,
	AddressType: config.Ipv4,
	Address:     []byte{0, 0, 0, 0},
	Port:        []byte{0,0},
}
// 连接请求 成功响应
var readConnectSuccessResponse = config.ConnectResponse{
	Version:     config.Version,
	Response:    config.Zero,
	Reserve:     config.Zero,
	AddressType: config.Ipv4,
	Address:     []byte{0, 0, 0, 0},
	Port:        []byte{0,0},
}

/*处理连接请求*/
func HandlerConnect(conn *net.TCPConn) (err error) {
	// 读取请求
	request, err := readConnectRequest(conn)
	if err != nil {
		SendMessage(conn, readConnectFailedResponse)
		return
	}

	var targetConn *net.TCPConn
	// 连接到目标服务器
	switch config.Config.Pattern {
	case config.Common:
		// 普通socks5，直接连接客户端请求的目标服务器
		targetConn, err = ConnectToTarget(request.Target)
		if err != nil {
			SendMessage(conn, readConnectFailedResponse)
			return
		}
		defer targetConn.Close()

	case config.CS:
		// 先连接到另一台服务器，再连接到客户端请求的目标服务器
		return errors.New("暂不支持")
	default:
		log.Fatalln("模式错误,当前模式:", config.Config.Pattern)
	}

	// 连接成功
	// ip和端口应该是可以不返回的
	//response.Address, response.Port = Ip2Bytes(targetConn.LocalAddr().String())
	// 发送响应
	if err = SendMessage(conn, readConnectSuccessResponse); err != nil {
		return
	}
	go io.Copy(targetConn, conn)
	io.Copy(conn, targetConn)

	return
}

/*读取连接请求*/
func readConnectRequest(conn *net.TCPConn) (request *config.ConnectRequest, err error) {
	// 读取字节
	var buf []byte
	buf, err = ReadMessage(conn, 512)
	if err != nil {
		return
	}
	// 转换为读取器
	reader := bytes.NewReader(buf)

	// 必须先创建该对象
	request = new(config.ConnectRequest)

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
	case config.Ipv4:
		request.AddressBytes, n, err = ReadMessageByReader(reader, 4)
		if err != nil {
			return
		}
		if n != 4 {
			err = errors.New("长度不正确")
			return
		}
		// 域名，下个字节表示域名长度
	case config.DomainName:
		_, request.AddressBytes, err = ReadByLen(reader)
		if err != nil {
			return
		}
		// ipv6
	case config.Ipv6:
		request.AddressBytes, n, err = ReadMessageByReader(reader, 16)
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
	portBytes, n, err = ReadMessageByReader(reader, 2)
	if err != nil {
		return
	}
	if n != 2 {
		err = errors.New("长度不正确")
		return
	}
	binary.Read(bytes.NewReader(portBytes), binary.BigEndian, &request.Port)

	if request.AddressType == config.DomainName { // 将目标地址[]byte转string
		request.Address = string(request.AddressBytes)
	} else {
		request.Address = Bytes2Ip(request.AddressBytes)
	}
	log.Println("connect:",request.Address)
	request.Target = net.JoinHostPort(request.Address, strconv.Itoa(int(request.Port))) // 拼接完成目标地址

	return
}


func BytesToInt(bys []byte) int {
	byteBuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(byteBuff, binary.BigEndian, &data)
	return int(data)
}
/*连接到目标服务器*/
func ConnectToTarget(target string) (conn *net.TCPConn, err error) {
	var dial net.Conn
	dial, err = net.DialTimeout("tcp", target, config.ConnectTimeout)
	if err != nil {
		return
	}
	conn = dial.(*net.TCPConn)
	return
}

/*读取消息,指定默认读取长度*/
func ReadMessage(conn *net.TCPConn,defaultLen int) ( []byte, error) {
	// 读取字节
	buf := make([]byte,defaultLen)
	if err := conn.SetReadDeadline(time.Now().Add(config.ReadTimeout));err != nil {
		return nil,err
	}
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
func ReadMessageByReader(reader *bytes.Reader,defaultLen int) ( []byte, int, error) {
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

/*发送消息*/
func SendMessage(conn *net.TCPConn,data interface{})(err error) {
	// 设置超时时间
	if err = conn.SetWriteDeadline(time.Now().Add(config.WriteTimeout));err != nil {
		return
	}

	switch data.(type) {
	case config.Byteable:
		_,err = conn.Write(data.(config.Byteable).ToBytes())
	default:
		// 发送对象， 该方法的data支持定长的类型（切片等都需要定长）
		err = binary.Write(conn,binary.LittleEndian,data)
	}
	return
}

/*读取1个字节的长度字节，并读取对应长度的后续字节*/
func ReadByLen(reader *bytes.Reader) (length byte,data []byte,err error) {
	length,err = reader.ReadByte()
	if err != nil {
		return
	}
	var n int
	data,n,err = ReadMessageByReader(reader,int(length))
	if n != int(length) {
		err = errors.New("长度不正确")
		return
	}
	return
}

/*[]byte转ip 字符串*/
func Bytes2Ip(data []byte) string {
	return net.IPv4(data[0],data[1],data[2],data[3]).String()
}

/*
	地址字符串转ip([]byte) 和 port(uint16)
	类似 192.168.1.111:54698, 如果没有端口，返回0

	注意，该方法返回的切片是定长的
*/
func Ip2Bytes(addr string) (ip []byte,port uint16) {
	// 转为 ["192.168.1.111","54698"] 或 （如果没有端口）["192.168.1.111"]
	arr := strings.Split(addr,":")

	ip = make([]byte,4,4)
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
