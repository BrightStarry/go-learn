package verify

import (
	"zx/ipProxyPool/config"
	"github.com/pkg/errors"
	"zx/ipProxyPool/util"
	"time"
	"log"
	"fmt"
)

/**
	校验器
 */
const(
)

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
	{url:"https://example.com/",selector:"head > title",value:"Example Domain"},
	//{url:"https://www.baidu.com/",selector:"head > title",value:"百度一下，你就知道"},
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
	// 通过格式
	level1
	// 通过http
	level2
	// 通过https
	level3
	// 通过翻墙
	level4
)

 var verifyMethods = [...]func(*config.ProxyIp)error{
	 verifyFormat,
	 verifyHttp,
	 verifyHttps,
	 verifyJump,
 }


 /**
 	启动检验器
  */
func StartVerifier() {
	for v:= range config.ObtainerOutChan{
		fmt.Println(v)
		go verify(v)
	}
}

/**
	校验
 */
func verify(proxyIp *config.ProxyIp) {
	log.Println("正在校验---")
	var l = level0
	defer func() {
		if p := recover(); p != nil{
			if l > level2 {
				log.Println("校验成功---")
				config.VerifiedChan <- proxyIp
			}
		}
	}()
	for _,f := range  verifyMethods{
		if err := f(proxyIp);err != nil {
			break
		}
		l++
	}
	if l > level2 {
		log.Println("校验成功---")
		config.VerifiedChan <- proxyIp
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
	response :=util.GetByProxy(httpsVerifyUrls[0].url,proxyIp.Url)
	document := util.ResponseToDocument(response)
	value:= util.GetTextBySelector(document, httpsVerifyUrls[0].selector)
	if value != httpsVerifyUrls[0].value{
		return errors.New("https校验失败")
	}
	proxyIp.Protocol = config.HttpsFlag
	return nil
}

/**
	翻墙校验
 */
func verifyJump(proxyIp *config.ProxyIp) error{
	response :=util.GetByProxy(jumpVerifyUrls[0].url,proxyIp.Url)
	document := util.ResponseToDocument(response)
	value:= util.GetTextBySelector(document, jumpVerifyUrls[0].selector)
	if value != jumpVerifyUrls[0].value{
		return errors.New("翻墙校验失败")
	}
	proxyIp.IsJump = true
	return nil
}

