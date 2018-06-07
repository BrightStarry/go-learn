package main

import (
	"jump/util"
	"flag"
	"errors"
	"log"
)




func main() {
	util.StartTcpServer(util.Config.Port)
}

/**
	 初始化
 */
func init() {
	err := readExternalParam()
	if err != nil {
		panic(errors.New("参数读取失败:" +err.Error()))
	}
	log.Println("准备启动，当前参数:",util.Config)
}

/*读取外部参数*/
func readExternalParam() (err error) {
	// 端口
	var port string
	flag.StringVar(&port, "port", "8081", "服务器端口")

	// 认证模式
	var method string
	flag.StringVar(&method, "method", "no", "认证方式 no:无需认证; pwd:用户名/密码认证")
	var username string
	var password string
	// 用户名/密码， 如果是该认证模式
	flag.StringVar(&username, "username", util.Config.Username, "认证用户名")
	flag.StringVar(&password, "pwd", util.Config.Password, "认证密码")
	flag.Parse()

	switch method {
	case util.UnMethodName:
		util.Config.Method = util.UnMethod
	case util.PwdMethodName:
		util.Config.Method = util.PwdMethod
	default:
		return errors.New("无法识别该认证方式")
	}
	util.Config.Port = port
	util.Config.Username = username
	util.Config.Password = password
	return
}




