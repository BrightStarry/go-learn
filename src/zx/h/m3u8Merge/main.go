package main

import (
	"zx/h/m3u8Merge/util"
	"path/filepath"
	"io/ioutil"
	"os"
	"sync"
	"zx/h/m3u8Merge/myLog"
	"errors"
	"strings"
	"strconv"
	"github.com/spf13/viper"
)

/**
	m3u8文件下载器
 */
 // 程序根目录
var rootDir string
// 下载的ts目录
 var rootTsDir string
 // mp4视频结果目录
var resultDir  string
// key目录
var keyDir string
func main() {
	defer func() {
		if err := recover(); err != nil {
			myLog.Error("error:%v" ,err)
		}
		util.GetParam("exit")
	}()

	util.ReadConfig("m3u8Merge")

	// key目录
	keyDir = viper.GetString("keyDir")
	// 获取rootTSDir
	rootDir = viper.GetString("rootDir") + string(os.PathSeparator)
	rootTsDir = rootDir + "ts"+string(os.PathSeparator)
	resultDir = rootDir +"result" +string(os.PathSeparator)
	// 是否转码
	isTranscoding := viper.GetString("isTranscoding")
	// 线程数
	//threadNum := viper.GetInt("threadNum")

	tsDirs,err := util.GetFileName(rootTsDir)
	if err != nil {
		panic("目录读取异常:" + err.Error())
	}
	if tsDirs==nil || len(tsDirs) == 0 {
		panic("目录为空.")
	}
	myLog.Info("待处理目录:\n" + strings.Join(tsDirs,"\n"))

	for i,tsDir := range tsDirs {
		myLog.Info("开始处理第%d个视频,目录:%s",i+1,tsDir)
		// 番号
		number :=filepath.Base(tsDir)
		// 解密后输出路径
		//outPath := rootTsDir+number +"-out"+string(os.PathSeparator)
		//// 解密
		//decrypt(tsDir,outPath,threadNum)
		//// 合并
		//merge(outPath)
		tsOutPath := resultDir + number +".ts"
		decryptAndMerge(tsDir,tsOutPath)
		if isTranscoding != "0" {
			// 转码
			transcodingAndRemoveTemp(number,tsOutPath)
		}


		myLog.Info("第%d个视频处理成功,目录:%s",i+1,tsDir)
	}
	myLog.Info("处理完成目录:\n" + strings.Join(tsDirs,"\n"))
}

/**
解密 并合并
 */
 func decryptAndMerge(dir,outPath string) {
	 myLog.Info("开始解密合并." )
	 // 获取ts文件名
	 filePaths := util.GetAllFileName(dir)
	 filePaths = shellSort(filePaths)// 排序
	 // 读取key
	 keyFilePath := keyDir +string(os.PathSeparator) +filepath.Base(dir) + ".key"
	 keyBytes,err := ioutil.ReadFile(keyFilePath)
	 if err!= nil {
		 panic("读取key异常：" + err.Error())
	 }
	 outFile, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC,0644)
	 defer outFile.Close()
	 if err != nil {
		 panic("创建文件异常:" + err.Error())
	 }

	 bytesChannel := make(chan []byte,32)
	 go func() {
		 for _,item := range filePaths {
			 tempBytes,err := ioutil.ReadFile(item)
			 if err != nil {
				 panic("解密合并异常:" + err.Error())
			 }
			 bytesChannel <- tempBytes
		 }
		 close(bytesChannel)
	 }()

	 var allBytes []byte
	 i:=1
	 for tempBytes:= range bytesChannel {
		 tempBytes = util.AESDecrypt(tempBytes,keyBytes)
		 allBytes = append(allBytes, tempBytes...)
		 // 每x个文件，或最后一个文件 ，全部写入
		 if i % 64== 0 || i == len(filePaths){
			 _,err = outFile.Write(allBytes)
			 if err != nil {
				 panic("解密合并异常:" + err.Error())
			 }
			 allBytes = []byte{}
			 myLog.Info("处理中...")
		 }
		 i++
	 }

	 myLog.Info("解密合并成功" )
	 // 成功则删除该文件夹
	 if err =os.RemoveAll(dir); err != nil {
		 myLog.Error("删除文件夹失败:%v",err)
	 }
 }

