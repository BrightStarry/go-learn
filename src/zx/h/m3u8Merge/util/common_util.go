package util

import (
	"bufio"
	"os"
	"fmt"
	"io/ioutil"
	"errors"
	"io"
	"golang.org/x/text/transform"
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"github.com/spf13/viper"
)




/**
读取配置文件.默认从根目录
 */
func ReadConfig(name string) {
	// 读取配置文件
	viper.SetConfigName(name)
	viper.AddConfigPath("./")
	//viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic("读取配置异常:" +err.Error())
	}
}
/**
获取参数
 */
func GetParam(fmtString string) string {
	input := bufio.NewScanner(os.Stdin)
	fmt.Println(fmtString)
	input.Scan()
	return input.Text()
}

/**
获取目录文件，不遍历子目录
 */
 func GetFileName(dirPath string)(fileNames []string,err error) {
	 // 读取目录
	 fileInfo,err := ioutil.ReadDir(dirPath)
	 if err != nil {
		return
	 }
	 for _,i := range fileInfo{
		 if i.Name() == "System Volume Information"{
			 continue
		 }
		 fileNames = append(fileNames, dirPath + string(os.PathSeparator) + i.Name())
	 }
	 return
 }

/*遍历出文件夹中所有文件名*/
func GetAllFileName(dirPath string)(fileNames []string){
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
			temp := GetAllFileName(dirPath + string(os.PathSeparator) + i.Name())
			fileNames = append(fileNames, temp...)
		}else{
			fileNames = append(fileNames, dirPath + string(os.PathSeparator) + i.Name())
		}
	}
	return
}

/**
判断文件是否存在
 */
func FileIsExist(filePath string)  bool{
	_,err := os.Stat(filePath)
	// 没有异常表示，存在
	return  err ==nil || os.IsExist(err)
}
/**
判断文件是否存在,不存在则创建
 */
func CreateFileByFileIsExist(filePath string) (err error) {
	if FileIsExist(filePath) {
		return
	}
	return os.MkdirAll(filePath, os.ModePerm)
}

/**
执行命令并输出到控制台
 */
 func RunCMD(name string, arg ...string) (err error) {
	 cmd := exec.Command(name,arg...)
	 //显示运行的命令
	 fmt.Println(cmd.Args)
	 // 获取子进程标准输出
	 stdout, _ := cmd.StdoutPipe()
	 // 执行命令
	 err = cmd.Start()
	 if err != nil {
	 	return
	 }
	 // 读取子进程
	 reader := bufio.NewReader(stdout)
	 for {
		 line, err := reader.ReadString('\n')
		 if err != nil || io.EOF == err {
			 break
		 }
		 // 转换CMD的编码为GBK
		 reader := transform.NewReader(
			 bytes.NewReader([]byte(line)),
			 simplifiedchinese.GBK.NewDecoder(),
		 )
		 d, _ := ioutil.ReadAll(reader)
		 // 将子进程的内容输出
		 print(string(d))
	 }
	 return
 }

