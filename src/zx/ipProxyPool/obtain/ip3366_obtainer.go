package obtain

import (
	"zx/ipProxyPool/config"
	"fmt"
	"zx/ipProxyPool/util"
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

/**
	ip3366 云代理
 */
type Ip3366Obtainer struct {
	*WebObtainer
}

func (this *Ip3366Obtainer) IncrementObtain() int {
	//增量获取1
	return ip3366Obtain(this.Url, 1)
}

func (this *Ip3366Obtainer) InitObtain() int {
	// 初始化时 10页
	return ip3366Obtain(this.Url, 10)
}
func (this *Ip3366Obtainer) GetWebObtainer() *WebObtainer {
	return this.WebObtainer
}

func ip3366Obtain(url string, count int) int{

	var proxyIps []*config.ProxyIp
	// 该网站由于需要一次性访问过多数据,并且经常超时,所以捕获异常
	defer func() {
		// 发送数据
		util.AsyncProxyIpsToChan(config.WaitVerifyChan, proxyIps...)
		if err:= recover();err != nil{
			// 再次抛出
			panic(err)
		}
	}()

	// 有四种类型
	for i := 1; i <= 4; i++ {
		// 每种类型有x页
		for j:=1; j <= count;j++{
			u := fmt.Sprintf(url,i,j)
			document := util.GetOfDocument(u)
			// 获取表格元素
			trs := document.Find("#list > table > tbody > tr")
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
					Url:url,
					From:config.Ip3366,
				}
				proxyIps = append(proxyIps, proxyIp)
			})
			time.Sleep(5 * time.Second)
		}
	}
	util.AsyncProxyIpsToChan(config.WaitVerifyChan, proxyIps...)
	return len(proxyIps)
}
