package main

import (
	"jump/util"
	"flag"
	"errors"
)




func main() {
	util.StartTcpServer(util.Config.Port)
}

/*读取外部参数*/
func readExternalParam() (err error) {
	// 端口
	flag.StringVar(&util.Config.Port, "port", "8081", "服务器端口，默认8081")

	// 认证模式
	var method string
	flag.StringVar(&method, "method", "no", "认证方式 no:无需认证; pwd:用户名/密码认证")
	switch method {
	case util.UnMethodName:
		util.Config.Method = util.UnMethod
	case util.PwdMethodName:
		util.Config.Method = util.PwdMethod
	default:
		return errors.New("无法识别该认证方式")
	}

	// 用户名/密码， 如果是该认证模式
	flag.StringVar(&util.Config.Username, "username", "zx", "认证用户名，默认zx")
	flag.StringVar(&util.Config.Password, "pwd", "zx", "认证密码,默认123456")



}


