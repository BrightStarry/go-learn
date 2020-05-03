package main

import (
	"encoding/base64"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"zx/h/exportKey/util"
	"os"
	"errors"
	"strings"
	"sync"
	"regexp"
	"github.com/spf13/viper"
	"path/filepath"
)
var reg = regexp.MustCompile("URI=\"(.+)\"")
func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("error:", err)
		}
		util.GetParam("exit")
	}()

	// 获取所有m3u8文件
	fileNames := getAllFileName(config.M3u8Dir)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(fileNames))
	for _, item := range fileNames {
		fileName := item
		go func() {
			defer func() {
				waitGroup.Done()
				if err := recover(); err != nil {
					log.Errorln("任务:",fileName,",内部异常:",err)
				}
			}()

			log.Println("任务:",fileName,",开始")
			fileBytes, err := ioutil.ReadFile(fileName)
			if err != nil {
				log.Errorln("读取文件异常:",err)
				return
			}
			m3u8Str := string(fileBytes)
			temp:=reg.FindStringSubmatch(m3u8Str)
			if len(temp) != 2 {
				log.Errorln("任务:",fileName,",提取异常,key格式不正确:",temp)
				return
			}
			keyStr := temp[1]
			if !strings.HasPrefix(keyStr, "base64:") {
				log.Errorln("任务:",fileName,",提取异常,key格式不正确:",temp)
				return
			}
			base64Key := strings.TrimPrefix(keyStr,"base64:")

			base64Encoder :=base64.StdEncoding
			bytes, err := base64Encoder.DecodeString(base64Key)
			if err != nil {
				log.Errorln("任务:",fileName,",base64解码异常:",err)
				return
			}

			nameExt := filepath.Base(fileName)
			ext := filepath.Ext(fileName)
			number := strings.TrimSuffix(nameExt, ext)
			err = ioutil.WriteFile(config.KeyDir + string(os.PathSeparator) + number + ".key", bytes, 0666)
			if err != nil {
				log.Errorln("任务:",fileName,",写入key文件异常:",err)
				return
			}
		}()
	}

	waitGroup.Wait()
	log.Println("任务完成.")





}
var config Config
type Config struct {
	// m3u8文件目录
	M3u8Dir string
	// 提取出来的key存放目录
	KeyDir string
}

func init() {
	configReader := viper.New()
	configReader.SetConfigName("exportKey")
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
			// 提取m3u8
			if strings.Contains(i.Name(), "m3u8") {
				fileNames = append(fileNames, dirPath + string(os.PathSeparator) + i.Name())
			}
		}
	}
	return
}
