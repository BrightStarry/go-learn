package main

import (
	"zx/h/videoSub/util"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
	"path/filepath"
	"strconv"
	"sync"
)

/**
视频截图
F:\芽森しずく\SHKD-682 原档 中字 牝畜に堕ちた未亡人 芽森しずく.mp4
E:\新建文件夹
F:\芽森しずく\SHKD-682 原档 中字 牝畜に堕ちた未亡人 芽森しずく.srt
 */

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("error:", err)
		}
		util.GetParam("exit")
	}()

	videoPath := util.GetParam("视频路径:")
	outDir := util.GetParam("截图输出文件夹:")
	txtPath := util.GetParam("字幕路径:")
	outDir += string(os.PathSeparator)

	if !util.FileIsExist(videoPath){
		log.Panicln("视频不存在:",videoPath)
	}
	if !util.FileIsExist(outDir){
		log.Panicln("截图输出文件夹不存在:",videoPath)
	}

	txtBytes, err := ioutil.ReadFile(txtPath)
	if err != nil {
		log.Panicln("读取字幕失败：",err)
	}
	timeTxt := string(txtBytes)

	txtExt := filepath.Ext(txtPath)

	timeArr := make([]string,0)
	if strings.Contains(txtExt, "srt") {
		// 提取时间
		var getTimeReg = regexp.MustCompile(`\d{1,2}:\d{1,2}:\d{1,2},\d{0,3}`)
		temp := getTimeReg.FindAllStringSubmatch(timeTxt,-1)
		for _, v := range temp {
			timeArr = append(timeArr,strings.Replace( v[0],",",".",1))
		}
	}else if  strings.Contains(txtExt, "ass") ||  strings.Contains(txtExt, "ssa") {
		// 提取时间
		var getTimeReg = regexp.MustCompile(`\d:\d{2}:\d{2}\.\d{2}`)
		temp := getTimeReg.FindAllStringSubmatch(timeTxt,-1)
		for _, v := range temp {
			timeArr = append(timeArr,v[0])
		}
	}else if  strings.Contains(txtExt, "vtt") {
		var getTimeReg = regexp.MustCompile(`\d{2}:\d{2}:\d{2}\.\d{3}`)
		temp := getTimeReg.FindAllStringSubmatch(timeTxt,-1)
		for _, v := range temp {
			timeArr = append(timeArr,v[0])
		}
	}else {
		log.Panicln("暂不支持该格式字幕.")
	}


	if len(timeArr)< 1 {
		log.Panicln("无时间节点")
	}

	log.Println("读取到时间节点",len(timeArr),"个,开始截图...")

	group := sync.WaitGroup{}
	for i := 0; i < len(timeArr); i++ {
		if i%2 == 0 {
			tempI := i
			group.Add(1)
			go func() {
				defer group.Done()
				util.StartCMD("cmd","/c","ffmpeg","-ss",timeArr[tempI],"-i",
					videoPath  ,
					"-vframes","1","-f","image2","-y",
					outDir + strconv.Itoa((tempI+2)/2)+ "-" + strings.Replace(timeArr[tempI],":",",",-1) + ".jpg")
			}()

			if i %10 == 0 {
				group.Wait()
			}
		}
	}

	group.Wait()
	log.Println("ok.")


}
