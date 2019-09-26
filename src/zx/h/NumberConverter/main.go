package main

import (
	"os"
	"bufio"
	"fmt"
	"io/ioutil"
	"errors"
	"strings"
	"regexp"
	"path/filepath"
	"github.com/gpmgo/gopm/modules/log"
)

/**
番号转换
dmm的番号转换成常规番号
 */

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error:" ,err)
		}
		getParam("exit")
	}()

	dir := getParam("请输入目录:")
	fileNames := getAllFileName(dir)

	for i:=0;i<len(fileNames);i++{
		// 解析出目录名 和 文件名
		currentdir, fileName := filepath.Split(fileNames[i])
		// 如果已经转换，跳过
		if strings.Contains(fileName, "-") {
			log.Warn("已经转换:" + fileName)
			continue
		}

		// 获取文件后缀
		ext  := filepath.Ext(fileName)
		// 没有后缀的文件名
		onlyFileName := strings.TrimSuffix(fileName,ext)

		/**
		开始转换
		 */
		numberReg := regexp.MustCompile("(h_)?(\\d*)((t28)|([A-Za-z]+))0*([\\d]+)([A-Za-z])?(_part)?(\\d)?")
		tempNumber := numberReg.FindStringSubmatch(fileName)

		if len(tempNumber) < 10 {
			log.Warn("无法解析:" + fileName)
			continue
		}

		pre := tempNumber[3]
		suf := tempNumber[6]
		index := tempNumber[9]
		switch len(suf) {
		case 1:
			suf = "00" + suf
		case 2:
			suf = "0" + suf
		}
		newFileName := pre + "-" + suf + "~" + index + "~" + onlyFileName + ext
		err := os.Rename(currentdir + fileName, currentdir + newFileName)
		if err != nil {
			log.Error("err! 重命名失败:" + err.Error() +  "\t newFileName:" + newFileName  )
			continue
		}
		fmt.Println("success! " + newFileName)
	}
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



/**
获取参数
 */
func getParam(fmtString string) string {
	input := bufio.NewScanner(os.Stdin)
	fmt.Println(fmtString)
	input.Scan()
	return input.Text()
}

