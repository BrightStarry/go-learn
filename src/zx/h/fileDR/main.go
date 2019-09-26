package main

import (
	"bufio"
	"os"
	"fmt"
	"io/ioutil"
	"errors"
	"regexp"
	"strings"
)

/**
番号文件去重
 */

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error:" ,err)
		}
		getParam("exit")
	}()

	dir1 := getParam("目录1:")
	// 获取文件
	fileNames1 := getAllFileName(dir1)
	// 获取番号
	fileNo1 := make([]no,len(fileNames1))
	for i,temp := range fileNames1 {
		fileNo1[i] = getNO(temp)

	}

	for{
		dir2 := getParam("目录2:")
		fileNames2 := getAllFileName(dir2)
		fileNo2 := make([]no,len(fileNames2))
		for i,temp := range fileNames2 {
			fileNo2[i] = getNO(temp)
		}

		fmt.Println("以下为重复番号:")
		// 比较
		for _, i1 := range fileNo1 {
			// no为空时退出该次比较
			if i1.isNull(){
				continue
			}
			for _,i2 := range fileNo2 {
				if i1.equals(&i2){
					fmt.Println(i1.pre +" " +i1.suf + "\t")
					break
				}
			}
		}
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
			fileNames = append(fileNames, i.Name())
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

type no struct {
	pre string
	suf string
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
	fc2 = "fc2"
	ZERO = "0"
)

/**
提取番号
 */
var getNOReg = regexp.MustCompile("^([A-Za-z\\d]+|[\\d]+)[-_\\s]?([\\d]+)")
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
		return no{FC2, strings.TrimLeft(temp[0],ZERO)}
	}

	// 处理其他番号
	temp := getNOReg.FindStringSubmatch(name)
	// 格式错误，直接返回空对象
	if len(temp) < 3 {
		fmt.Println("格式错误:" +name)
		return
	}
	return no{temp[1], strings.TrimLeft(temp[2],ZERO)}
}