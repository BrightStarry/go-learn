package main

import (
	"github.com/parnurzeal/gorequest"
	"time"
	"net/http"
	"strconv"
	"zx/h/115/util"
	"encoding/json"
	"zx/h/115/myLog"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"path/filepath"
	"regexp"
	"github.com/spf13/viper"
	"sync"
	"fmt"
	"encoding/hex"
	"bytes"
	"encoding/binary"
)


/**
 指定115目录， 获取番号， 转为普通番号，再自动获取片名，重命名
 */
var keyword = []string{":","*","?","/","\\","\"","<",">","|"}
// 默认网址s
var defaultUrls = []string{"http://www.f37b.com","https://www.dmmbus.co"}
// 默认网址后缀
// javlibrary http://www.m34z.com/cn/vl_searchbyid.php?keyword=番号
// dmmbus https://www.dmmbus.co/番号
var defaultUrlSufs = []string{"/cn/vl_searchbyid.php?keyword=","/"}
// 默认获取片名的方法
var defaultGetNameFuns = []func(*goquery.Document)string {
	func(doc *goquery.Document)string {
		return util.GetTextBySelector(doc,"#video_title > h3 > a")
	},
	func(doc *goquery.Document)string {
		return util.GetTextBySelector(doc,"body > div.container > h3")
	},
}
// 解析dmm番号
var numberReg = regexp.MustCompile("(h_)?(\\d*)((t28)|([A-Za-z]+))0*([\\d]+)([A-Za-z])?(_part)?(\\d)?")



func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error:", err)
		}
		util.GetParam("exit")
	}()

	// 读取配置文件
	viper.SetConfigName("115")
	viper.AddConfigPath("./")
	//viper.AddConfigPath("./src/zx/h/115")
	//viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic("读取配置异常:" +err.Error())
	}
	cid := viper.Get("cid").(string)
	cookie := viper.Get("cookie").(string)
	offset := viper.Get("offset").(string)
	isDMM := viper.Get("isDMM").(bool)

	baseRequest := gorequest.New().
		AppendHeader("Cookie",cookie).
		//AppendHeader("User-Agent","Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36 115Browser/9.2.1").
		Retry(3, 0 * time.Second, http.StatusBadRequest, http.StatusInternalServerError)

	resp, body, errs := baseRequest.Clone().
		Get("https://webapi.115.com/files?aid=1&cid="+cid+"&o=user_ptime&asc=0&offset="+offset+"&show_dir=1&limit=115&code=&scid=&snap=0&natsort=1&source=&format=json&type=&star=&is_share=&suffix=&custom_order=&fc_mix=&is_q=").
		End()
	if len(errs) != 0 {
		panic(errs)
	}
	if resp.StatusCode != http.StatusOK {
		panic("http异常:" + strconv.Itoa(resp.StatusCode))
	}
	// json转map
	var filesMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &filesMap); err!= nil {
		panic(err)
	}
	errorStr := filesMap["error"].(string)
	if errorStr != "" {
		fmt.Println("115响应异常:",errorStr)
		return
	}



	waitGroup := sync.WaitGroup{}
	// data下就是每个文件，遍历
	for _, i := range filesMap["data"].([]interface{}) {
		goI := i
		go func() {
			defer func(){
				waitGroup.Done()
				if err := recover(); err != nil {
					myLog.Error("内部错误:%v",err)
				}
			}()
			waitGroup.Add(1)
			item := goI.(map[string]interface{})
			fidInterface,ok := item["fid"]
			// 不是文件跳过
			if !ok{
				return
			}
			fid := fidInterface.(string)
			// 文件名
			name := item["n"].(string)
			name = strings.TrimSuffix(name,filepath.Ext(name))
			oldName := name
			// 如果包含空格，则跳过
			if strings.Contains(name, " ")  {
				myLog.Warn("文件名包含空格，可能已经重命名:%s" ,name)
				return
			}
			if isDMM{
				// dmm番号转普通番号
				tempNumber := numberReg.FindStringSubmatch(name)
				pre := tempNumber[3]
				suf := tempNumber[6]
				index := tempNumber[9]
				switch len(suf) {
				case 1:
					suf = "00" + suf
				case 2:
					suf = "0" + suf
				}
				name = pre + "-" + suf + index

			}


			// 解析番号
			no := util.GetNO(name)
			if no.IsNull() {
				myLog.Warn("番号解析有误:%s" , name)
				return
			}
			// 给番号添加0
			suf := no.Suf
			switch len(suf) {
			case 1:
				suf = "00" + suf
			case 2:
				suf = "0" + suf
			}
			// 正常番号
			number := no.Pre + "-" +suf

			avName := ""
			var newAVName string
			for j:=0; j < len(defaultUrls) && avName == ""; j++ {
				// 当前url
				currentUrl := defaultUrls[j] + defaultUrlSufs[j] + number
				document,err1 := util.GetOfDocument(currentUrl)
				if err1 != nil {
					continue
				}
				// 获取片名
				avName = defaultGetNameFuns[j](document)
				if avName == "" {
					continue
				}
				newAVName = avName + "~" + oldName
			}

			if newAVName == "" {
				myLog.Error("当前文件名:%s,获取番号失败",name)
				return
			}

			// 替换特殊字符
			for i := 0; i< len(keyword);i++  {
				newAVName = strings.Replace(newAVName,keyword[i]," ",-1)
			}

			resp, body, errs :=  baseRequest.Clone().
				Post("https://webapi.115.com/files/edit").
				Type("form").
				Send(`{"fid":"`+fid+`","file_name":"`+newAVName+`"}`).
				End()
			if len(errs) != 0 {
				myLog.Error("当前文件名:%s,重命名响应异常:%v",name,errs)
			}
			if resp.StatusCode != http.StatusOK {
				myLog.Error("当前文件名:%s,重命名响应http异常:%d",name,resp.StatusCode)
			}
			// json转map
			var renameResultMap map[string]interface{}
			if err := json.Unmarshal([]byte(body), &renameResultMap); err!= nil {
				myLog.Error("当前文件名:%s,重命名响应异常:%s",name,body)
			}
			stateInterface,ok := renameResultMap["state"]
			if !ok{
				myLog.Error("当前文件名:%s,重命名响应异常:%s",name,body)
				return
			}
			if state := stateInterface.(bool);!state{
				myLog.Error("当前文件名:%s,重命名响应异常:%s",name,body  )
				return
			}
			myLog.Info("success! %s",newAVName)
		}()
	}

	waitGroup.Wait()
	myLog.Info("done!")
}

func u2s(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}