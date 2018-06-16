package config

import (
	"time"
	"net/url"
	"fmt"
)

/**
	基本结构体
 */




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
	Delay time.Duration `json:"delay"`
	// 是否可翻墙
	IsJump bool `json:"isJump"`
}
