package verify

import (
	"zx/ipProxyPool/config"
	"github.com/pkg/errors"
	"zx/ipProxyPool/util"
	"time"
)

/**
	校验器
 */

/**
	校验网站
 */
type verifyUrl struct{
	// 网址
	url string
	// selector,指向网站上一个比较醒目的元素
	selector string
	// 该元素的值
	value string
}

/**
	http校验网站
 */
var httpVerifyUrls = []verifyUrl{
	{url:"http://example.com/",selector:"head > title",value:"Example Domain"},
}

/**
	https校验网站
 */
var httpsVerifyUrls = []verifyUrl{
	//{url:"https://example.com/",selector:"head > title",value:"Example Domain"},
	{url:"https://www.baidu.com/",selector:"head > title",value:"百度一下，你就知道"},
}

/**
	翻墙校验网站
 */
 var jumpVerifyUrls = []verifyUrl{
	 {url:"https://www.google.com/",selector:"head > title",value:"Google"},
 }

const(
	// 默认
	level0 = iota
	// 通过格式校验
	level1
	// 通过http, 取消http校验
	//level2
	// 通过https
	level3
	// 通过翻墙
	level4
)

 /**
 	校验方法数组
  */
 var verifyMethods = [...]func(*config.ProxyIp)error{
	 verifyFormat,
	 //verifyHttp,
	 verifyHttps,
	 verifyJump,
 }


 /**
 	启动检验器
  */
func StartVerifier() {
	// 限制并发
	for i := 0; i < config.Config.VerifierThreadNum; i++{
		asyncVerify()
	}
}

/**
	异步校验
 */
func asyncVerify() {
	go func() {
		for v := range config.WaitVerifyChan {
			d := Distinct(v)
			if d {
				Verify(v)
			}

		}
	}()
}

/**
	去重
	return true:表示不存在, false:表示重复
 */
func Distinct(ip *config.ProxyIp) bool {
	if _, ok := config.ProxyIpDistinctMap.Load(ip.Url.Host); ok{
		return false
	}
	config.ProxyIpDistinctMap.Store(ip.Url.Host,nil)
	return true
}

/**
	校验
 */
func Verify(proxyIp *config.ProxyIp) {
	var l = level0
	defer func() {
		// 此处捕获异常,不做任何处理
		if p := recover(); p != nil{
			//fmt.Println("校验异常:",p)
		}
		// 实现了https,加入通道
		if l >= level3 {
			config.VerifiedChan <- proxyIp
		}
	}()
	// 循环使用校验方法校验
	for _,f := range  verifyMethods{
		if err := f(proxyIp);err != nil {
			return
		}
		l++
	}

}

/**
	格式校验
 */
func verifyFormat(proxyIp *config.ProxyIp) error{
	if proxyIp.Url.Host == "" || proxyIp.Url.Port() == "" {
		return errors.New("ip或port为空")
	}
	return nil
}

/**
	http校验
 */
func verifyHttp(proxyIp *config.ProxyIp) error{
	start := time.Now()
	response :=util.GetByProxy(httpVerifyUrls[0].url,proxyIp.Url)
	elapsed := time.Since(start)
	document := util.ResponseToDocument(response)
	value:= util.GetTextBySelector(document, httpVerifyUrls[0].selector)
	if value != httpVerifyUrls[0].value{
		return errors.New("http校验失败")
	}
	// 记录验证时间和延迟
	proxyIp.LastVerifyTime = start
	proxyIp.Delay = elapsed
	proxyIp.Protocol = config.HttpFlag
	return nil
}

/**
	https校验
 */
func verifyHttps(proxyIp *config.ProxyIp)error {
	start := time.Now()
	response :=util.GetByProxy(httpsVerifyUrls[0].url,proxyIp.Url)
	elapsed := time.Since(start)
	document := util.ResponseToDocument(response)
	value:= util.GetTextBySelector(document, httpsVerifyUrls[0].selector)
	if value != httpsVerifyUrls[0].value{
		return errors.New("https校验失败")
	}
	// 记录验证时间和延迟
	proxyIp.LastVerifyTime = start
	proxyIp.Delay = elapsed
	proxyIp.Protocol = config.HttpsFlag
	return nil
}

/**
	翻墙校验
 */
func verifyJump(proxyIp *config.ProxyIp) error{
	start := time.Now()
	response :=util.GetByProxy(jumpVerifyUrls[0].url,proxyIp.Url)
	elapsed := time.Since(start)
	document := util.ResponseToDocument(response)
	value:= util.GetTextBySelector(document, jumpVerifyUrls[0].selector)
	if value != jumpVerifyUrls[0].value{
		return errors.New("翻墙校验失败")
	}
	// 记录验证时间和延迟
	proxyIp.LastVerifyTime = start
	proxyIp.Delay = elapsed
	proxyIp.IsJump = true
	return nil
}

