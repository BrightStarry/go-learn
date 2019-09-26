package config

import (
	"net/http"
	"time"
	"sync"
	"fmt"
	"net/url"
	"net/http/cookiejar"
	"log"
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
func Init() {

	// 初始化请求头
	InitDefaultHeader()
	// 初始化client
	InitDefaultClient()
	// 初始化校验器client pool
	InitVerifierClient()

	WaitVerifyChan = make(chan *ProxyIp, Config.VerifierThreadNum*100)
	VerifiedChan = make(chan *ProxyIp, 512)
	ReVerifyChan = make(chan *ProxyIp, Config.ReVerifyThreadNum*100)
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
	WebPort string
	// 爬虫超时时间
	SpiderTimeout time.Duration
	// 校验器超时时间
	VerifierTimeout time.Duration
	// 校验器并发数
	VerifierThreadNum int
	// ip重校验间隔
	ReVerifyInterval time.Duration
	// ip重校验线程数
	ReVerifyThreadNum int
}

func (this *SystemConfig) String() string {
	return fmt.Sprintln("web服务端口:",this.WebPort,",爬虫超时时间:",this.SpiderTimeout,",校验器超时时间:",this.VerifierTimeout,
		",校验器并发数:",this.VerifierThreadNum,",ip重校验间隔:",this.ReVerifyInterval,",ip重校验线程数:",this.ReVerifyThreadNum)
}

/**
	初始化系统参数
 */
func InitSystemConfig() {
	Config = &SystemConfig{
		// web服务端口
		WebPort: "9999",
		// 爬虫超时时间
		SpiderTimeout: 15 * time.Second,
		// 校验器超时时间
		VerifierTimeout: 5 * time.Second,
		// 校验器并发数
		VerifierThreadNum: 128,
		// 入库ip重新校验间隔
		ReVerifyInterval: 5 * time.Minute,
		// ip重校验线程数
		ReVerifyThreadNum: 16,
	}
}

/**
	初始化client
 */
func InitDefaultClient() {
	DefaultClient = &http.Client{}
	// 超时时间
	DefaultClient.Timeout = Config.SpiderTimeout
	// 构建cookie
	cookie, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln("构建cookie失败:", err)
	}
	DefaultClient.Jar = cookie
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
    Windows

Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0
Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.835.163 Safari/535.1
Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)
 */
func InitDefaultHeader() {
	// 设置默认请求头
	DefaultHeader = &map[string][]string{
		"User-Agent":                {"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"},
		"Connection":                {"keep-alive"},
		"Accept-Encoding":           {"gzip, deflate"},
		"Accept-Language":           {"zh-CN,zh;q=0.9"},
		"Upgrade-Insecure-Requests": {"1"},
		"Cache-Control":             {"max-age=0"},
	}
}

// http/https
const(
	HttpFlag = 0
	HttpsFlag = 1
	SocksFlag = 2
	Http = "http"
	Https = "https"
)


// 匿名级别
const(
	// 普通
	Normal = 0
	// 匿名
	Anonymity = 1
)

// 来自-所有要爬取的网站
const(
	Xici = iota
	Ip66
	Kuaidaili
	Ip3366
	Ip89
)

//其他
const(
	HttpStatusSpiderBad = 521
)

/**
	代理ip
 */
type ProxyIp struct {
	// ip-port- https/http
	Url *url.URL
	// 协议, http/ https/ socks4/5
	Protocol uint8
	// 最后验证时间
	LastVerifyTime time.Time
	// 最后验证延迟毫秒数
	Delay time.Duration
	// 是否可翻墙
	IsJump bool
	// 类型,普通:0  匿名:1
	Type uint8
	// 来自哪个网站
	From uint8
}

func (this *ProxyIp) String() string {
	return fmt.Sprintf("host:%v,协议:%v,延迟:%v,翻墙:%v,最后验证时间：%v",this.Url.Host,this.Protocol,this.Delay.String(),this.IsJump,this.LastVerifyTime)
}


/**
   ip返回对象
*/
type IpDTO struct{
	// [ip]:[port]
	Host string `json:"host"`
	// 协议 http or https
	Protocol string `json:"protocol"`
	// 最后验证时间
	LastVerifyTime time.Time `json:"lastVerifyTime"`
	// 最后验证延迟毫秒数
	Delay float64 `json:"delay"`
	// 是否可翻墙
	IsJump bool `json:"isJump"`
}

