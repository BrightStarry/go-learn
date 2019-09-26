package main

import (
	"bufio"
	"os"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	path2 "path"
	"strconv"
	"os/exec"
	"regexp"
	"io/ioutil"
)

/**
视频剪辑
 */
func main() {


	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error:" ,err)
		}
		getParam("exit")
	}()
	for{
		flagStr := getParam("1:剪辑；2:合并；")
		flag,err := strconv.Atoi(flagStr)
		if err == nil && (flag ==1||flag==2){
			if flag == 1{
				clip()
			}else{
				merge()
			}
		}
	}

}

const(
	DefaultTime = "00:00:00"
	TempDirectory="C:" + string(os.PathSeparator)+"temp"+ string(os.PathSeparator)
	TempFile = "video_clip_merga_temp.txt"
	MergaOutputFile = "merga"
)
var videoSuf  = []string{".mp4",".mkv",".wmv",".avi",".MOV",".ASF",".asx",".MPEG",".mpg",".ISO",".3GP",".FLV",".F4V",".RMVB",".dat",".vob",".m2ts"}



/**
合并
 */
 func merge() {
	 // 获取视频文件路径
	 inputDirectory := getParam("请输入视频目录：")+ string(os.PathSeparator)
	 if inputDirectory == "" {
		 panic(errors.New("视频目录不能为空!" ))
	 }
	 // 获取目录下所有视频文件
	 fileNames :=getvideoFile(inputDirectory)
	 fmt.Println("该目录视频列表如下:")
	 for i:=0; i<len(fileNames); i++ {
		 fmt.Println(strconv.Itoa(i) + ":\t" + fileNames[i])
	 }
	 // 获取选定的视频
	 var videoFileNames []string
	 videoFileNoStr := getParam("请输入视频编号(用空格分隔)：")
	 videoFileNoStrs := strings.Split(strings.Trim(videoFileNoStr," ")," ")
	 for _, value := range videoFileNoStrs{
		 valueInt,_ :=strconv.Atoi(value)
		 videoFileNames = append(videoFileNames,fileNames[valueInt])
	 }

	 // 拼接临时文件字符
	 var tempContent string
	 for _,value :=range videoFileNames{
	 	tempContent = tempContent + "file '" +  inputDirectory + value + "'\r\n"
	 }
	 // 写入临时文件
	 tempFilePath := getNotExistFilePath(TempDirectory,TempFile)
	 os.Mkdir(TempDirectory, os.ModePerm)
	 createFile(tempFilePath)
	 defer deleteFile(tempFilePath)
	 if err :=ioutil.WriteFile(tempFilePath, []byte(tempContent),0666);err != nil {
	 	panic("写入临时文件失败:"+ err.Error())
	 }

	 // 输出路径
	 outputDirectory := getParam("请输入视频输出目录(为空则为当前视频目录,默认输出文件名为merga.视频后缀)：")
	 if outputDirectory == "" {
		 outputDirectory = inputDirectory
	 }else {
		 outputDirectory = outputDirectory + string(os.PathSeparator)
	 }

	 // 获取视频后缀
	 fileSuffix :=  path2.Ext(videoFileNames[0])
	 outFilePath := getNotExistFilePath(outputDirectory,MergaOutputFile + fileSuffix)
	 fmt.Println("请等待...")
	 cmd := exec.Command("ffmpeg.exe",
		 "-f","concat",
		 "-safe","0",
		 "-i",tempFilePath,
		 "-c", "copy",
		 "-copyts",outFilePath)
	 if err := cmd.Run(); err != nil {
		 panic("合并命令异常:" + err.Error())
	 }
	 fmt.Println("done")


 }

/**
 剪辑
 */
