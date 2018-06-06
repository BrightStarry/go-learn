package util

/*网络传输对象*/

/*
	客户端握手请求对象
*/
type HandshakeRequest struct {
	// 版本，  0x05
	Version uint8
	// Methods字段占用字节数, 0x01
	NMethods uint8
	/*
		客户端支持的认证方式,
		0x00 不认证; 0x01 通用安全服务应用程序接口(GSSAPI); 0x02用户名/密码(USERNAME/PASSWORD); 0xff 没有可接受方法
	*/
	Methods uint8
}


/*
	握手响应对象
*/
type HandshakeResponse struct{
	// 版本, 0x05
	Version uint8
	// 返回服务端支持的认证方法, 如果客户端的所有认证方式，服务端都不支持，则返回0xff, 目前返回0x02
	Method uint8
}

/*
	客户端密码认证请求
*/
type PwdAuthenticationRequest struct{
	// 一个无意义标识, 0x01
	Pointless uint8
	// 用户名长度
	UsernameLength uint8
	// 用户名,根据UsernameLength确定字节长度
	Username []byte
	// 密码长度
	PasswordLength uint8
	// 密码，根据PasswordLength确定长度
	Password []byte
}

/*
	密码认证响应
*/
type PwdAuthenticationResponse struct{
	// 一个无意义标识, 0x01
	Pointless uint8
	// 验证结果 0x00成功； 其余失败
	Result uint8
}

/*
	客户端连接请求
*/
type ConnectRequest struct{
	// 版本
	Version uint8
	// 命令 0x01：CONNECT 建立 TCP 连接; 0x02: BIND 上报反向连接地址; 0x03：关联 UDP 请求
	CMD uint8
	// 保留字段0x00
	Reserve uint8
	// 地址类型 01:ipv4; 03:域名; 0x04:ipv6
	AddressType uint8
	// 目标地址，根据不同地址类型有不同长度 0x01：4 个字节的 IPv4 地址；0x03：1 个字节表示域名长度，紧随其后的是对应的域名；0x04：16 个字节的 IPv6 地址
	AddressBytes []byte
	// 目标地址，将[]byte转为string
	Address string
	// 目标端口
	Port uint16
	// 目标，直接将地址和端口拼接
	Target string

}

/*
	连接请求响应
*/
type ConnectResponse struct{
	// 版本
	Version uint8
	/**
		响应
		'00'成功; '01'一般的socket服务异常; '02': 规定不允许该连接；
		'03'网络异常； '04'：主机异常；'05'连接拒绝； '06' TTL过期； '07'命令不支持；
		'08'地址类型不支持
	*/
	Response uint8
	// 保留字段0x00
	Reserve uint8
	// 地址类型 01:ipv4; 03:域名; 0x04:ipv6
	AddressType uint8
	// 服务端的ip
	Address []byte
	// 服务端连接到目标服务器时，服务端的端口
	Port uint16
}

/*用户信息*/
type User struct{
	Ip string
}


