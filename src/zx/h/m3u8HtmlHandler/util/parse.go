package util

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"compress/gzip"
	"net/url"
	"net"
	"io/ioutil"
	"io"
	"zx/ipProxyPool/config"
	"errors"
)

/**
	解析相关
 */

/**
   response 转 string,并关闭response
*/
func ResponseToStr(response *http.Response) (string,error){
	defer response.Body.Close()
	reader,err1 := ResponseToReader(response)
	if err1 != nil {
		err1  = errors.New("url:"+response.Request.URL.Host+" response转string异常:"+err1.Error())
		return "",err1
	}
	bytes, err2 := ioutil.ReadAll(reader)
	if err2 != nil {
		err2 = errors.New("url:"+response.Request.URL.Host+" response转string异常:"+err2.Error())
	}
	return string(bytes),err2


}

/**
	将response解析为reader
 */
func ResponseToReader(response *http.Response) (io.Reader,error){
	var reader io.Reader
	var err error
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			err = errors.New("url:"+response.Request.URL.Host+" gzip解码异常:"+err.Error())
		}
	} else {
		reader = response.Body
	}
	return reader,err
}


 /**
 	将response解析为document,并关闭response
  */
 func ResponseToDocument(response *http.Response) (*goquery.Document,error){
	 defer  response.Body.Close()
	 reader,err1 := ResponseToReader(response)
	 if err1 != nil {
		 err1 = errors.New("url:" + response.Request.URL.Host + " 转为document异常:" + err1.Error())
		 return nil,err1
	 }
	 doc, err2 := goquery.NewDocumentFromReader(reader)
	 if err2 != nil {
		 err2 = errors.New("url:" + response.Request.URL.Host + " 转为document异常:" + err2.Error())
	 }
	 return doc,err2
 }


 /**
 	根据selector，获取指定元素的text
  */
func GetTextBySelector(doc *goquery.Document,selector string) (string) {
	return doc.Find(selector).Text()
}

/**
	根据selector，获取指定元素
 */
func GetElement(doc *goquery.Document,selector string) (element *goquery.Selection) {
	return doc.Find(selector)
}


/**
	根据protocol/ ip/ port/spearation解析出 url.Url
 */
 func ParseToUrlOfSeparation(protocol string,ip string,port string,spearation string) (*url.URL,error) {
	 return  url.Parse(protocol + spearation + net.JoinHostPort(ip, port))
 }

/**
   根据protocol/ ip/ port解析出 url.Url
*/
func ParseToUrl(protocol string,ip string,port string) (*url.URL,error) {
	return ParseToUrlOfSeparation(protocol,ip,port,"")
}

/**
	根据ip/ port解析出 url.Url
 */
func ParseToUrlOfHttp(ip string,port string) (*url.URL,error) {
	return ParseToUrlOfSeparation(config.Http,ip,port,"://")
}


