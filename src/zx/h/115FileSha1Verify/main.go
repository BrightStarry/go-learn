package main

import (
	"zx/h/115FileSha1Verify/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"errors"
	"github.com/parnurzeal/gorequest"
	"time"
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"
	"path/filepath"
	"io"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

/**
校验115下载的文件sha1是否和其原文件一致
 */

var config Config
var baseRequest *gorequest.SuperAgent
var searchUrlFormat = `https://webapi.115.com/files/search?offset=0&limit=115&search_value=%s&date=&aid=1&cid=0&pick_code=&type=&source=&format=json`
func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("error:", err)
		}
		util.GetParam("exit")
	}()

	destDir := util.GetParam("需要校验的目录（直接回车则使用配置文件中的目录）:")
	if strings.TrimSpace(destDir)!= "" {
		config.VerifyDir = strings.TrimSpace(destDir)
	}

	filePaths := getAllFileName(config.VerifyDir)

	var errResult  []string

	for i, file := range filePaths {
		log.Println("开始处理第",strconv.Itoa(i+1),"个文件:",file)
		sha1Str,err := getSh1(file)
		if err != nil {
			log.Warn("获取sha1失败:",err,",文件:",file)
			errResult = append(errResult, file)
			continue
		}
		log.Println("当前文件sha1:",sha1Str)

		filenameWithExt := filepath.Base(file)
		datas,err := doSearch(filenameWithExt)
		if err != nil {
			log.Warn("115搜索异常:",err,",文件:",file)
			errResult = append(errResult, file)
			continue
		}
		verifyFlag := false
		for j, data := range datas {
			log.Println("比对搜索到的第",strconv.Itoa(j+1),"个文件:",data.name)
			if strings.EqualFold(sha1Str, data.sha1)  {
				verifyFlag = true
				log.Println("校验成功.文件:",file)
				break
			}
			if j == config.SearchLimit-1 {
				break
			}
		}
		if !verifyFlag {
			log.Warn("校验失败，没有相匹配的sha1,文件:",file)
			errResult = append(errResult, file)
		}
	}

	log.Println("校验失败文件如下:")
	for _, item := range errResult {
		log.Println(item)
	}
}

/**
获取文件sha1
 */
 func getSh1(filePath string) (sha1Str string,err error) {
	 file, err := os.Open(filePath)
	 if err != nil {
		 return "",errors.New("打开文件失败" + err.Error())
	 }
	 defer file.Close()
	 h := sha1.New()
	 _, err = io.Copy(h, file)
	 if err != nil {
		 return
	 }
	 return hex.EncodeToString(h.Sum(nil)),nil
 }

/**
请求搜索url并获取响应
 */
 func doSearch(searchKey string)( datas []Data,err error){
	 resp, body, errs := baseRequest.Clone().Get(fmt.Sprintf(searchUrlFormat,searchKey)).End()
	 if len(errs) != 0 {
		 return nil,errs[0]
	 }
	 if resp.StatusCode != http.StatusOK {
		 return nil,errors.New("http异常:" + strconv.Itoa(resp.StatusCode))
	 }
	 // json转map
	 var resultMap map[string]interface{}
	 if err = json.Unmarshal([]byte(body), &resultMap); err!= nil {
		 return
	 }
	 errorStr := resultMap["error"].(string)
	 if errorStr != "" {
		 return nil,errors.New("115响应异常:" + errorStr)
	 }

	 switch data := resultMap["data"].(type) {
	 case []interface{}:
		 for _, item := range data {
			 itemMap := item.(map[string]interface{})
			 // 不是文件跳过
			 if _,ok := itemMap["fid"];!ok{
				 continue
			 }
			 name := itemMap["n"].(string)
			 sha1Str := itemMap["sha"].(string)
			 datas = append(datas, Data{name,sha1Str})
			 if len(datas) >= config.SearchLimit-1 {
			 	break
			 }
		 }
	 case map[string]interface{}:
	 	 dateItemInterface,ok := data["0"]
		 for i := 1; ok && i <= config.SearchLimit;i++  {
			 dataTemp := dateItemInterface.(map[string]interface{})
			 // 不是文件跳过
			 if _,ok2 := dataTemp["fid"];!ok2{
				 continue
			 }
			 name := dataTemp["n"].(string)
			 sha1Str := dataTemp["sha"].(string)
			 datas = append(datas, Data{name,sha1Str})

			 dateItemInterface,ok = data[strconv.Itoa(i)]
		 }
	 default:
		 return nil, errors.New("data数据类型无法处理")
	 }
	 if len(datas) <= 0 {
		 return nil,errors.New("搜索结果为0")
	 }
	 return
 }

/*遍历出文件夹中所有文件名,不包含目录*/
func getAllFileName(dirPath string)(fileNames []string){
	// 读取目录
	fileInfo,err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(errors.New("目录读取异常:" + err.Error()))
	}
	for _,i := range fileInfo{
		if i.IsDir(){
			if i.Name() == "System Volume Information"{
				continue
			}
			//temp := getAllFileName(dirPath + string(os.PathSeparator) + i.Name())
			//fileNames = append(fileNames, temp...)
		}else{
			fileNames = append(fileNames, dirPath + string(os.PathSeparator) + i.Name())
		}
	}
	return
}
type Data struct {
	name string
	sha1 string
}
type Config struct {
	Cookie string
	SearchLimit int
	VerifyDir string
}
func init() {
	util.InitConfig("115FileSha1Verify",&config)

	baseRequest = gorequest.New().
		AppendHeader("Cookie",config.Cookie).
		Timeout(10*time.Second).
		Retry(2, 0 * time.Second, http.StatusBadRequest, http.StatusInternalServerError,http.StatusRequestTimeout)
}