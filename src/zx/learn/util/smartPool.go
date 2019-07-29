package main

import (
	"io/ioutil"
	"os"
	"errors"
)

func main() {

}


/*遍历出文件夹中所有文件,*/
func getAllFileName(dirPath string)(fileNames []string,err error){
	// 读取目录
	files,err := ioutil.ReadDir(dirPath)
	for _,i := range files {
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