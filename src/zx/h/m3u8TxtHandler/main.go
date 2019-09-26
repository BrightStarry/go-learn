package main

import (
	"zx/h/m3u8TxtHandler/myLog"
	"zx/h/m3u8TxtHandler/util"
	"io/ioutil"
	"strings"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"os"
	"sync"
)

/**
读取txt，从 链接中获取m3u8和key
 */

func main() {
	defer func() {
		if err := recover(); err != nil {
			myLog.Error("error:%v" ,err)
		}
		util.GetParam("exit")
	}()

	util.ReadConfig("m3u8TxtHandler")
	 // m3u8和key文件保存目录
	downloadDir := viper.GetString("downloadDir")
	// txt文件路径
	txtDir:=  viper.GetString("txtDir")
	// 代理url
	proxyURL :=  viper.GetString("proxyURL")
	// key dir
	keyDir := viper.GetString("keyDir")

	buf,err :=ioutil.ReadFile(txtDir)
	if err != nil {
		panic(err)
	}
	s := string(buf)
	sArr := strings.Split(s,"\r\n")

	var linkArr []string
	for _, i := range sArr {
		if i == "" {
			continue
		}
		link := strings.Split(i,",")[1]
		myLog.Info(link)
		linkArr = append(linkArr, link)
	}

	downloadFunc := func(name,link string) error {
		resp, _, errs :=gorequest.New().Proxy(proxyURL).Get(link).End()
		if len(errs)!= 0 {
			return errs[0]
		}
		raw := resp.Body
		resultBytes,err := ioutil.ReadAll(raw)
		if err != nil {
			return err
		}
		raw.Close()
		if err = ioutil.WriteFile( name,resultBytes,0666); err != nil {
			return err
		}
		return nil
	}

	wg := sync.WaitGroup{}
	for _,item := range linkArr {
		link:= item
		go func() {
			wg.Add(1)
			defer func() {
				wg.Done()
				if err := recover(); err != nil {
					myLog.Error("内部错误:%v",err)
				}
			}()
			number := strings.Split(link,"/")[4]
			myLog.Info("开始下载%s",number)

			// 下载地址
			m3u8Path := downloadDir + string(os.PathSeparator) +number+".m3u8"
			keyPath := downloadDir + string(os.PathSeparator) +number+".key"
			err := downloadFunc(m3u8Path,link)
			if err != nil {
				myLog.Error("当前:%s,下载m3u8异常:%v",number,err)
				return
			}
			err = downloadFunc(keyPath,strings.Replace(link,"0.m3u8","key.bin",-1))
			if err != nil {
				myLog.Error("当前:%s,下载key异常:%v",number,err)
				return
			}
			// 读取m3u8文件
			m3u8Bytes,err := ioutil.ReadFile(m3u8Path)
			if err!= nil {
				myLog.Error("读取m3u8文件异常:%v",err)
				return
			}
			m3u8Str := strings.Replace(string(m3u8Bytes),"key.bin","file@" +keyDir +string(os.PathSeparator)+number+".key",1)

			err = ioutil.WriteFile(m3u8Path,[]byte(m3u8Str),0666)
			if err!= nil {
				myLog.Error("写入m3u8文件异常:%v",err)
				return
			}
			myLog.Info("%s处理完成.",number)
		}()
	}

	wg.Wait()
	myLog.Info("done")
}