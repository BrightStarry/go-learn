package obtian

import (
	"zx/ipProxyPool/config"
	"log"
	"github.com/PuerkitoBio/goquery"
)

/**
	西刺高匿
 */
type XiciAnonymity struct {
	BaseObtainer
}

func (this *XiciAnonymity) SetWebInfo(webInfo *config.WebInfo) {
	this.WebInfo = webInfo
}

func (this *XiciAnonymity) InitObtain() *[]config.ProxyIp{

	u := this.WebInfo.Url + "1"
	response, err := config.DefaultClient.Get(u)
	if err != nil {
		log.Println(this.WebInfo.Name,"-访问:",u,"-异常:",err)
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println(this.WebInfo.Name,"-读取响应异常:",err)
	}
	tbody := document.Find("#ip_list > tbody")
	log.Println(tbody)
	return nil
}