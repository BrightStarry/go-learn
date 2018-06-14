package config

import (
	"net/http"
	"time"
	"net/http/cookiejar"
	"log"
)

/**
	参数相关
 */
// 保存所有网站信息
var WebInfos  []WebInfo
// 系统配置
var Config *SystemConfig
// 默认请求头
var DefaultHeader *map[string][]string
// client
var DefaultClient *http.Client

/**
   初始化方法
*/
func Init() {
	// 初始化系统参数
	InitSystemConfig()
	// 初始化网站信息
	InitWebInfos()
	// 初始化client
	DefaultClient = InitDefaultClient()
}


/**
	系统参数
 */
type SystemConfig struct {
	// 爬虫默认超时时间
	SpiderTimeout time.Duration
}

/**
	初始化系统参数
 */
func InitSystemConfig() {
	Config = &SystemConfig{
		SpiderTimeout: 10 * time.Second,
	}
}

/**
	初始化client
 */
func InitDefaultClient() *http.Client{
	DefaultClient := &http.Client{}
	// 超时时间
	DefaultClient.Timeout =  Config.SpiderTimeout
	// 构建cookie
	cookie,err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln("构建cookie失败:",err)
	}
	DefaultClient.Jar = cookie

	// 设置默认请求头
	DefaultHeader = &map[string][]string{
		"User-Agent" : {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36"},
		"Accept" : {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"},
		"Connection" : {"keep-alive"},
		"Accept-Encoding": {"gzip, deflate"},
		"Accept-Language": {"zh-CN,zh;q=0.9"},
		"Upgrade-Insecure-Requests" : {"1"},
		"Cache-Control" : {"max-age=0"},
	}
	return DefaultClient
}
