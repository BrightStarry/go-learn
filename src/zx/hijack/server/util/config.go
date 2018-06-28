package util

import "fmt"

/*配置类*/

/*系统配置*/
type Configuration struct{
	// 端口
	Port string
	// pac文件路径
	PacPath string
}

/*toString方法*/
func (this Configuration) String() string {
	return fmt.Sprintln("启动端口:",this.Port,",pac文件路径:",this.PacPath)
}


// 配置
var Config = NewDefaultConfig()

/*构造默认的系统配置*/
func NewDefaultConfig() *Configuration {
	return &Configuration{
		Port: "9000",
		//PacPath:"C:\\code\\goSpace\\go-learn\\src\\zx\\hijack\\resources\\test.pac",
		PacPath:"./test.pac",
	}
}