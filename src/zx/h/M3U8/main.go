package main

import (
	"fmt"
	"os"
	"bufio"
	"errors"
	"sync"
	"io/ioutil"
	"sync/atomic"
	"strings"
	"regexp"
	"encoding/base64"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error:" ,err)
		}
		getParam("exit")
	}()

	// 线程数
	threadNum := 1024
	dir := getParam("请输入当前m3u8和key文件的目录：")
	if dir == "" {
		panic("目录不能为空!")
	}

	// 读取所有文件
	fileInfos := dirents(dir)

	// 接收文件的通道
	fileChannel := make(chan os.FileInfo)

	// 向通道发送单个文件
	go func() {
		for _, item := range fileInfos {
			fileChannel <- item
		}
		close(fileChannel)
	}()

	count := int64(0)

	// 保存 keyName 和 其base64编码的key
	var keyMap sync.Map

	base64Encoder :=base64.StdEncoding
	reg := regexp.MustCompile(`(.+URI=["]?)[\w@:\\.]*(["]?.+)`)

	// 并发修改
	waitGroup := sync.WaitGroup{}
	for i := 0; i < threadNum; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for fileInfo := range fileChannel {
				atomic.AddInt64(&count,1)

				filePath := dir + string(os.PathSeparator) + fileInfo.Name()

				// 保存key文件内容
				if strings.HasSuffix(fileInfo.Name(),"key") {
					bytes,err := ioutil.ReadFile(filePath)
					if err!= nil {
						panic(errors.New("读取文件异常:" + err.Error()))
					}
					keyMap.LoadOrStore(strings.TrimSuffix(fileInfo.Name(),".key"),"base64:"  + base64Encoder.EncodeToString(bytes))
					continue
				}
				// 修改m3u8文件
				if strings.HasSuffix(fileInfo.Name(),"m3u8") {
					// 读取key文件
					keyPath := strings.Replace(filePath,"m3u8","key",1)
					bytes,err := ioutil.ReadFile(keyPath)
					if err!= nil {
						panic(errors.New("读取文件异常:" + err.Error()))
					}
					key := "base64:"  + base64Encoder.EncodeToString(bytes)


					// 文件内容
					fileStr := readFileToString(filePath)
					// 替换
					fileStr = reg.ReplaceAllString(fileStr,"${1}" + key + "${2}")
					writeFile(filePath,fileStr)
					fmt.Printf("文件:%s,修改成功.\n",filePath)
					continue
				}
				panic("遍历目录异常，包含了其他元素")
			}
		}()
	}
	waitGroup.Wait()
	fmt.Printf("成功,遍历总文件个数:%d \n",count )


}

// dirents 返回 dir 目录中的条目
func dirents(dir string) []os.FileInfo {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(errors.New("遍历文件异常：" + err.Error() ))
		return nil
	}
	return fileInfos
}

/**
读取文件
 */
 func readFileToString(fileName string)string{
 	result,err :=ioutil.ReadFile(fileName)
 	if err != nil{
 		panic("读取文件失败:"+ err.Error())
	}
	return string(result)
 }

 /**
 写入文件
  */
  func writeFile(fileName,fileContent string) {
		err := ioutil.WriteFile(fileName,[]byte(fileContent),0666)
	  if err != nil {
		  panic("写入文件失败:"+ err.Error())
	  }
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
