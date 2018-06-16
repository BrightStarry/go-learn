package obtain

import (
	"zx/ipProxyPool/util"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"zx/ipProxyPool/config"
	"log"
	"time"
)

/**
	快代理 高匿
 */

type KuaidailiNnObtainer struct {
	*WebObtainer
}

func (this *KuaidailiNnObtainer) IncrementObtain() int {
	return kuaidailiOtain(this.Url,2)
}

func (this *KuaidailiNnObtainer) InitObtain() int {
	return kuaidailiOtain(this.Url,10)
}

func (this *KuaidailiNnObtainer) GetWebObtainer() *WebObtainer {
	return this.WebObtainer
}

/**
	快代理 普通
 */
type KuaidailiCommonObtainer struct {
	*WebObtainer
}

func (this *KuaidailiCommonObtainer) IncrementObtain() int {
	return kuaidailiOtain(this.Url,2)
}

func (this *KuaidailiCommonObtainer) InitObtain() int {
	return kuaidailiOtain(this.Url,10)
}

func (this *KuaidailiCommonObtainer) GetWebObtainer() *WebObtainer {
	return this.WebObtainer
}



func kuaidailiOtain(url string,count int) int{
	//proxyIps := make([]*config.ProxyIp,count * 14)
	var proxyIps []*config.ProxyIp
	for i := 1; i <= count; i++ {
		u := url + strconv.Itoa(i)
		doc := util.GetOfDocument(u)
		trs := doc.Find("#list > table > tbody > tr")
		trs.Each(func(j int, item *goquery.Selection) {
			tds := item.Find("td")
			ip := tds.Eq(0).Text()
			port := tds.Eq(1).Text()
			url, err := util.ParseToUrlOfHttp(ip, port)
			if err != nil {
				log.Println("url:", u, " 构造ipPort异常，当前值:", ip,port)
				return
			}
			proxyIp := &config.ProxyIp{
				Url:      url,
				Type:     config.Anonymity,
			}
			proxyIps = append(proxyIps,proxyIp)
		})
		// 该网站短时间访问过于频繁会返回-10
		time.Sleep(4 * time.Second)
	}
	util.AsyncProxyIpsToChan(config.WaitVerifyChan, proxyIps...)
	return len(proxyIps)
}

