package main

import (
	"io/ioutil"
	"os"
	"bufio"
	"fmt"
	"strings"
	"regexp"
	"errors"
	"path/filepath"
	"zx/h/getNameByNumber/util"
	"zx/h/getNameByNumber/myLog"
	"github.com/PuerkitoBio/goquery"
	"sync"
)


var keyword = []string{":","*","?","/","\\","\"","<",">","|"}
// 默认网址s
var defaultUrls = []string{"http://www.f37b.com","https://www.dmmbus.co"}
// 默认网址后缀
// javlibrary http://www.m34z.com/cn/vl_searchbyid.php?keyword=番号
// dmmbus https://www.dmmbus.co/番号
var defaultUrlSufs = []string{"/cn/vl_searchbyid.php?keyword=","/"}
// 默认获取片名的方法
var defaultGetNameFuns = []func(*goquery.Document)string {
	func(doc *goquery.Document)string {
		return util.GetTextBySelector(doc,"#video_title > h3 > a")
	},
		func(doc *goquery.Document)string {
			return util.GetTextBySelector(doc,"body > div.container > h3")
		},
}


/**
通过番号从 网站 获取片名
 */
func main() {
	defer func() {
		if err := recover(); err != nil {
			myLog.Error("error:" ,err)
		}
		getParam("exit")
	}()

	// 获取外部参数
	//importExtParam()
	dir := getParam("请输入目录:")

	// 获取文件
	fileNames := getAllFileName(dir)
	// 解析番号
	fileNo := make([]no,len(fileNames))
	for i,temp := range fileNames {
		fileNo[i] = getNO(filepath.Base(temp))
	}
	waitGroup := sync.WaitGroup{}
	for i:=0;i<len(fileNames);i++ {
		// 解析出目录名 和 文件名
		currentdir, fileName := filepath.Split(fileNames[i])
		// 当前番号
		no := fileNo[i]
		waitGroup.Add(1)
		// 异步处理
		go func() {
			defer func(){
				waitGroup.Done()
				if err := recover(); err != nil {
					myLog.Error("当前文件名:" + fileName + ",内部错误:",err)
				}
			}()

			// 如果包含空格，则跳过
			if strings.Contains(fileName, " ") {
				myLog.Warn("文件名有空格，可能已经重命名:" + fileName)
				return
			}


			if no.isNull() {
				myLog.Warn("番号解析有误:" + fileName)
				return
			}
			// 给番号添加0
			suf := no.suf
			switch len(suf) {
			case 1:
				suf = "00" + suf
			case 2:
				suf = "0" + suf
			}
			// 正常番号
			number := no.pre + "-" +suf

			avName := ""
			var newFileName string
			for j:=0; j < len(defaultUrls) && avName == ""; j++ {
				// 当前url
				currentUrl := defaultUrls[j] + defaultUrlSufs[j] + number
				document,err1 := util.GetOfDocument(currentUrl)
				if err1 != nil {
					continue
				}
				// 获取片名
				avName = defaultGetNameFuns[j](document)
				if avName == "" {
					continue
				}
				newFileName = avName + "~" + fileName
			}

			if avName == "" {
				myLog.Error("当前文件名:" + fileName + ",获取番号失败")
				return
			}

			// 替换特殊字符
			for i := 0; i< len(keyword);i++  {
				newFileName = strings.Replace(newFileName,keyword[i]," ",-1)
			}

			err2 := os.Rename(currentdir + fileName, currentdir + newFileName)
			if err2 != nil {
				myLog.Error("err! 重命名失败:" + err2.Error() +  "\t newFileName:" + newFileName  )
				return
			}
			myLog.Info("success! " + newFileName)
		}()
	}
	waitGroup.Wait()
	myLog.Info("done!")


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
			return
		}
		return no{FC2, strings.TrimLeft(temp[0],ZERO)}
	}

	// 处理其他番号
	temp := getNOReg.FindStringSubmatch(name)
	// 格式错误，直接返回空对象
	if len(temp) < 3 {
		return
	}
	return no{temp[1], strings.TrimLeft(temp[2],ZERO)}
}


/**
	获取外部参数
 */
func importExtParam() {

}