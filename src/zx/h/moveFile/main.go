package main

import (
	"zx/h/moveFile/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
)
var config Config
func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("error:", err)
		}
		util.GetParam("exit")
	}()

	wmvPath := util.GetParam("wmv路径:")
	wmvPath = strings.Trim(wmvPath,`"`)

	wmvFileWithSuf := filepath.Base(wmvPath)
	wmvFileExt := filepath.Ext(wmvPath)
	wmvFilenameOnly := strings.TrimSuffix(wmvFileWithSuf,wmvFileExt)

	savePre := config.WmvSaveDir +string(os.PathSeparator)+ wmvFilenameOnly+string(os.PathSeparator)
	if err :=os.MkdirAll(savePre,0666);err!= nil {
		log.Panicln("创建wmv文件失败:",err)
	}
	if err :=os.Rename(wmvPath,savePre+ wmvFileWithSuf);err!= nil {
		log.Panicln("移动wmv文件失败:",err)
	}

	newHdsBytes,err := ioutil.ReadFile(config.NewHdsFilePath)
	if err!= nil {
		log.Panicln("复制生成的hds文件,读取失败:",err)
	}
	if err = ioutil.WriteFile(savePre + filepath.Base(config.NewHdsFilePath),newHdsBytes,0666);err!= nil {
		log.Panicln("复制生成的hds文件，写入失败:",err)
	}

	oldHdsBytes,err := ioutil.ReadFile(config.OldHdsFilePath)
	if err!= nil {
		log.Panicln("复制初始hds文件,读取失败:",err)
	}
	if err = ioutil.WriteFile(config.NewHdsFilePath,oldHdsBytes,0666);err!= nil {
		log.Panicln("复制初始hds文件，写入失败:",err)
	}
}

type Config struct {
	WmvSaveDir string
	NewHdsFilePath string
	OldHdsFilePath string
}

func init() {
	initConfig()
}
func initConfig() {
	configReader := viper.New()
	configReader.SetConfigName("moveFile")
	configReader.AddConfigPath("./")
	configReader.SetConfigType("yaml")
	if err := configReader.ReadInConfig(); err != nil {
		log.Panicln("读取配置异常:" + err.Error())
	}
	if err := configReader.Unmarshal(&config); err != nil {
		log.Panicln("读取配置异常:" + err.Error())
	}
}

