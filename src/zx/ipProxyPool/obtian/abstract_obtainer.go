package obtian

import "zx/ipProxyPool/config"

/**
	获取者接口
 */
type Obtainer interface {
	// 设置 webInfo
	SetWebInfo(webInfo *config.WebInfo)
	// 初始获取全部ip方法
	InitObtain() *[]config.ProxyIp
	// 增量获取ip方法
	IncrementObtain() *[]config.ProxyIp
}

/**
	获取者通用参数
 */
 type BaseObtainer struct{
 	// 网站信息
 	WebInfo *config.WebInfo
 }
