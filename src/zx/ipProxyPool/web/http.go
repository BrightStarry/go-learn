package web

import (
	"net/http"
	"log"
	"encoding/json"
	"zx/ipProxyPool/config"
	"strconv"
	"zx/ipProxyPool/util"
)

/**
	web服务,包括监控/接口相关
 */

/**
	启动web服务
 */
func SyncStartWebServer() {
	// 监控器
	http.HandleFunc("/monitor", errWrapper(monitorHandler))
	http.HandleFunc("/get", errWrapper(getIpHandler))
	if err := http.ListenAndServe(":"+config.Config.WebPort, nil); err != nil {
		log.Panicln("web服务异常:", err)
	}
}

/**
	监控
 */
func monitorHandler(w http.ResponseWriter, r *http.Request) error {
	// 返回json头
	w.Header().Set("content-type", "text/json; charset=utf-8")
	result := config.MonitorParam{
		IpPoolLen:         len(config.ProxyIpStore.Queue),
		ReVerifyChanLen:   len(config.ReVerifyChan),
		//VerifiedChanLen:   len(config.VerifiedChan),
		WaitVerifyChanLen: len(config.WaitVerifyChan),
	}
	bytes, _ := json.Marshal(result)
	w.Write(bytes)
	return nil
}

/**
	获取ip
 */
func getIpHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	l := r.Form.Get("len")
	length, err := strconv.Atoi(l)
	// 如果len为nil或格式有误，默认为最大
	if err != nil {
		length = 9999999
	}
	isJump := r.Form.Get("isJump")


	// 获取ip
	//ips := store.GetIpsAtLast(length)
	ips := config.ProxyIpStore.Queue[:]
	if length > len(ips){
		length = len(ips)
	}
	// 根据延迟从小到大排序
	ips = util.Sort(ips, len(ips))[:length]

	// 如果只需要翻墙服务器
	if isJump == "1"{
		var tempIps []*config.ProxyIp
		for _,v:= range ips{
			if v.IsJump{
				tempIps = append(tempIps, v)
			}
		}
		ips = tempIps
	}

	var result []*config.IpDTO
	var p string
	for _, v := range ips {
		if v.Protocol == config.HttpFlag {
			p = config.Http
		} else if v.Protocol == config.HttpsFlag {
			p = config.Https
		}
		result = append(result, &config.IpDTO{
			Host:           v.Url.Host,
			Protocol:       p,
			LastVerifyTime: v.LastVerifyTime,
			Delay:          v.Delay.Seconds(),
			IsJump:         v.IsJump,
		})
	}

	bytes, _ := json.Marshal(result)
	w.Write(bytes)
	return nil
}

/**
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
