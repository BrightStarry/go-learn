package util

import (
	"time"
	"fmt"
)

// 认证模式
const(
	// 无需认证
	UnMethod = 0x00
	UnMethodName = "no"
	// 密码认证
	PwdMethod = 0x02
	PwdMethodName = "pwd"
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
	// 用户名
	Username string
	// 密码
	Password string
	// 启动端口
	Port string
	// 当前认证模式
	Method byte


}

/*toString方法*/
func (this Configuration) String() string {
	if this.Method == UnMethod {
		return fmt.Sprintln("认证模式:",UnMethodName,",端口:",this.Port)
	}else{
		return fmt.Sprintln("认证模式:",PwdMethodName,",端口:",this.Port,",用户名:",this.Username,",密码:",this.Password)
	}

}

/*构造默认的系统配置*/
func NewDefaultConfig() *Configuration {
	return &Configuration{
		Username: "zx",
		Password: "123456",
		Port:     "8081",
		Method:   UnMethod,
	
	}
}