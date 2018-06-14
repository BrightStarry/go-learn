package obtian

import (
	"zx/ipProxyPool/config"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"fmt"
	"compress/gzip"
	"io"
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

	u1, _ := url.Parse("xicidaili.com")
	config.DefaultClient.Jar.SetCookies(u1,[]*http.Cookie{{Name:"_free_proxy_session",Value:"BAh7B0kiD3Nlc3Npb25faWQGOgZFVEkiJTJlYmJiODlkMGJmZDc1YTdkZWJlYThkMWIwYjc0YmI0BjsAVEkiEF9jc3JmX3Rva2VuBjsARkkiMXpZQVkrdDhpamNFYzl2b0MvWWFQUHZPbEp5bXBGWmZwcUsxM0dqaWdnQW89BjsARg%3D%3D--7a3ee1f1f7ca29b6216d01bfbc08f2f75bba5177"}})

	u := this.WebInfo.Url + "1"
	request,err := http.NewRequest(http.MethodGet,u,nil)
	if err != nil {
		log.Println(this.WebInfo.Name,"-构建请求:",u,"-异常:",err)
	}
	request.Header = *config.DefaultHeader
	request.Header.Add("Referer", "http://www.xicidaili.com/nt/")
	request.Header.Add("Host", "www.xicidaili.com")
	response, err := config.DefaultClient.Do(request)

	if err != nil {
		log.Println(this.WebInfo.Name,"-访问:",u,"-异常:",err)
	}
	defer response.Body.Close()

	var reader *io.ReadCloser
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			
		}
	} else {
		reader = &response.Body
	}

	//document, err := goquery.NewDocumentFromReader(response.Body)
	//if err != nil {
	//	log.Println(this.WebInfo.Name,"-读取响应异常:",err)
	//}
	//totalPageEle := document.Find("#body > div.pagination > a:nth-child(13)")
	//fmt.Println(document.Html())
	//log.Println(totalPageEle.Text())
	b,_ :=ioutil.ReadAll(reader)
	fmt.Println(string(b))
	return nil
}