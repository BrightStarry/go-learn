package main

import (
	log "github.com/sirupsen/logrus"
	"zx/h/m3u8HtmlHandler/util"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"strconv"
	"github.com/spf13/viper"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"io/ioutil"
	"path/filepath"
	"os"
)
var config Config
/**
从html网页中获取m3u8
 */
func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("error:" ,err)
		}
		util.GetParam("exit")
	}()

	resp, _, errs := gorequest.New().Clone().Proxy(config.ProxyUrl).Get(config.Url).End()
	if len(errs) != 0 {
		log.Panicln(errs)
	}
	if resp.StatusCode != http.StatusOK {
		log.Panicln("http异常:" + strconv.Itoa(resp.StatusCode))
	}

	doc, err := util.ResponseToDocument(resp)
	if err != nil {
		log.Panicln("获取document异常：",err)
	}
	h1Elements := util.GetElement(doc, "h1")
	aElements := util.GetElement(doc, "a")
	imgElements := util.GetElement(doc, "img")

	h1Elements.Each(func(i1 int, h1Ele1 *goquery.Selection) {
		i := i1
		h1Ele := h1Ele1
		go func() {
			m3u8Url, m3u8UrlExist := aElements.Eq(i * 2).Attr("href")
			if !m3u8UrlExist {
				log.Panicln("格式有误:",h1Elements.Text())
			}
			imgUrl, imgUrlExist := imgElements.Eq(i).Attr("src")
			if !imgUrlExist {
				log.Panicln("格式有误:",h1Elements.Text())
			}
			number := h1Ele.Text()

			keyUrl := strings.Replace(m3u8Url,"0.m3u8","key.bin",-1)
			log.Println("index:",i,",番号:",number,",m3u8下载路径:", m3u8Url,",img下载路径:", imgUrl,",key下载路径:", keyUrl)

			imgExt := filepath.Ext(imgUrl)
			pathPre := config.M3u8Path + string(os.PathSeparator) + number
			// 下载
			if err = download(pathPre + ".m3u8", m3u8Url);err != nil {
				log.Errorln("番号:",number,",下载异常:",err)
			}
			if err = download(pathPre + ".key", keyUrl);err != nil {
				log.Errorln("番号:",number,",下载异常:",err)
			}
			if err = download(pathPre + imgExt, imgUrl);err != nil {
				log.Errorln("番号:",number,",下载异常:",err)
			}
		}()
	})



}

func download(path,url string)error {
	resp, _, errs :=gorequest.New().Proxy(config.ProxyUrl).Get(url).End()
	if len(errs)!= 0 {
		return errs[0]
	}
	raw := resp.Body
	resultBytes,err := ioutil.ReadAll(raw)
	if err != nil {
		return err
	}
	raw.Close()
	if err = ioutil.WriteFile( path,resultBytes,0666); err != nil {
		return err
	}
	return nil
}

type Config struct {
	// 要解析的url
	Url string
	// m3u8文件保存位置
	M3u8Path string
	// 代理url
	ProxyUrl string
}

func init() {
	configReader := viper.New()
	configReader.SetConfigName("m3u8HtmlHandler")
	configReader.AddConfigPath("./")
	configReader.SetConfigType("yaml")
	if err := configReader.ReadInConfig(); err != nil {
		log.Panicln("读取配置异常:" + err.Error())
	}
	if err := configReader.Unmarshal(&config); err != nil {
		log.Panicln("读取配置异常:" + err.Error())
	}

	log.WithFields(log.Fields{
		"config":config,
	}).Info("当前参数.")
}
