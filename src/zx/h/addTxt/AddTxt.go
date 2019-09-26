package main

import (
	"os"
	"io/ioutil"
	"bufio"
	"fmt"
	"errors"
	"strings"
	"github.com/gpmgo/gopm/modules/log"
	"path"
)

const(
	TXT = " 原档"
)

/**
给某目录下的所有文件增加文字

文件夹名字不能包含空格
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
		currentdir, fileName := path.Split(fileNames[i])

		// 如果已经重命名,跳过
		if strings.Contains(fileName, TXT) {
			continue
		}

		splitArr :=strings.Split(fileName," ")
		if len(splitArr) < 2 {
			log.Error("err! 文件名格式有误,当前文件名:" + fileName)
			continue
		}
		newFileName := ""
		for j := 0; j < len(splitArr); j++ {
			if j == 0 {
				newFileName += splitArr[j] + TXT
			}else{
				newFileName += " " + splitArr[j]
			}
		}


		err := os.Rename(currentdir + fileName, currentdir + newFileName)
		if err != nil {
			log.Error("err! 重命名失败:" + err.Error() +  "\t old:" + currentdir + fileName  )
			continue
		}
		fmt.Println("success! old:" + fileName + "\t" + "new:" + newFileName)
	}

	log.Info("done!")
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