func clip() {
	// 获取视频文件路径
	inputDirectory := getParam("请输入视频目录：")+ string(os.PathSeparator)
	if inputDirectory == "" {
		panic(errors.New("视频目录不能为空!" ))
	}
	// 获取目录下所有视频文件
	fileNames :=getvideoFile(inputDirectory)
	fmt.Println("该目录视频列表如下:")
	for i:=0; i<len(fileNames); i++ {
		fmt.Println(strconv.Itoa(i) + ":\t" + fileNames[i])
	}
	var videoFileName string
	for{
		videoFileNoStr := getParam("请输入视频编号：")
		videoFileNo,err := strconv.Atoi(videoFileNoStr)
		if err == nil {
			if videoFileNo >= 0 && videoFileNo < len(fileNames) {
				videoFileName = fileNames[videoFileNo]
				break
			}
		}
	}


	// 手动输入视频文件名
	//videoFileName := getParam("请输入视频文件名（带后缀）：")
	//if videoFileName == "" {
	//	panic(errors.New("视频文件名不能为空!" ))
	//}

	// 根据编号选择
	inFilePath := inputDirectory  + videoFileName
	flag :=isPathExist(inFilePath)
	if !flag {
		panic(errors.New("视频文件路径:"+inFilePath+" 有误!" ))
	}

	// 开始时间
	startTime := getParam("请输入开始时间(格式: 00:00:00[.000]，为空为从头开始)：")
	if startTime == "" {
		startTime = DefaultTime
	}

	//结束时间
	endTime := getParam("请输入结束时间(格式: 00:00:00[.000]，为空为视频时长):")

	// 输出路径
	outputDirectory := getParam("请输入视频输出目录(为空则为当前视频目录)：")
	//输出文件名
	if outputDirectory == "" {
		outputDirectory = inputDirectory
	}else{
		outputDirectory = outputDirectory + string(os.PathSeparator)
	}
	outFilePath := getNotExistFilePath(outputDirectory,videoFileName)
	fmt.Println("请等待...")
	var cmd *exec.Cmd
	if endTime =="" {
		cmd = exec.Command("ffmpeg.exe",
			"-ss",startTime,
			"-i",inFilePath,
			"-c","copy",
			"-copyts",outFilePath)
	}else{
		cmd = exec.Command("ffmpeg.exe",
			"-ss",startTime,
			"-i", inFilePath ,
			"-to",endTime,
			"-c","copy",
			"-copyts",outFilePath)
	}
	if err := cmd.Run(); err != nil {
		panic("剪辑命令异常：" + err.Error())
	}

	fmt.Println("done!")

}

/**
获取参数
 */
 func getParam(fmtString string) string {
	 input := bufio.NewScanner(os.Stdin)
	 fmt.Println(fmtString)
	 input.Scan()
	 return input.Text()
 }

 /**
 判断路径是否存在
  */
func isPathExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

/**
获取目录下未重复文件编号
例如 abp-001 则从abp-001(1)开始递增，直到寻到第一个不存在的文件名
 */
 func getNotExistFilePath(path string, name string) string{
 	/**
 	todo 有个很奇葩的问题， 字符串拼接后，如果拼接了"(0)" ,就少最后一个字符，去掉(就没事)
 	 */
	 fileSuffix :=  path2.Ext(name)//获取文件后缀
	 filenameOnly := strings.TrimSuffix(name, fileSuffix)//获取文件名
	 for i := 0; i < 10000; i++ {
	 	tempPath := path + filenameOnly +strconv.Itoa(i)+fileSuffix
		 if !isPathExist(tempPath){
		 	return tempPath
		 }
	 }
	 // 呵呵
	 panic(errors.New("???" ))
 }

/**判断一个文件是否是视频文件*/
var isvideoReg = regexp.MustCompile("\\.[\\w]+$")
func isvideo(name string)bool{
	suf := isvideoReg.FindString(name)
	for _,temp := range videoSuf{
		if strings.EqualFold(temp,suf){
			return true
		}
	}
	return false
}

/**
读取目录下视频文件
 */
 func getvideoFile(dirPath string)(fileNames []string){
	 // 读取目录
	 fileInfo,err := ioutil.ReadDir(dirPath)
	 if err != nil{
	 	panic("读取目录失败:" +err.Error())
	 }
	 for _,i := range fileInfo{
		 if !i.IsDir() && isvideo(i.Name()){
			 fileNames = append(fileNames, i.Name())
		 }
	 }
	 return
 }

/**
创建文件
 */
 func createFile(name string) {
	 file,err:=os.Create(name)
	 defer file.Close()
	 if err!=nil{
		 panic("创建文件失败:"+err.Error())
	 }
 }

 /**
 删除文件
  */
  func deleteFile(mame string) {
	  err := os.Remove(mame)
	  if err != nil {
		  panic("删除文件失败:"+err.Error())
	  }
  }



