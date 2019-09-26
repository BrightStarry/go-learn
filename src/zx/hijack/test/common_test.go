package test

import (
	"os/exec"
	"testing"
	"net/http"
	"log"
	"io/ioutil"
)

// 通用测试类

/**
	测试自动配置pac脚本
 */
func TestAutoConfigPAC(t *testing.T) {
	cmd := exec.Command("reg","add","HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		"/v","AutoConfigURL", "/t","REG_SZ", "/d","http://127.0.0.1:9000/pac", "/f")
	cmd.Run()
	cmd2 := exec.Command("ipconfig","/flushdns")
	cmd2.Run()
}

/**
	测试自动配置代理服务器
 */
func TestAutoConfigProxy(t *testing.T) {
	cmd := exec.Command("reg","add","HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		"/v","ProxyEnable", "/t","REG_SZ", "/d","1", "/f")
	cmd.Run()

	cmd2 := exec.Command("reg","add","HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		"/v","ProxyServer", "/t","REG_SZ", "/d","192.168.0.1:8088", "/f")
	cmd2.Run()
	cmd3 := exec.Command("ipconfig","/flushdns")
	cmd3.Run()
}

/**
测试重启资源管理器，让http代理服务器生效
taskkill /f /im explorer.exe & start explorer.exe

 */
func TestRestarResourceMaster(t *testing.T)  {
	cmd := exec.Command("taskkill","/f","/im","explorer.exe")
	cmd.Run()

	cmd2 := exec.Command("explorer.exe")
	cmd2.Start()


}

/**
	测试启动web服务，返回pac脚本
 */
func TestWeb(t *testing.T) {
	SyncStartWebServer()
}

var pacByte = readPac()

/**
	启动web服务
 */
func SyncStartWebServer() {

	http.HandleFunc("/pac", errWrapper(pac))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panicln("web服务异常:", err)
	}
}

/**
读取pac
 */
func readPac() []byte{
	bytes, err := ioutil.ReadFile("..\\resources\\test.pac")
	if err != nil {
		log.Println(err)
	}
	return bytes
}

/**
	返回pac文件
 */
 func pac(w http.ResponseWriter, r *http.Request)error{
	 w.Header().Set("content-type", "application/x-ns-proxy-autoconfig")
	 w.Write(pacByte)
	 return nil
 }


/*
http处理方法
*/
type customHttpHandler func(w http.ResponseWriter, r *http.Request) error

/**
	统一异常处理
 */
func errWrapper(handler customHttpHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 异常捕获,此处返回 500系统异常
		defer func() {
			if err := recover(); err != nil {
				log.Println("web服务异常:", err)
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		// 调用业务方法
		err := handler(w, r)
		if err == nil {
			return
		}
		// 如果是自定义异常
		if customErr, ok := err.(CustomError); ok {
			http.Error(w, customErr.Message(), http.StatusBadRequest)
			return
		}
		// 否则返回 500
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)

	}
}

/**
	自定义异常接口
 */
type CustomError interface {
	// 整合接口,让userError接口必须实现error接口的Error()方法
	error
	Message() string
}

/**
	自定义异常类型
 */
type ServiceError string

func (this ServiceError) Error() string {
	return this.Message()
}
func (this ServiceError) Message() string {
	return string(this)
}

