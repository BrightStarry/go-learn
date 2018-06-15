package util

import (
	"zx/ipProxyPool/config"
)

/**
	异步将任意数量得 proxyIp 放入 chan
 */
func AsyncProxyIpsToChan(c chan *config.ProxyIp,ip ...*config.ProxyIp) {
	go func() {
		for _,v := range ip{
			if v == nil {
				continue
			}
			c <- v
		}
	}()
}