/**
不转码时，把result.ts移动到result目录，并删除out文件夹
 */
 func processNotTranscoding(outPath,number string) {
 	movePath := resultDir + number + ".ts"
	if err :=util.StartCMD("cmd","/c","move",
		outPath+"result.ts",movePath);err != nil {
			panic("移动文件失败:"+ err.Error())
	}

	 if util.FileIsExist(movePath) {//如果ts文件存在
		 if err := os.RemoveAll(outPath);err != nil {
			 myLog.Error("删除outPath失败:%v",err)
		 }
	 }
 }

/**
	转码
*/
func transcodingAndRemoveTemp(number string,tsOutPath string) {
	myLog.Info("开始转码.")

	vedioPath := resultDir + number +".mp4"
	util.StartCMD("cmd","/c",
		"ffmpeg",
		"-i",tsOutPath,
		"-c","copy",
		"-bsf:a","aac_adtstoasc",
		"-y", vedioPath)
	myLog.Info("转码成功.")

	// 删除文件
	if util.FileIsExist(vedioPath) {//如果mp4文件存在
		if err := os.RemoveAll(tsOutPath);err != nil {
			myLog.Error("删除tsOutPath失败:%v",err)
		}
	}
}

/**
合并
 */
 func merge(outPath string) {

	 myLog.Info("开始合并.")
	 util.StartCMD("cmd","/c", "copy","/B",
		 outPath+"*.ts",  outPath+"result.ts")

	 //tempPath := outPath+"temp.ts"
	 //for i,offset := 0,256;i <=len(filePaths);i += offset{
	 //	if i + 64 >= len(filePaths){
		//	offset = len(filePaths) - i - 1
		//}
		// tempStr := strings.Join(filePaths[i:i+offset],"+")
		// if i != 0 {
		//	 tempStr = tempPath + "+" + tempStr
		// }
		// util.StartCMD("cmd","/c",
		//	 "copy","/B",
		//	 tempStr,  tempPath)
	 //}


	 /**
	 分段合并成多个临时文件
	  */
	 //index := 0
	 //for i,offset:= 0,128;i <=len(filePaths);i += offset{
		//index+=1
	 //	if i + 64 >= len(filePaths){
	 //		offset = len(filePaths) - i - 1
	 //	}
		// tempStr := strings.TrimRight(strings.Join(filePaths[i:i+offset],"+"),"+")
		// tempPath := outPath+"temp"+ strconv.Itoa(index) +".ts"
	 //
		// util.StartCMD("cmd","/c",
		// 	 "copy","/B",
		// 	 tempStr,  tempPath)
	 //}

	 /**
	 将多个临时文件合并
	  */
	 //allTempPath:= ""
	 //for i:=1; i<=index;i++{
		// allTempPath += outPath+"temp"+ strconv.Itoa(i) +".ts" + "+"
	 //}
	 //allTempPath = strings.TrimSuffix(allTempPath,"+")
	 //util.StartCMD("cmd","/c", "copy","/B",
		// allTempPath,  outPath+"result.ts")
	 //// 成功后删除
	 //for i:=1; i<=index;i++{
	 //	if err := os.Remove(outPath+"temp"+ strconv.Itoa(i) +".ts" ); err != nil {
	 //		log.Error("临时文件删除失败:%s",outPath+"temp"+ strconv.Itoa(i) +".ts" )
		//}
	 //}



	//// 以追加模式打开文件，当文件不存在时生成文件
	//file, err := os.OpenFile(outPath+"result.ts", os.O_RDWR|os.O_CREATE|os.O_TRUNC,0644)
	//defer file.CloseQueue()
	//if err != nil {
	//	 panic("合并异常:" + err.Error())
	//}
	//var allBytes []byte
	//for i,item := range filePaths {
	//	tempBytes,err := ioutil.ReadFile(item)
	//	if err != nil {
	//		panic("合并异常:" + err.Error())
	//	}
	//	allBytes = append(allBytes, tempBytes...)
	//	// 每x个文件，或最后一个文件 ，全部写入
	//	if (i!=0 && (i % 512) == 0) || i == len(filePaths)-1{
	//		_,err = file.Write(allBytes)
	//		if err != nil {
	//			panic("合并异常:" + err.Error())
	//		}
	//		allBytes = []byte{}
	//		myLog.Info("追加数据成功...")
	//	}
	//}

	 //var allBytes []byte
	 //for _,item := range filePaths {
		// tempBytes,err := ioutil.ReadFile(item)
		// 	if err != nil {
		// 		panic("合并异常:" + err.Error())
		// 	}
		// 	allBytes = append(allBytes,tempBytes...)
	 //}
	 //if err := ioutil.WriteFile(outPath+"result.ts",allBytes,0666);err != nil {panic("合并异常:" + err.Error())}

	 myLog.Info("合并成功.")
 }

