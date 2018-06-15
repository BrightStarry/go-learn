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

func (this *Ip66CommonObtainer) IncrementObtain() {
	//增量获取300
	ip66Obtain(this.Url, 300)
}

func (this *Ip66CommonObtainer) InitObtain() {
	// 初始化时，提取1001
	ip66Obtain(this.Url, 1001)
}



/**
	66 https
 */
type Ip66HttpsObtainer struct {
	*WebObtainer
}

func (this *Ip66HttpsObtainer) IncrementObtain() {
	//增量获取200
	ip66Obtain(this.Url, 200)
}

func (this *Ip66HttpsObtainer) InitObtain() {
	// 初始化时，提取1001
	ip66Obtain(this.Url, 200)
}

/**
	通用提取方法
 */
func ip66Obtain(url string, initSum int) {
	u := url + strconv.Itoa(initSum)
	response := util.Get(u)
	doc := util.ResponseToDocument(response)
	bodyEle := doc.Find("body")
	body, _ := bodyEle.Html()
	arr := strings.Split(body, "<br/>\n\t\t")
	proxyIps := make([]*config.ProxyIp, initSum-1)
	//遍历,舍弃最后一个ip，因为还需要切割
	for i := 0; i < len(arr)-1; i++ {
		ip, port, err := net.SplitHostPort(arr[i])
		if err != nil {
			log.Println("url:", u, " 解析ioPort异常,当前值:", arr[i])
			return
		}
		url, err := util.ParseToUrlOfHttp(ip, port)
		if err != nil {
			log.Println("url:", u, " 构造ipPort异常，当前值", arr[i])
			return
		}
		proxyIps = append(proxyIps,&config.ProxyIp{
			Url:  url,
			//Protocol:config.HttpFlag,
			Type: config.Normal,
		})
	}
	util.AsyncProxyIpsToChan(config.ObtainerOutChan, proxyIps...)
}

