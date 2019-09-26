package test

import (
	"testing"
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
	"net/http/cookiejar"
	"net/url"
	"time"
	"log"
)

/*
	用go发起htto请求
*/

func TestHttpClient(t *testing.T) {
	u,_ := url.Parse("https://www.google.com")

	ipProxy,err := url.Parse("http://127.0.0.1:8080")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ipProxy)

	/**
		设置http代理
	 */
	urlProxy,_ := url.Parse("http://127.0.0.1:8080")
	//proxyClient := &http.Client{
	//	Transport: &http.Transport{
	//		Pattern:http.ProxyURL(urlProxy),
	//	},
	//}
	fmt.Println(urlProxy)

	/**
		设置socks5代理
	 */
	 //dialer,_ := proxy.SOCKS5("tcp","207.148.25.48:8081",nil,proxy.Direct)
	 //socks5Transport := &http.Transport{Dial:dialer.Dial}
	 //client := &http.Client{
		//Transport:socks5Transport ,
	 //}


	// 构造请求客户端(默客户端所有参数都为nil)
	// 客户端线程安全，官方建议最好重用一个
	client := &http.Client{}
	// 设置超时时间
	client.Timeout =  10 *  time.Second


	// 自定义重定向方法,初始请求是req,每次重定向的请求保存到via分片中,只要返回的error==nil,就会自动重定向
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return nil
	}


	// 设置cookie,其中New()本应传入的参数是指cookie的作用域范围
	// 该cookie会自动保存网站写入的cookie
	jar,_ := cookiejar.New(nil)
	client.Jar = jar
	// 设置cookie，需要设置cookie的作用域范围
	client.Jar.SetCookies(u,[]*http.Cookie{{Name: "a", Value: "b"},{Name: "a", Value: "b"}})


	// 构造请求
	request,err := http.NewRequest(http.MethodGet,"https://www.imooc.com/",strings.NewReader("name=aaa"))
	if err != nil {
		panic(err)
	}



	// 设置请求头
	request.Header.Set("x","x")

	// 进行请求
	response,err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 如下语句可以直接获取整个响应的原始报文,true表示将body部分也输出.
	//httputil.DumpResponse(response,true)

	// 或者直接(实际上内部用一个默认的client，进行了NewRequest，Do等同上的操作)
	//response,err := http.Get("https://www.jianshu.com/p/757d133021de")

	// 进行post form提交，否则构造post请求时需要追加content-type
	//client.PostForm()


	// 输出返回到控制台
	//if _,err=io.Copy(os.Stdout,response.Body);err != nil {
	//	panic(err)
	//}

	// 直接读取为字节
	body,_ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	// 用goquery将其解析为doc
	//doc,err := goquery.NewDocumentFromReader(response.Body)
	//if err != nil {
	//	panic(err)
	//}
	// 用selector选择dom
	//ele := doc.Find("body > div.note > div.post > div.article > div.show-content > div > h1:nth-child(16)")
	//fmt.Println(ele.Text())

	// 打印cookies
	fmt.Println(client.Jar.Cookies(u))


}
