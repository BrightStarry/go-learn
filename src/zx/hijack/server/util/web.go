package util

import (
	"net/http"
	"log"
	"io/ioutil"
)

var pacByte []byte

/**
	启动web服务
 */
func SyncStartWebServer() {
	readPac() // 读取配置文件
	http.HandleFunc("/pac", errWrapper(pac))
	if err := http.ListenAndServe(":" + Config.Port, nil); err != nil {
		log.Panicln("web服务异常:", err)
	}
}

/**
	从本地读取pac
 */
func readPac(){
	bytes, err := ioutil.ReadFile(Config.PacPath)
	if err != nil {
		log.Fatalln("读取pac文件失败:",err)
	}
	pacByte = bytes
}

/**
	返回pac文件
 */
func pac(w http.ResponseWriter, r *http.Request)error{
	//w.Header().Set("content-type", "application/x-ns-proxy-autoconfig")
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