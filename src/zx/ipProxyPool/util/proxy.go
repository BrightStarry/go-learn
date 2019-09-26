package util

import (
	"os/exec"
	"net/url"
)

/**
	代理设置相关
 */
 const(
	 ProxyEnable = "1"
	 ProxyDisable = ""
 )


 /**
 取消代理
  */
  func CancelProxy(){
	  cmd := exec.Command("reg","add","HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		  "/v","ProxyEnable", "/t","REG_SZ", "/d",ProxyDisable, "/f")
	  cmd.Run()
	  flushProxy()
  }

  /**
  设置代理
   */
func SetProxy(url *url.URL){
	//先取消，消除缓存
	CancelProxy()
	cmd := exec.Command("reg","add","HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		"/v","ProxyEnable", "/t","REG_SZ", "/d","1", "/f")
	cmd.Run()

	cmd2 := exec.Command("reg","add","HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		"/v","ProxyServer", "/t","REG_SZ", "/d",url.Host, "/f")
	cmd2.Run()

	flushProxy()
}

 /**
 刷新代理
  */
 func flushProxy(){
	 cmd := exec.Command("taskkill","/f","/im","explorer.exe")
	 cmd.Run()
	 cmd2 := exec.Command("explorer.exe")
	 cmd2.Start()
 }
