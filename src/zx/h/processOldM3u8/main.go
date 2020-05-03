package main

import (
	log "github.com/sirupsen/logrus"
	"zx/h/processOldM3u8/util"
	"io/ioutil"
	"os"
	"errors"
	"strconv"
	"encoding/base64"
	"strings"
	"path/filepath"
)
const(
	KEYBIN = "key.bin"
)
/**
处理 0.m3u8和key.bin格式的
 */
func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("error:", err)
		}
		util.GetParam("exit")
	}()

	rootDir := util.GetParam("输入根目录:")
	// 获取根目录下下一级文件夹
	allDirs := getDirName(rootDir,false)
	if len(allDirs) < 1 {
		log.Panicln("根目录下无文件夹...")
	}

	// 循环每个目录
	for _, dirItem := range allDirs {
		processDir2(dirItem)
	}

}
/**
处理单个目录
处理 0.m3u8和key.bin格式的
 */
 func processDir1(dir string) {
 	dir += string(os.PathSeparator)
 	log.Println("当前处理目录:",dir)
 	 // 获取子文件夹名
	 dirNames := getDirName(dir,true)
	 if len(dirNames) < 1 {
	 	log.Panicln("该目录下无子目录...")
	 }
	 // 循环子文件夹，每个文件夹名即 dmm番号
	 for _, number := range dirNames {
		 log.Println("当前处理番号:",number)
		 // 当前目录
	 	currentDir := dir + number + string(os.PathSeparator)
	 	// 当前目录 + 番号前缀
	 	currentDirAndNumber := currentDir+number

	 	// 读取key
		 keyBytes,err := ioutil.ReadFile(currentDir + KEYBIN)
		 if  err != nil {
			 log.Panicln("读取key异常:",err)
		 }
		 if  len(keyBytes) != 16 {
			 log.Panicln("读取key异常，key不是16个字节:",currentDir + KEYBIN)
		 }
		 keyBase64 := key2Base64(keyBytes)
		// 处理单个m3u8
		if util.FileIsExist(currentDir+ "0.m3u8") {
			if err:=os.Rename(currentDir+ "0.m3u8",currentDirAndNumber+ ".m3u8");err != nil {
				log.Panicln("重命名异常:",err)
			}
			// 替换m3u8
			replaceKeyByM3u8(currentDirAndNumber+ ".m3u8",keyBase64)
			if err:=os.Rename(currentDir+ KEYBIN,currentDirAndNumber+ ".key");err != nil {
				log.Panicln("重命名异常:",err)
			}
		}else{
			// 处理多段m3u8
			for j:=1;;j++{
				// 当前目录 + 番号前缀 + 番号片段索引
				currentDirAndNumberIndex := currentDirAndNumber + "_part"+ strconv.Itoa(j)
				if util.FileIsExist(currentDir+ strconv.Itoa(j) +".m3u8") {
					if err:=os.Rename(currentDir+ strconv.Itoa(j) +".m3u8",
						currentDirAndNumberIndex+".m3u8");err != nil {
						log.Panicln("重命名异常:",err)
					}
					// 替换m3u8
					replaceKeyByM3u8(currentDirAndNumberIndex+".m3u8",keyBase64)
					if err := ioutil.WriteFile(currentDirAndNumberIndex+".key",keyBytes,0666);err != nil {
						log.Panicln("写入key异常:",err)
					}

				}else{
					// 不存在了，删除key.bin 跳出循环
					if err := os.Remove(currentDir + KEYBIN);err != nil {
						log.Panicln("删除key.bin失败:",err)
					}
					break
				}
			}
		}
	 }

	 /**
	 提取所有m3u8和key
	  */
	 allFiles := getAllFileName(dir)
	 for _, item := range allFiles {
		 ext := filepath.Ext(item)
		 if ext == ".m3u8"||ext == ".key" {
		 	// 获取文件名
		 	fileName := filepath.Base(item)
			if err:=os.Rename(item,dir + fileName );err!= nil {
				log.Panicln("移动文件失败:",err)
			}
		 }
	 }
	 // 删除其他文件
	 delFiles := getAllFileName(dir)
	 for _, item := range delFiles {
		 ext := filepath.Ext(item)
		 if ext == ".m3u8"||ext == ".key" {
		 	continue
		 }
		 if err := os.Remove(item);err != nil {
			 log.Panicln("删除文件失败:",err)
		 }
	 }

	 // 删除文件夹
	 delDirs := getDirName(dir,false)
	 for _, item := range delDirs {
		 if err := os.Remove(item);err != nil {
		 	log.Panicln("删除文件夹失败:",err)
		 }
	 }
 }

