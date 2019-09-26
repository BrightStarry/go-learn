package main

import (
	"flag"
	"errors"
	"log"
	"zx/jump/client/config"
	"zx/jump/util"
)

func main() {
	util.StartTcpServer(config.Config.Port,config.Handler)
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
}

/*读取外部参数*/
func readExternalParam() (err error) {
	flag.StringVar(&config.Config.Port, "port", config.Config.Port, "服务器端口")
	flag.StringVar(&config.Config.Pattern, "pattern", config.Config.Pattern, "代理模式 common:socks5; cs:client-server")
	method := flag.String( "method", config.UnMethodName, "认证方式 no:无需认证; pwd:用户名/密码认证")
	flag.StringVar(&config.Config.Username, "username", config.Config.Username, "用户名")
	flag.StringVar(&config.Config.Password, "pwd", config.Config.Password, "密码")
	flag.StringVar(&config.Config.Server, "server", config.Config.Server, "服务器地址,用于CS代理模式， ip:port")
	flag.StringVar(&config.Config.ServerPwd, "spwd", config.Config.ServerPwd, "服务器密码")
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




