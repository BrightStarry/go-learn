package main

import (
	"zx/ipProxyPool/verify"
	"log"
	"time"
	"zx/ipProxyPool/obtain"
	"zx/ipProxyPool/store"
	"zx/ipProxyPool/web"
)

/**
	启动
 */
func main() {
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

	web.SyncStartWebServer()
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
			asyncIncrementTicker(&val, &c)
		}(v)
	}
}

/**
	增量获取定时器
 */
func asyncIncrementTicker(o *obtain.Obtainer, c *<-chan time.Time) {
	// 每当收到信号
	go func() {
		for tick := range *c {
			func() {
				// 进行异常捕获
				defer func() {
					if err := recover(); err != nil {
						log.Println((*o).GetWebObtainer().Name, " 增量获取失败:", err)
					}
				}()
				log.Println((*o).GetWebObtainer().Name, " 进行增量获取:", tick)
				length := (*o).IncrementObtain()
				log.Println((*o).GetWebObtainer().Name, " 增量获取数:", length)
			}()
		}
	}()
}
