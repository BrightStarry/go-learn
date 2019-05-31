package main

import (
	"flag"
	"io/ioutil"
	"fmt"
	"os"
	"strings"
	"regexp"
	"strconv"
	"errors"
)

/*
搜索思路：
1.从配置文件读取字幕路径和视频路径，视频路径默认为空，如果为空，则要在运行后输入视频地址。
2.读取视频文件名，去重.
3.读取所有字幕文件名
4.列出有字幕的文件
	1.从英文开始读取，到其他字符结束，为番号头
	2.读取任意字符(空或其他字符)，读取到数字，在从数字开始读取到其他任意字符结束，为番号编号


*/


/*系统配置*/
type config struct{
	// 字幕路径
	subPath string
	// 视频路径
	avPath string
}
var defaultConfig = new(config)

var avSuf  = []string{".mp4",".mkv",".wmv",".avi",".MOV",".ASF",".asx",".MPEG",".mpg",".ISO",".3GP",".FLV",".F4V",".RMVB",".dat",".vob",".m2ts"}
var subSuf = []string{".SSA",".ASS",".SMI",".SRT",".SUB",".LRC",".SST",".vtt"}
type no struct {
	pre string
	suf string
}


func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error:" ,err)
		}
		var temp string
		fmt.Scan(&temp)
	}()
	err := readExternalParam()
	if err != nil {
		panic(errors.New("参数读取失败:" + err.Error()))
	}
	if defaultConfig.subPath == "" {
		panic(errors.New("字幕路径为空" ))
	}
	if defaultConfig.avPath == "" {
		panic(errors.New("视频路径为空" ))
	}
	// 输出参数
	fmt.Println("视频路径:" + defaultConfig.avPath +"; 字幕路径:"+ defaultConfig.subPath)

	// 读取视频文件
	avFiles,err := getAllFileName(defaultConfig.avPath)

	fmt.Println("读取视频目录文件总数:" + strconv.Itoa(len(avFiles)))

	// 提取视频文件
	for i:=0;i<len(avFiles);{
		if !isAV(avFiles[i]) {
			avFiles = append(avFiles[:i],avFiles[i+1:]...)
		}else{
			i++
		}
	}

	fmt.Println("读取视频目录视频文件总数:" + strconv.Itoa(len(avFiles)))

	// 获取视频番号
	avFileN := make([]no,len(avFiles))
	for i,temp := range avFiles {
		avFileN[i] = getNO(temp)
	}

	// 读取字幕
	subFiles,err := getAllFileName(defaultConfig.subPath)
	fmt.Println("读取字幕目录文件总数:" + strconv.Itoa(len(subFiles)))
	// 提取字幕文件
	for i:=0;i<len(subFiles);{
		if !isSub(subFiles[i]) {
			subFiles = append(subFiles[:i],subFiles[i+1:]...)
		}else{
			i++
		}
	}
	fmt.Println("读取字幕目录字幕文件总数:" + strconv.Itoa(len(subFiles)))

	// 获取字幕番号
	subFileN :=  make([]no,len(subFiles))
	for i,temp := range subFiles {
		subFileN[i] = getNO(temp)
	}

	fmt.Println("以下为有字幕视频:")
	// 比较
	for _, av := range avFileN {
		for _,sub := range subFileN {
			if av == sub{
				fmt.Println(av.pre +"-" +av.suf)
				break
			}
		}
	}













}


/*读取外部参数*/
func readExternalParam() (err error) {
	flag.StringVar(&defaultConfig.subPath,"subPath", "", "字幕路径")
	flag.StringVar(&defaultConfig.avPath , "avPath", "", "视频路径")
	flag.Parse()
	return
}

/*遍历出文件夹中所有文件名*/
func getAllFileName(dirPath string)(fileNames []string,err error){
	// 读取目录
	fileInfo,err := ioutil.ReadDir(dirPath)
	for _,i := range fileInfo{
		if i.IsDir(){
			temp,err := getAllFileName(dirPath + string(os.PathSeparator) + i.Name())
			if err != nil {
				panic(errors.New("目录读取异常:" + err.Error()))
			}
			fileNames = append(fileNames, temp...)
		}else{
			fileNames = append(fileNames, i.Name())
		}

	}
	return
}

/**判断一个文件是否是视频文件*/
func isAV(name string)bool{
	reg := regexp.MustCompile("\\.[\\w]+$")
	suf := reg.FindString(name)
	for _,temp := range avSuf{
		if strings.EqualFold(temp,suf){
			return true
		}
	}
	return false
}

/**判断一个文件是否是字幕文件*/
func isSub(name string)bool{
	reg := regexp.MustCompile("\\.[\\w]+$")
	suf := reg.FindString(name)
	for _,temp := range subSuf{
		if strings.EqualFold(temp,suf){
			return true
		}
	}
	return false
}



/**提取番号
所有番号 [\\w]+.?[\\d]+
前缀
后缀
*/
func getNO(name string)(n no){
	reg := regexp.MustCompile("^[A-Za-z]+|[\\d]+")
	temp := reg.FindAllString(name,2)
	return no{temp[0], temp[1]}
}