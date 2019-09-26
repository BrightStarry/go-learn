package main

import (
	"zx/ipProxyPool/verify"
	"log"
	"time"
	"zx/ipProxyPool/obtain"
	"zx/ipProxyPool/store"
	"zx/ipProxyPool/web"
	"zx/ipProxyPool/config"
	"flag"
)

/**
	启动
 */
func main() {
	// 初始化系统参数
	config.InitSystemConfig()

	// 获取外部参数,修改默认系统参数
	importExtParam()

	log.Println("参数:",config.Config)

	// 配置初始化
	config.Init()

	// 启动校验器
	verify.StartVerifier()

	// 启动存储器
	store.StartStorage()

	// 启动定时校验任务
	store.StartVerifyTicker()

	// 初始化各代理,获取ip
	asyncInitObtainer()

	// 启动增量获取定时器
	startObtainerIncrementTicker()

	// 启动web服务
	web.SyncStartWebServer()
}

/**
	获取外部参数
 */
func importExtParam() {

	// 端口
	flag.StringVar(&config.Config.WebPort, "port", config.Config.WebPort, "服务器端口")
	// 爬虫超时时间
	var spiderTimeoutSecond int
	flag.IntVar(&spiderTimeoutSecond, "stt", int(config.Config.SpiderTimeout.Seconds()), "爬虫超时时间(单位:s)")


	// 校验器超时时间
	var verifyTimeoutSecond int
	flag.IntVar(&verifyTimeoutSecond, "vtt", int(config.Config.VerifierTimeout.Seconds()), "ip校验超时时间(单位:s)")

	// ip校验并发数
	flag.IntVar(&config.Config.VerifierThreadNum,"vc",config.Config.VerifierThreadNum,"ip校验并发数")

	// ip重校验间隔
	var reVerifyIntervalMinute int
	flag.IntVar(&reVerifyIntervalMinute,"rvi",int(config.Config.ReVerifyInterval.Minutes()),"ip重校验间隔(每次重校验3分之一,单位:min)")

	// ip重校验并发数
	flag.IntVar(&config.Config.ReVerifyThreadNum,"rvc",config.Config.ReVerifyThreadNum,"ip重校验并发数")

	flag.Parse()

	config.Config.SpiderTimeout = time.Duration(spiderTimeoutSecond) * time.Second
	config.Config.VerifierTimeout = time.Duration(verifyTimeoutSecond) * time.Second
	config.Config.ReVerifyInterval = time.Duration(reVerifyIntervalMinute) * time.Minute



}

/**
	执行各ip代理初始化方法
 */
func asyncInitObtainer() {
	for _, v := range obtain.WebObtainers {
		// 将初始化方法放在内部函数中
		// 只有这样值传递,才不会被后续循环影响
		go func(val obtain.Obtainer) {
			// 进行异常捕获
			defer func() {
				if err := recover(); err != nil {
					log.Println(val.GetWebObtainer().Name, " 初始化获取失败:", err)
				}
			}()
			// 执行初始化
			log.Println(val.GetWebObtainer().Name, " 进行初始化获取")
			length := val.InitObtain()
			log.Println(val.GetWebObtainer().Name, " 初始化获取数:", length)
		}(v)
	}
}

/**
	启动各obtainer的增量获取定时器
 */
func startObtainerIncrementTicker() {
	for _, v := range obtain.WebObtainers {
		func(val obtain.Obtainer) {
			// 创建定时通道
			c := time.Tick(val.GetWebObtainer().Interval)
			// 启动定时器
			asyncIncrementTicker(val, &c)
		}(v)
	}
}

/**
	增量获取定时器
 */
func asyncIncrementTicker(o obtain.Obtainer, c *<-chan time.Time) {
	// 每当收到信号
	go func() {
		for  range *c {
			func() {
				// 进行异常捕获
				defer func() {
					if err := recover(); err != nil {
						log.Println(o.GetWebObtainer().Name, " 增量获取失败:", err)
					}
				}()
				log.Println(o.GetWebObtainer().Name, " 进行增量获取")
				length := o.IncrementObtain()
				log.Println(o.GetWebObtainer().Name, " 增量获取数:", length)
			}()
		}
	}()
}
