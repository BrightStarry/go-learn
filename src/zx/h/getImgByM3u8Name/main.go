package main

import (
	log "github.com/sirupsen/logrus"
	"zx/h/getImgByM3u8Name/util"
	"io/ioutil"
	"os"
	"errors"
	"path/filepath"
	"strings"
	"github.com/parnurzeal/gorequest"
	"fmt"
	"time"
	"net/http"
	"sync"
)

/**
根据m3u8名字获取封面图
 */
var commonImgUrl = "http://pics.dmm.co.jp/mono/movie/adult/%s/%spl.jpg"
var soplImgUrl = "http://pics.dmm.co.jp/mono/movie/adult/%sso/%ssopl.jpg"
func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("error:", err)
		}
		util.GetParam("exit")
	}()

	m3u8Dir := util.GetParam("请输入目录:")
	//遍历
	files := getAllFileName(m3u8Dir)
	group := sync.WaitGroup{}
	for i, item1 := range files {
		item := item1
		ext := filepath.Ext(item)
		if ext != ".m3u8" {
			continue
		}
		go func() {
			defer group.Done()
			group.Add(1)
			currentDir := filepath.Dir(item) + string(os.PathSeparator)
			fileName := filepath.Base(item)
			onlyFileName := strings.TrimSuffix(fileName,ext)
			if strings.Contains(onlyFileName, "_part") {
				onlyFileName = strings.Split(onlyFileName,"_")[0]
			}

			// 验证图片是否存在
			if util.FileIsExist(currentDir + onlyFileName + ".jpg") {
				bytes,_ := ioutil.ReadFile(currentDir + onlyFileName + ".jpg")
				if len(bytes) > 3000 {
					log.Warn("图片已存在:",item)
					return
				}
			}

			log.Println("开始下载:",item)
			url := fmt.Sprintf(commonImgUrl,strings.Replace(onlyFileName,"00","",1),strings.Replace(onlyFileName,"00","",1))

			resultBytes,err := getImgBytes(url)
			if err != nil {
				log.Warn("下载图片:",item,"异常:",err)
			}
			// 小于3000 说明爬取有错误
			if len(resultBytes) < 3000 {
				url = fmt.Sprintf(soplImgUrl,strings.Replace(onlyFileName,"00","",1),strings.Replace(onlyFileName,"00","",1))
				resultBytes,err = getImgBytes(url)
				if err != nil {
					log.Warn("下载图片:",item,"异常:",err)
					return
				}
			}
			if err = ioutil.WriteFile(currentDir + onlyFileName + ".jpg",resultBytes,0666); err != nil {
				log.Warn("保存图片:",item,"异常:",err)
				return
			}
		}()

		if i % 30 == 0 {
			group.Wait()
		}
	}
	group.Wait()
	log.Println("ok")
}

/**
获取图片bytes
 */
 func getImgBytes(url string)([]byte,error) {
	 resp, _, errs := gorequest.New().Timeout(20 *time.Second).
		 Retry(3, 0 * time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		 Get(url).End()
	 if len(errs) > 0 {
		 return nil,errs[0]
	 }
	 raw := resp.Body
	 raw.Close()
	 resultBytes,err := ioutil.ReadAll(raw)
	 if err != nil {
		 return nil,err
	 }
	 return resultBytes,nil
 }

/*遍历出文件夹中所有文件名*/
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
			temp := getAllFileName(dirPath + string(os.PathSeparator) + i.Name())
			fileNames = append(fileNames, temp...)
		}else{
			fileNames = append(fileNames, dirPath + string(os.PathSeparator) + i.Name())
		}
	}
	return
}