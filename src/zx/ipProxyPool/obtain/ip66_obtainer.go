package obtain

import (
	"zx/ipProxyPool/config"
	"zx/ipProxyPool/util"
	"strings"
	"net"
	"log"
	"strconv"
	"errors"
	"github.com/robertkrimen/otto"
	"net/http"
	"net/url"
)

/**
	66普通
	http://www.66ip.cn/mo.php?tqsl=100

	ps:提取数量应该无限制
 */
type Ip66CommonObtainer struct {
	*WebObtainer
}

func (this *Ip66CommonObtainer) IncrementObtain() int {
	//增量获取300
	return ip66Obtain(this.Url, 2000)
}

func (this *Ip66CommonObtainer) InitObtain() int {
	// 初始化时
	return ip66Obtain(this.Url, 5000)
}
func (this *Ip66CommonObtainer) GetWebObtainer() *WebObtainer {
	return this.WebObtainer
}


/**
	66 匿名
	ps:限制了最多300
 */
type Ip66AnonymityObtainer struct {
	*WebObtainer
}

func (this *Ip66AnonymityObtainer) IncrementObtain() int {
	//增量获取200
	return ip66Obtain(this.Url, 300)
}

func (this *Ip66AnonymityObtainer) InitObtain() int {
	// 初始化时不提取，为了和ip66普通代理错开时间，防止写入cookie时冲突
	//return ip66Obtain(this.Url, 300)
	return 0
}

func (this *Ip66AnonymityObtainer) GetWebObtainer() *WebObtainer {
	return this.WebObtainer
}

/**
	通用提取方法
 */
func ip66Obtain(thisUrl string, initSum int) int{
	u := thisUrl + strconv.Itoa(initSum)

	request := util.BuildGet(u)
	response, err := config.DefaultClient.Do(request)
	if err != nil {
		panic(errors.New("url:" + u + " 请求异常:"+err.Error()))
	}
	doc := util.ResponseToDocument(response)
	// 处理521状态码
	if response.StatusCode == config.HttpStatusSpiderBad{
		// 获取第一段js
		script1 := doc.Find("script").Text()
		script1 = strings.Replace(script1,"eval","result=",-1)
		vm := otto.New()
		_,err = vm.Run(script1)
		if err!=nil {
			panic(err.Error())
		}
		// 获取第二段js
		result, err := vm.Get("result")
		if  err != nil {
			panic(err.Error())
		}
		//处理第二段js
		flag1 :="='__jsl"
		flag2 := "Path=/;'"
		script2,_ := result.ToString()
		script2 = "result" + script2[strings.Index(script2,flag1):strings.Index(script2,flag2) + len(flag2)]


		result2,err := vm.Run(script2)
		//result2, err := vm2.Get("result")
		if  err != nil {
			panic(err)
		}
		cookieValue,_ := result2.ToString()


		cookieValueArr := strings.Split(cookieValue,";")
		metaUrl,_ :=url.Parse("http://www.66ip.cn")
		config.DefaultClient.Jar.SetCookies(metaUrl,[]*http.Cookie{
			{Name: "__jsl_clearance", Value: cookieValueArr[0][len("__jsl_clearance="):]},
		})
		doc = util.GetOfDocument(u)
	}





	bodyEle := doc.Find("body")
	body, _ := bodyEle.Html()
	arr := strings.Split(body, "<br/>\n\t\t")
	var proxyIps []*config.ProxyIp
	//遍历,舍弃最后一个ip，因为还需要切割
	for i := 0; i < len(arr)-1; i++ {
		ip, port, err := net.SplitHostPort(arr[i])
		if err != nil {
			log.Println("thisUrl:", u, " 解析ioPort异常,当前值:", arr[i])
			return 0
		}
		ipUrl, err := util.ParseToUrlOfHttp(ip, port)
		if err != nil {
			log.Println("thisUrl:", u, " 构造ipPort异常，当前值", arr[i])
			return 0
		}
		proxyIps = append(proxyIps,&config.ProxyIp{
			Url:  ipUrl,
			Type: config.Normal,
			From:config.Ip66,
		})
	}
	util.AsyncProxyIpsToChan(config.WaitVerifyChan, proxyIps...)
	return len(proxyIps)
}

