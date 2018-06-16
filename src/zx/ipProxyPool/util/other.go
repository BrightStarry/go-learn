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

func Sort(arr []*config.ProxyIp, size int) []*config.ProxyIp {
	in := size / 2 // 起始增量,为一半元素
	// 循环到增量为0 (1/2=0)
	for in >= 1 {
		// 循环 in次，循环所有子数组
		for i1:=0; i1<in; i1++{
			// 对每个子数组进行插入排序
			// 增量为in时， 0，0+in,0+in+in为一个子数组， 1,1+in,1+in+in是一个子数组
			for i2:= i1+in;i2 < size;i2+=in {
				temp := arr[i2]
				i3 := i2 - in
				for ;i3 >= i1 && arr[i3].Delay > temp.Delay; i3 -= in {
					arr[i3+in] = arr[i3]
				}
				arr[i3+in] =temp
			}
		}
		in = in / 2
	}
	return arr
}

