package config

import (
	"net/http"
	"time"
)

/**
	参数相关
 */

// 系统配置
var Config *SystemConfig
// 默认请求头
var DefaultHeader *map[string][]string
// 爬虫默认client
var DefaultClient *http.Client
// 校验器client
var VerifierClient *http.Client
// obtainer 入库 通道
var ObtainerOutChan = make(chan *ProxyIp, 128)
// 校验通过 通道
var VerifiedChan = make(chan *ProxyIp, 128)

/**
   初始化方法
*/
func init() {
	// 初始化系统参数
	InitSystemConfig()
	// 初始化请求头
	InitDefaultHeader()
	// 初始化client
	InitDefaultClient()
}

/**
	系统参数
 */
type SystemConfig struct {
	// 默认超时时间
	SpiderTimeout time.Duration
	// 校验器超时时间
	VerifierTimeout time.Duration
}

/**
	初始化系统参数
 */
func InitSystemConfig() {
	Config = &SystemConfig{
		// 爬虫超时时间
		SpiderTimeout: 5 * time.Second,
		// 校验器超时时间
		VerifierTimeout: 5 * time.Second,
	}
}

/**
	初始化client
 */
func InitDefaultClient() {
	DefaultClient = &http.Client{}
	// 超时时间
	DefaultClient.Timeout = Config.SpiderTimeout
	// 构建cookie  暂不构建cookie
	//cookie, err := cookiejar.New(nil)
	//if err != nil {
	//	log.Fatalln("构建cookie失败:", err)
	//}
	//DefaultClient.Jar = cookie
}

/**
	初始化校验器client
 */
func InitVerifierClient() {
	VerifierClient = &http.Client{}
	// 超时时间
	DefaultClient.Timeout = Config.SpiderTimeout
}

/**
	初始化默认请求头
 */
func InitDefaultHeader() {
	// 设置默认请求头
	DefaultHeader = &map[string][]string{
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"},
		"Connection":                {"keep-alive"},
		"Accept-Encoding":           {"gzip, deflate"},
		"Accept-Language":           {"zh-CN,zh;q=0.9"},
		"Upgrade-Insecure-Requests": {"1"},
		"Cache-Control":             {"max-age=0"},
	}
}
