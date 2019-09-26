package obtain

import (
	"strconv"
	"zx/ipProxyPool/util"
	"strings"
	"zx/ipProxyPool/config"
	"net"
	"log"
)

/**
	89ip
 */
type Ip89Obtainer struct {
	*WebObtainer
}

func (this *Ip89Obtainer) IncrementObtain() int {
	//增量获取300
	return ip89Obtain(this.Url, 301)
}

func (this *Ip89Obtainer) InitObtain() int {
	// 初始化时
	return ip89Obtain(this.Url, 501)
}
func (this *Ip89Obtainer) GetWebObtainer() *WebObtainer {
	return this.WebObtainer
}

func ip89Obtain(url string, initSum int) int{
	u := url + strconv.Itoa(initSum)
	doc := util.GetOfDocument(u)
	bodyEle := doc.Find("body")
	body, _ := bodyEle.Html()
	arr := strings.Split(body, "<br/>")
	var proxyIps []*config.ProxyIp
	//遍历,舍弃第一个ip，因为还需要切割
	for i := 2; i < len(arr); i++ {
		ip, port, err := net.SplitHostPort(arr[i])
		if err != nil {
			log.Println("url:", u, " 解析ioPort异常,当前值:", arr[i])
			return 0
		}
		url, err := util.ParseToUrlOfHttp(ip, port)
		if err != nil {
			log.Println("url:", u, " 构造ipPort异常，当前值", arr[i])
			return 0
		}
		proxyIps = append(proxyIps,&config.ProxyIp{
			Url:  url,
			Type: config.Normal,
			From: config.Ip89,
		})
	}
	util.AsyncProxyIpsToChan(config.WaitVerifyChan, proxyIps...)
	return len(proxyIps)
}