/**
解密
 */
 func decrypt(dir ,outPath string,threadNum int) {
	 myLog.Info("开始解密." )
	 // 获取ts文件名
	 filePaths := util.GetAllFileName(dir)
	 // 读取key
	 keyFilePath := keyDir +string(os.PathSeparator) +filepath.Base(dir) + ".key"
	 keyBytes,err := ioutil.ReadFile(keyFilePath)
	 if err!= nil {
		 panic("读取key异常：" + err.Error())
	 }


	 err = util.CreateFileByFileIsExist(outPath)
	 if err!= nil {
		 panic("创建文件异常：" + err.Error())
	 }

	 // 设置线程池
	 threadPool := util.ThreadPool{}
	 threadPool.Init(threadNum, func(args []interface{}) error {
		 filePath := args[0].(string)
		 // 读取加密视频
		 videoBytes,err := ioutil.ReadFile(filePath)
		 if err!= nil {
			 return errors.New("读取加密视频异常：" + err.Error())
		 }
		 decryptBytes := util.AESDecrypt(videoBytes,keyBytes)

		 /**
		 修改文件名 ,将1-9变成001，002，003, 100-999变成0100,0999
		   */
		 fileName := filepath.Base(filePath)
		 // 获取文件编号
		 number := strings.TrimSuffix(strings.Split(fileName,"_")[2],".ts")
		 switch len(number) {
		 case 1:
			 number = "000" + number
		 case 2:
			 number = "00" + number
		 case 3:
		 	number = "0" + number
		 }
		 err = ioutil.WriteFile(outPath + number + ".ts",decryptBytes,0666)
		 if err != nil {
			 return errors.New("写入解密视频异常：" + err.Error())
		 }
		 return nil
	 })
	 // 启动
	 threadPool.Start()
	 // 获取结果,必须放在入队前，否则会死锁
	 errResults := make([]util.Result,0)
	 waitGroup := sync.WaitGroup{}
	 waitGroup.Add(1)
	 go func() {
		 defer waitGroup.Done()

		 // 获取结果
		 for i:=0;i< len(filePaths);i++{
			 result := threadPool.Take()
			 if !result.Success {
				 errResults = append(errResults,result)
			 }
		 }
	 }()

	 //任务入队
	 for i,item := range filePaths{
		 threadPool.Put(i,[]interface{}{item})
	 }
	 // 关闭
	 threadPool.CloseQueue()
	 // 等待获取结果
	 waitGroup.Wait()
	 // 打结果
	 if len(errResults)> 0 {
		 myLog.Error("失败结果如下:%v" ,errResults)
		 panic("解密失败.")
	 }else{
		 myLog.Info("解密成功" )
		 // 成功则删除该文件夹
		 if err =os.RemoveAll(dir); err != nil {
		 	myLog.Error("删除文件夹失败:%v",err)
		 }
	 }
 }

 /**
 排序，只针对filePaths
  */
func shellSort(arr []string) []string {
	subFunc := func(s string)(string) {
		tempArr := strings.Split(filepath.Base(s),"_")
		return strings.TrimSuffix(tempArr[2],".ts")
	}

	equalsFunc := func(s1,s2 string)bool {
		i1, err := strconv.Atoi(subFunc(s1))
		i2, err := strconv.Atoi(subFunc(s2))
		if err != nil {
			panic(err)
		}
		return i1 > i2
	}

	size := len(arr)
	in := size / 2 // 起始增量,为一半元素
	// 循环到增量为0 (1/2=0)
	for in >= 1 {
		// 循环 in次，循环所有子数组
		for i1:=0; i1<in; i1++{
			// 对每个子数组进行插入排序
			// 增量为in时， 0，0+in,0+in+in为一个子数组， 1,1+in,1+in+in是一个子数组
			for i2:= i1+in;i2 < size;i2+=in {
				temp := arr[i2]
				i3 := i2 - in
				for ;i3 >= i1 && equalsFunc(arr[i3],temp); i3 -= in {
					arr[i3+in] = arr[i3]
				}
				arr[i3+in] =temp
			}
		}
		in = in / 2
	}
	return arr
}