/**
处理单个目录
处理 cid.m3u8和key.bin格式的
*/
func processDir2(dir string) {
	dir += string(os.PathSeparator)
	log.Println("当前处理目录:",dir)
	// 获取子文件夹名
	dirNames := getDirName(dir,true)
	if len(dirNames) < 1 {
		log.Panicln("该目录下无子目录...")
	}
	// 循环子文件夹，每个文件夹名即 dmm番号
	for _, number := range dirNames {
		log.Println("当前处理番号:",number)
		// 当前目录
		currentDir := dir + number + string(os.PathSeparator)
		// 当前目录 + 番号前缀
		currentDirAndNumber := currentDir+number

		// 读取key
		keyBytes,err := ioutil.ReadFile(currentDir + KEYBIN)
		if  err != nil {
			log.Panicln("读取key异常:",err)
		}
		if  len(keyBytes) != 16 {
			log.Panicln("读取key异常，key不是16个字节:",currentDir + KEYBIN)
		}
		keyBase64 := key2Base64(keyBytes)
		// 处理单个m3u8
		if util.FileIsExist(currentDirAndNumber +".m3u8") {
			// 替换m3u8
			replaceKeyByM3u8(currentDirAndNumber+ ".m3u8",keyBase64)
			if err:=os.Rename(currentDir+ KEYBIN,currentDirAndNumber+ ".key");err != nil {
				log.Panicln("重命名异常:",err)
			}
		}else{
			// 处理多段m3u8
			for j:=1;;j++{
				// 当前目录 + 番号前缀 + 番号片段索引
				currentDirAndNumberIndex := currentDirAndNumber + "_part"+ strconv.Itoa(j)
				if util.FileIsExist(currentDir+ strconv.Itoa(j) +".m3u8") {
					if err:=os.Rename(currentDir+ strconv.Itoa(j) +".m3u8",
						currentDirAndNumberIndex+".m3u8");err != nil {
						log.Panicln("重命名异常:",err)
					}
					// 替换m3u8
					replaceKeyByM3u8(currentDirAndNumberIndex+".m3u8",keyBase64)
					if err := ioutil.WriteFile(currentDirAndNumberIndex+".key",keyBytes,0666);err != nil {
						log.Panicln("写入key异常:",err)
					}

				}else{
					// 不存在了，删除key.bin 跳出循环
					if err := os.Remove(currentDir + KEYBIN);err != nil {
						log.Panicln("删除key.bin失败:",err)
					}
					break
				}
			}
		}
	}

	/**
	提取所有m3u8和key
	 */
	allFiles := getAllFileName(dir)
	for _, item := range allFiles {
		ext := filepath.Ext(item)
		if ext == ".m3u8"||ext == ".key" {
			// 获取文件名
			fileName := filepath.Base(item)
			if err:=os.Rename(item,dir + fileName );err!= nil {
				log.Panicln("移动文件失败:",err)
			}
		}
	}
	// 删除其他文件
	delFiles := getAllFileName(dir)
	for _, item := range delFiles {
		ext := filepath.Ext(item)
		if ext == ".m3u8"||ext == ".key" {
			continue
		}
		if err := os.Remove(item);err != nil {
			log.Panicln("删除文件失败:",err)
		}
	}

	// 删除文件夹
	delDirs := getDirName(dir,false)
	for _, item := range delDirs {
		if err := os.Remove(item);err != nil {
			log.Panicln("删除文件夹失败:",err)
		}
	}
}

 /**
 读取指定M3u8,替换key为base64，并重新写入
  */
  func replaceKeyByM3u8(m3u8Path,keyBase64 string) {
	  m3u8Bytes,err := ioutil.ReadFile(m3u8Path)
	  if err != nil {
		  log.Panicln("读取m3u8异常:",err)
	  }
	  m3u8Str := strings.Replace(string(m3u8Bytes),KEYBIN,"base64:" + keyBase64,1)
	  if err = ioutil.WriteFile(m3u8Path, []byte(m3u8Str), 0666);err != nil {
		  log.Panicln("写入m3u8异常:",err)
	  }
  }

/**
读取key并转为base64
 */
  func key2Base64(bytes []byte)string {
	  return base64.StdEncoding.EncodeToString(bytes)
  }

/*遍历出文件夹中下一级目录*/
func getDirName(dirPath string,isOnlyName bool)(fileNames []string){
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
			if isOnlyName {
				fileNames = append(fileNames, i.Name())
			}else{
				fileNames = append(fileNames, dirPath + string(os.PathSeparator) + i.Name())
			}
		}
	}
	return
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
