package obtain

import (
	"testing"
	"zx/ipProxyPool/config"
)


/**
	 测试
 */

func TestXiciAnonymity_InitObtain(t *testing.T) {
	config.InitSystemConfig()
	config.Init()

	//
	//request,_ := http.NewRequest(http.MethodGet,"https://www.aicoin.net.cn/currencies",nil)
	//request.Header = *config.DefaultHeader
	//request.Header.Add("Host","www.aicoin.net.cn")
	//request.Header.Add("Referer","https://www.aicoin.net.cn/")
	//response, _ := config.DefaultClient.Do(request)
	//
	//
	//
	//time.Sleep(10  * time.Second)
	//
	//request,_ = http.NewRequest(http.MethodGet,"https://www.aicoin.net.cn/api/ping",nil)
	//request.Header = *config.DefaultHeader
	//request.Header.Add("Host","www.aicoin.net.cn")
	//request.Header.Add("Referer","https://www.aicoin.net.cn/")
	//u,_ := url.Parse("aicoin.net.cn")
	//c:=config.DefaultClient.Jar.Cookies(u)
	//for _,v
	//request.Header.Add("X-XSRF-TOKEN","")
	//response,_ = config.DefaultClient.Do(request)
	//result := util.ResponseToStr(response)
	//fmt.Println(result)
	//
	//fmt.Println("---------------------------------------------------------------------------------------")
	//
	//time.Sleep(10  * time.Second)
	//response = util.Get("https://www.aicoin.net.cn/api/data/getMc?type=all")
	//result = util.ResponseToStr(response)
	//fmt.Println(result)
	//
	//
	//time.Sleep(time.Hour)

}
