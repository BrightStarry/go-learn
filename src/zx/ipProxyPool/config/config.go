package config

import (
	"net/http"
	"time"
	"sync"
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
// 校验器client池(用池的原因是为了减少gc,复用client)
var VerifierClientPool *sync.Pool
// 待校验通道
var WaitVerifyChan chan *ProxyIp
// 校验通过 通道
var VerifiedChan chan *ProxyIp
// 重新校验通道
var ReVerifyChan chan *ProxyIp
// ip去重map
var ProxyIpDistinctMap *sync.Map
// ip存储
var ProxyIpStore = struct {
	Queue []*ProxyIp
	Lock  sync.RWMutex
}{Lock: sync.RWMutex{}}

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
	// 初始化校验器client pool
	InitVerifierClient()

	WaitVerifyChan = make(chan *ProxyIp, Config.VerifierThreadNum*10)
	VerifiedChan = make(chan *ProxyIp, Config.VerifiedChanBufNum)
	ReVerifyChan = make(chan *ProxyIp, Config.ReVerifyThreadNum*10)
	ProxyIpDistinctMap = &sync.Map{}
}

/**
	监控参数
 */
type MonitorParam struct {
	// ip池长度
	IpPoolLen int `json:"ipPoolLen"`
	// 待校验通道长度
	WaitVerifyChanLen int `json:"waitVerifyChanLen"`
	// 校验通过通道长度 该通道就一个去重.不需要
	//VerifiedChanLen int `json:"verifiedChanLen"`
	// 重新校验通道长度
	ReVerifyChanLen int `json:"reVerifyChanLen"`
}

/**
	系统参数
 */
type SystemConfig struct {
	// web服务端口
	WebPort int
	// 默认超时时间
	SpiderTimeout time.Duration
	// 校验器超时时间
	VerifierTimeout time.Duration
	// 校验器并发数
	VerifierThreadNum int
	// 校验通过通道缓冲数
	VerifiedChanBufNum int
	// ip重校验间隔
	ReVerifyInterval time.Duration
	// ip重校验线程数
	ReVerifyThreadNum int
}

/**
	初始化系统参数
 */
func InitSystemConfig() {
	Config = &SystemConfig{
		// web服务端口
		WebPort: 8080,
		// 爬虫超时时间
		SpiderTimeout: 15 * time.Second,
		// 校验器超时时间
		VerifierTimeout: 10 * time.Second,
		// 校验器并发数
		VerifierThreadNum: 64,
		// 校验通过通道缓冲数
		VerifiedChanBufNum: 32,
		// 入库ip重新校验间隔
		ReVerifyInterval: 10 * time.Minute,
		// ip重校验线程数
		ReVerifyThreadNum: 32,
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
	// 创建对象池
	VerifierClientPool = &sync.Pool{New: func() interface{} {
		return &http.Client{
			Transport: &http.Transport{},
			Timeout:   Config.VerifierTimeout,
		}
	}}
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
