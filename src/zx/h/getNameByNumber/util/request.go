package util

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"time"
	"net/http/cookiejar"
)

/**
	请求工具类
 */
var DefaultClient *http.Client

/**
	初始化client
 */
func init() {
	DefaultClient = &http.Client{}
	// 超时时间
	DefaultClient.Timeout = 15 * time.Second
	// 构建cookie
	cookie, err := cookiejar.New(nil)
	if err != nil {
		panic(errors.New("构建cookie失败:"))
	}
	DefaultClient.Jar = cookie
}
var DefaultHeader = map[string][]string{
"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36"},
"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"},
"Connection":                {"keep-alive"},
"Accept-Encoding":           {"gzip, deflate"},
"Accept-Language":           {"zh-CN,zh;q=0.9"},
"Upgrade-Insecure-Requests": {"1"},
"Cache-Control":             {"max-age=0"},
}
 /**
 	根据url，构建GET请求，并追加默认请求头
  */
func BuildGet(url string)(request *http.Request) {
	request,err := http.NewRequest(http.MethodGet,url,nil)
	if err != nil {
		panic(errors.New("url:" + url + " 构建请求异常:"+err.Error()))
	}
	request.Header = DefaultHeader
	return
}



/**
	发送默认get请求
 */
func Get(url string) (*http.Response,error){
	request := BuildGet(url)
	response, err := DefaultClient.Do(request)
	if err != nil {
		err = errors.New("url:" + url + " 请求异常:"+err.Error())
		return nil,err
	}
	if response.StatusCode != http.StatusOK{
		err = errors.New("url:" + url + " 请求异常:"+response.Status)
		return response,err
	}
	return response,nil
}



/**
	发起默认get请求，并解析为document
 */
func GetOfDocument(url string) (*goquery.Document,error) {
	response,err := Get(url)
	if err!= nil {
		return nil,err
	}
	return ResponseToDocument(response)
}

