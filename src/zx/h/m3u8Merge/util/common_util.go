package util

import (
	"bufio"
	"os"
	"fmt"
	"io/ioutil"
	"errors"
	"os/exec"
	"github.com/spf13/viper"
	"io"
	"zx/h/m3u8Merge/myLog"
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
		panic("读取配置异常:" + err.Error())
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
func GetFileName(dirPath string) (fileNames []string, err error) {
	// 读取目录
	fileInfo, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	for _, i := range fileInfo {
		if i.Name() == "System Volume Information" {
			continue
		}
		fileNames = append(fileNames, dirPath+string(os.PathSeparator)+i.Name())
	}
	return
}

/*遍历出文件夹中所有文件名*/
func GetAllFileName(dirPath string) (fileNames []string) {
	// 读取目录
	fileInfo, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(errors.New("目录读取异常:" + err.Error()))
	}
	for _, i := range fileInfo {
		if i.IsDir() {
			if i.Name() == "System Volume Information" {
				continue
			}
			temp := GetAllFileName(dirPath + string(os.PathSeparator) + i.Name())
			fileNames = append(fileNames, temp...)
		} else {
			fileNames = append(fileNames, dirPath+string(os.PathSeparator)+i.Name())
		}
	}
	return
}

/**
判断文件是否存在
 */
func FileIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	// 没有异常表示，存在
	return err == nil || os.IsExist(err)
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
执行命令并输出到控制台，异步
 */
//func StartCMD(name string, arg ...string) (err error) {
// cmd := exec.Command(name,arg...)
// //显示运行的命令
// fmt.Println(cmd.Args)
// // 获取子进程标准输出
// stdout, _ := cmd.StdoutPipe()
// defer stdout.CloseQueue()
// // 执行命令
// err = cmd.Start()
// if err != nil {
// 	return
// }
// // 读取子进程
// reader := bufio.NewReader(stdout)
// for {
//	line, err := reader.ReadString('\n')
//	if err != nil   {
//		if io.EOF != err {
//			myLog.Warn("执行cmd:%s,读取输出异常:%v",cmd.Args,err)
//		}
//		break
//	}
//	// 转换CMD的编码为GBK
//	lineBytes, _ := ioutil.ReadAll(transform.NewReader(
//		bytes.NewReader([]byte(line)),
//		simplifiedchinese.GBK.NewDecoder(),
//	))
//	// 将子进程的内容输出
//	fmt.Println(string(lineBytes))
//
// }
// return
//}

/**
执行命令并输出到控制台，异步
*/
func StartCMDBase(name string, arg ...string) (cmd *exec.Cmd, reader *bufio.Reader, err error) {
	// 将cmd命令行输出从GBK改到UTF-8，防止中文乱码
	//exec.Command("chcp", "65001").Run()

	cmd = exec.Command(name, arg...)
	//显示运行的命令
	fmt.Println(cmd.Args)
	// 获取子进程标准输出
	stdout, _ := cmd.StdoutPipe()
	//defer stdout.CloseQueue()
	// 执行命令
	err = cmd.Start()
	if err != nil {
		return
	}
	// 读取子进程
	reader = bufio.NewReader(stdout)
	return
}

/**
执行命令并输出到控制台，异步
*/
func StartCMD(name string, arg ...string) (err error) {
	cmd, reader, err := StartCMDBase(name, arg...)
	defer cmd.Wait()
	if err != nil {
		return err
	}
	// 处理输出
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if io.EOF == err {
				break
			}
			myLog.Warn("执行cmd:%s,读取输出异常:%v", cmd.Args, err)
		}
		fmt.Println(line)
	}
	return
}

/**
执行命令并输出到日志，异步
*/
func StartCMDToLog(loaPath, name string, arg ...string) (err error) {
	cmd, reader, err := StartCMDBase(name, arg...)
	defer func() {
		if err2 := cmd.Wait();err2 != nil {
			myLog.Error("命令执行失败,命令:%v,异常:%v",cmd.Args,err2)
		}
	}()
	if err != nil {
		return err
	}
	// 创建日志文件
	logFile, err := os.OpenFile(loaPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return errors.New("打开日志文件失败:" + err.Error())
	}
	defer logFile.Close()

	// TODO 可增加缓冲区
	writer := bufio.NewWriter(logFile)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			writer.Flush()
			if io.EOF == err {
				break
			}
			myLog.Warn("执行cmd:%s,读取输出异常:%v", cmd.Args, err)
		}
		if _,err = writer.WriteString(line);err != nil {
			myLog.Warn("记录单条日志异常,日志内容:%s,异常:%v",line,err)
		}
	}
	writer.Flush()
	return
}

/**
执行命令并输出到控制台，同步
*/
func RunCMD(flag, name string, arg ...string) (err error) {
	cmd := exec.Command(name, arg...)
	buf, err := cmd.Output()
	fmt.Printf("%s: %s\n", flag, buf)
	return err
}
