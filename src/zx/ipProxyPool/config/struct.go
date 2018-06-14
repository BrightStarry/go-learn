package config

import (
	"time"
	"fmt"
	"net/url"
)

/**
	基本结构体
 */

/**
	初始化网站信息
 */
func InitWebInfos() {
	a := NewDefaultWebInfo("http://www.xicidaili.com/nn/", "西刺国内高匿", 0, 5*time.Minute)

	*WebInfos = append(*WebInfos, *a)
}

/**
	目标网站信息
 */
type WebInfo struct {
	// 网站名-作日志打印
	Name string
	// 网址
	Url string
	// 爬取间隔
	Interval time.Duration
	// 最后标记-作增量使用
	LastLabel interface{}
	// 权重
	Weight uint8
}

func (this *WebInfo) String() string {
	return fmt.Sprintf("权重:%v,名称:%v,网址:%v,间隔:%v,最后标记:%v", this.Weight, this.Name, this.Url, this.Interval, this.LastLabel)
}

/**
	创建默认网站对象
	url:网址
	name:网站名
	interval:爬取间隔
	weight: 权重
 */
func NewDefaultWebInfo(url string, name string, weight uint8, interval time.Duration) *WebInfo {
	return &WebInfo{
		Url:      url,
		Name:     name,
		Interval: interval,
		Weight:   weight,
	}
}

// http/https
const(
	Http = 0
	Https = 1
)

// 匿名级别
const(
	// 普通
	Normal = 0
	// 匿名
	Anonymity = 1
)

/**
	代理ip
 */
 type ProxyIp struct {
 	// ip-port- https/http
 	Url *url.URL
	// 标识, http或https
	Flag uint8
	// 最后验证时间
	LastVerifyTime time.Time
	// 最后验证延迟毫秒数
	DelayMs int
	// 是否可翻墙
	IsJump bool
	// 类型,普通:0  匿名:1
	Type uint8
 }
