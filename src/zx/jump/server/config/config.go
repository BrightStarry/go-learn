package config

import (
	"time"
	"fmt"
)

/*配置类*/

// 加密方式
const(
	// 不加密
	NotEncrypt = "no"
)

// 默认配置
const(
	// 读取超时时间
	ReadTimeout = 15 * time.Second
	// 数据发送超时时间
	WriteTimeout = 15 * time.Second
)

var Config = NewDefaultConfig()

//服务端配置
type Configuration struct {
	// 代理服务端口
	JumpPort string
	// web服务端口
	WebPort string
	// 密码
	Pwd string
	// 加密方式
	Encrypt string


	// 密码字节
	PwdByte []byte
}

func (this *Configuration) String() string{
	return fmt.Sprintln("代理服务端口:",this.JumpPort,",web服务端口:",this.WebPort,",密码:",this.Pwd,",加密方式:",this.Encrypt)
}

/*构造默认的系统配置*/
func NewDefaultConfig() *Configuration {
	return &Configuration{
		JumpPort: "9703",
		WebPort:"9745",
		Pwd:"123456",
		Encrypt:"no",
	}
}



