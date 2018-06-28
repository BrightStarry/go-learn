package util

import (
	"time"
	"fmt"
)

// 认证模式
const(
	// 无需认证
	UnMethod = 0x00
)

// 地址类型
const(
	Ipv4 = 0x01
	DomainName = 0x03
	Ipv6 = 0x04
)

const(
	// 默认版本
	Version = 0x05
	// 默认值 0x01
	One = 0x01
	// 默认值 0x00
	Zero = 0x00
	// 表示不支持客户端的认证方式
	NotSupport = 0xff
	// 读取超时时间
	ReadTimeout = 15 * time.Second
	// 数据发送超时时间
	WriteTimeout = 15 * time.Second
)

// 配置
var Config = NewDefaultConfig()


/*系统配置*/
type Configuration struct{
	// 启动端口
	JumpPort string
	// web服务端口
	WebPort string
	// pac获取url
	PacUrl string

}

/*toString方法*/
func (this Configuration) String() string {
	return fmt.Sprintln("启动端口:",this.JumpPort,",web服务端口:",this.WebPort,",pac获取url:",this.PacUrl)
}

/*构造默认的系统配置*/
func NewDefaultConfig() *Configuration {
	return &Configuration{
		JumpPort: "8081",
		WebPort:  "9999",
		PacUrl: "http://106.14.7.29:9000/pac",
	}
}