package main

import (
	"flag"
	"zx/hijack/server/util"
	"log"
)

func main() {
	util.SyncStartWebServer()
}

/**
	 初始化
 */
func init() {
	readExternalParam() // 读取外部参数
	log.Println("准备启动，当前参数:",util.Config)
}

/**
	读取外部参数
 */
func readExternalParam() {
	// 端口
	var port string
	flag.StringVar(&port, "port", util.Config.Port, "端口")

	// pac文件路径
	var pacPath string
	flag.StringVar(&pacPath, "pac", util.Config.PacPath, "pac文件路径")

	// 数据存储路径
	var dataPath string
	flag.StringVar(&dataPath, "data", util.Config.DataPath, "数据文件路径")

	flag.Parse()

	util.Config.Port = port
	util.Config.PacPath = pacPath
}