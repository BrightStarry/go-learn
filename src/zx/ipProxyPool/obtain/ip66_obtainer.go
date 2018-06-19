package obtain

import (
	"zx/ipProxyPool/config"
	"zx/ipProxyPool/util"
	"strings"
	"net"
	"log"
	"strconv"
)

/**
	66普通
	http://www.66ip.cn/mo.php?tqsl=100
 */

type Ip66CommonObtainer struct {
	*WebObtainer
}

func (this *Ip66CommonObtainer) IncrementObtain() int {
	//增量获取300
	return ip66Obtain(this.Url, 301)
}

func (this *Ip66CommonObtainer) InitObtain() int {
	// 初始化时
	return ip66Obtain(this.Url, 501)
}
func (this *Ip66CommonObtainer) GetWebObtainer() *WebObtainer {
	return this.WebObtainer
}


/**
	66 https
 */
type Ip66HttpsObtainer struct {
	*WebObtainer
}

func (this *Ip66HttpsObtainer) IncrementObtain() int {
	//增量获取200
	return ip66Obtain(this.Url, 201)
}

func (this *Ip66HttpsObtainer) InitObtain() int {
	// 初始化时，提取1001
	return ip66Obtain(this.Url, 301)
}

func (this *Ip66HttpsObtainer) GetWebObtainer() *WebObtainer {
	return this.WebObtainer
}

/**
	通用提取方法
 */
func ip66Obtain(url string, initSum int) int{
	u := url + strconv.Itoa(initSum)
	doc := util.GetOfDocument(u)
	bodyEle := doc.Find("body")
	body, _ := bodyEle.Html()
	arr := strings.Split(body, "<br/>\n\t\t")
	var proxyIps []*config.ProxyIp
	//遍历,舍弃最后一个ip，因为还需要切割
	for i := 0; i < len(arr)-1; i++ {
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
			From:config.Ip66,
		})
	}
	util.AsyncProxyIpsToChan(config.WaitVerifyChan, proxyIps...)
	return len(proxyIps)
}

