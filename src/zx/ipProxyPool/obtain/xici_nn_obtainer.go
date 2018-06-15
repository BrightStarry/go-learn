package obtain

import (
	"strconv"
	"zx/ipProxyPool/util"
	"github.com/PuerkitoBio/goquery"
	"log"
	"zx/ipProxyPool/config"
	"time"
)

/**
	西刺Http
 */
type XiciHttpObtainer struct {
	*WebObtainer
}


func (this *XiciHttpObtainer) InitObtain() {
	xiciObtain(this.Url,10,config.Anonymity)
}

func (this *XiciHttpObtainer) IncrementObtain() {
	xiciObtain(this.Url,2,config.Anonymity)
}

/**
	西刺Https
 */
type XiciHttpsObtainer struct {
	*WebObtainer
}


func (this *XiciHttpsObtainer) InitObtain() {
	xiciObtain(this.Url,10,config.Normal)
}

func (this *XiciHttpsObtainer) IncrementObtain() {
	xiciObtain(this.Url,2,config.Normal)
}

func xiciObtain(url string,count int,t uint8) {
	proxyIps := make([]*config.ProxyIp,count * 100)
	// 遍历每一页
	for i := 0; i < count; i++ {
		u := url + strconv.Itoa(i)
		document := util.GetOfDocument(u)
		// 获取表格元素
		trs := document.Find("#ip_list > tbody > tr")
		trs.Each(func(j int, item *goquery.Selection) {
			if i == 0 {
				return
			}
			tds := item.Find("td")
			ip := tds.Eq(1).Text()
			port := tds.Eq(2).Text()
			url, err := util.ParseToUrlOfHttp(ip, port)
			if err != nil {
				log.Println("url:", u, " 构造ipPort异常，当前值:", ip,port)
				return
			}

			proxyIp := &config.ProxyIp{
				Url:url,
				//Protocol:protocol,
				Type:t,
			}
			proxyIps = append(proxyIps, proxyIp)
		})
		time.Sleep(7 * time.Second)
	}
	util.AsyncProxyIpsToChan(config.ObtainerOutChan, proxyIps...)
}



/**
	获取总页数
 */
//func getTotalPage(this *XiciHttpObtainer) int {
//	u := this.Url + "1"
//	doc := util.ResponseToDocument(util.Get(u))
//	//fmt.Println(doc.Html())
//	totalPageStr := util.GetTextBySelector(doc, "#body > div.pagination > a:nth-child(13)")
//	totalPage, err := strconv.Atoi(totalPageStr)
//	if err != nil {
//		log.Fatalln(this.Name, "获取总页数失败,当前获取值:", totalPageStr)
//	}
//	return totalPage
//}
