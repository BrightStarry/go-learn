package main

import (
	"zx/jump/util"
	"zx/jump/server/config"
	"flag"
	"log"
)

/*服务器*/

func main() {
	util.StartTcpServer(config.Config.JumpPort,config.Handler)
}

func init() {
	readExternalParam()
	log.Println("初始化，当前参数:  ",config.Config)
}


/*读取外部参数*/
func readExternalParam()  {
	flag.StringVar(&config.Config.JumpPort, "jport", config.Config.JumpPort, "代理服务端口")
	flag.StringVar(&config.Config.WebPort, "wport", config.Config.WebPort, "服务器端口")
	flag.StringVar(&config.Config.Pwd, "pwd", config.Config.Pwd, "密码")
	flag.StringVar(&config.Config.Encrypt, "encrypt", config.Config.Encrypt, "加密方式(no:不加密)")
	flag.Parse()

	config.Config.PwdByte = []byte(config.Config.Pwd)
}
