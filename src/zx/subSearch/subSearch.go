package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"regexp"
	"errors"
	"fmt"
	"strconv"
)

/*
搜索思路：
1.执行时获取命令行参数 视频路径和字幕路径
2.读取所有视频文件
3.读取所有字幕文件
4.列出有字幕的文件
	1.提取视频和字幕的番号，分为 番号前缀（例如ipx）和番号后缀（例如909,番号后缀统一去除所有开头的0，例如005.只记录5）
	2.不区分大小写比较番号对象


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
	isExist bool
}
/*no对象比较*/
func (s *no) equals(other *no)bool{
	if strings.EqualFold(s.pre,other.pre) && strings.EqualFold(s.suf,other.suf) {
		return true
	}
	return false
}
/*no对象判断是否为空*/
func (s *no) isNull() bool{
	if s == nil || s.pre == ""  || s.suf == "" {
		return true
	}
	return false
}
const (
	FC2 = "FC2"
	fc2="fc2"
	FLAG = " 中字 "
	ZERO = "0"
)


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
		// 不是视频文件，从分片中删除
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
		subFileN[i] =getNO(temp)
	}

	fmt.Println("以下为有字幕视频:")
	// 比较
	for _, av := range avFileN {
		// no为空时退出该次比较
		if av.isNull(){
			continue
		}
		for _,sub := range subFileN {
			if av.equals(&sub){
				if av.isExist{
					fmt.Println(av.pre +"-" +av.suf + "\t")
				}else{
					fmt.Println(av.pre +"-" +av.suf + "\t" + "未包含中字")
				}
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

			if i.Name() == "System Volume Information"{
				continue
			}

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
var isAVReg = regexp.MustCompile("\\.[\\w]+$")
func isAV(name string)bool{

	suf := isAVReg.FindString(name)
	for _,temp := range avSuf{
		if strings.EqualFold(temp,suf){
			return true
		}
	}
	return false
}

/**判断一个文件是否是字幕文件*/
var isSubReg = regexp.MustCompile("\\.[\\w]+$")
func isSub(name string)bool{
	suf := isSubReg.FindString(name)
	for _,temp := range subSuf{
		if strings.EqualFold(temp,suf){
			return true
		}
	}
	return false
}



/**提取番号
所有番号 [\\w]+.?[\\d]+
^[A-Za-z]+|[\\d]+
*/
var getNOReg = regexp.MustCompile("^([A-Za-z]+|[\\d]+)[-_]?([\\d]+)")
var getNORegFC2 = regexp.MustCompile("[\\d]{4,}")

func getNO(name string)(n no){
	// 处理fc2番号
	if strings.HasPrefix(name,FC2) || strings.HasPrefix(name,fc2){
		temp := getNORegFC2.FindAllString(name,1)
		// 格式错误，直接返回空对象
		if len(temp) < 1{
			fmt.Println("格式错误:" +name)
			return
		}
		return no{FC2, strings.TrimLeft(temp[0],ZERO),strings.Contains(name,FLAG)}
	}

	// 处理其他番号
	temp := getNOReg.FindStringSubmatch(name)
	// 格式错误，直接返回空对象
	if len(temp) < 3 {
		fmt.Println("格式错误:" +name)
		return
	}
	return no{temp[1], strings.TrimLeft(temp[2],ZERO),strings.Contains(name,FLAG)}
}