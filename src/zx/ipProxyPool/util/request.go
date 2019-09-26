package util

import (
	"net/http"
	"zx/ipProxyPool/config"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"net/url"
)

/**
	请求工具类
 */


 /**
 	根据url，构建GET请求，并追加默认请求头
  */
func BuildGet(url string)(request *http.Request) {
	request,err := http.NewRequest(http.MethodGet,url,nil)
	if err != nil {
		panic(errors.New("url:" + url + " 构建请求异常:"+err.Error()))
	}
	request.Header = *config.DefaultHeader
	return
}



/**
	发送默认get请求
 */
func Get(url string) (response *http.Response){
	request := BuildGet(url)
	response, err := config.DefaultClient.Do(request)
	if err != nil {
		panic(errors.New("url:" + url + " 请求异常:"+err.Error()))
	}
	if response.StatusCode != http.StatusOK{
		panic(errors.New("url:" + url + " 请求异常:"+response.Status))
	}
	return
}

/**
	用代理client发起get请求
 */
func GetByProxy(url string,proxy *url.URL)(response *http.Response) {
	request := BuildGet(url)
	// 从池中获取
	client := config.VerifierClientPool.Get().(*http.Client)
	defer config.VerifierClientPool.Put(client)
	// 设置代理
	client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxy)

	response, err := client.Do(request)
	if err != nil {
		panic(errors.New("url:" + url + " 请求异常:"+err.Error()))
	}
	return
}

/**
	发起默认get请求，并解析为document
 */
func GetOfDocument(url string) (*goquery.Document) {
	return ResponseToDocument(Get(url))
}

