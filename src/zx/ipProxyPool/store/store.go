package store

import (
	"zx/ipProxyPool/config"
	"log"
	"time"
	"zx/ipProxyPool/util"
	"zx/ipProxyPool/verify"
)

/**
	存储管理ip
 */

/**
 	启动存储器
 */
func StartStorage() {
	go func() {
		for ip := range config.VerifiedChan {
			log.Println("入库:", ip)
			put(ip)
		}
	}()
}

/**
	启动定时校验任务
 */
func StartVerifyTicker() {

	ticker := time.Tick(config.Config.ReVerifyInterval)
	go func() {
		for t := range ticker {
			log.Println(t, " 启动定时校验任务,当前待校验通道长度:", len(config.ReVerifyChan))
			verifyTask()
		}
	}()
}

/**
	校验任务
 */
func verifyTask() {
	// 获取到前三分之一的ip
	ips := getAndDelIpsAtFirst(3)
	// 写入待校验chan
	util.AsyncProxyIpsToChan(config.ReVerifyChan, ips...)
	for i := 0; i < config.Config.ReVerifyThreadNum; i++ {
		go func() {
			for ip := range config.ReVerifyChan {
				verify.Verify(ip)
			}
		}()
	}
}

/**
	加入队列
 */
func put(ip *config.ProxyIp) {
	config.ProxyIpStore.Lock.Lock()
	defer config.ProxyIpStore.Lock.Unlock()
	config.ProxyIpStore.Queue = append(config.ProxyIpStore.Queue, ip)
}

/**
	从队首获取并删除指定长度的ip
 */
func getAndDelIpsAtFirst(x int) ([]*config.ProxyIp) {
	// 写锁
	config.ProxyIpStore.Lock.Lock()
	defer config.ProxyIpStore.Lock.Unlock()
	l := len(config.ProxyIpStore.Queue)
	temp := config.ProxyIpStore.Queue[:l/x]
	config.ProxyIpStore.Queue = config.ProxyIpStore.Queue[l/x:]
	return temp
}

/**
 	从队尾开始 获取指定长度的ip
 */
func GetIpsAtLast(size int) []*config.ProxyIp{
	// 读锁
	config.ProxyIpStore.Lock.RLock()
	defer config.ProxyIpStore.Lock.RUnlock()

	// 长度越界,返回全部
	length := len(config.ProxyIpStore.Queue)
	if size > length {
		size = length
	}
	if length == 0 {
		return []*config.ProxyIp{}
	}
	return config.ProxyIpStore.Queue[length - size:]
}
