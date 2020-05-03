package util

import (
	"os"
	"io/ioutil"
	"bufio"
	"fmt"
	"os/exec"
	"io"
	log "github.com/sirupsen/logrus"
	"errors"
	"github.com/spf13/viper"
)


/**
初始化配置文件
 */
func InitConfig(configFileName string,config interface{}) {
	configReader := viper.New()
	configReader.SetConfigName(configFileName)
	configReader.AddConfigPath("./")
	configReader.SetConfigType("yaml")
	if err := configReader.ReadInConfig(); err != nil {
		log.Panicln("读取配置异常:" + err.Error())
	}
	if err := configReader.Unmarshal(config); err != nil {
		log.Panicln("读取配置异常:" + err.Error())
	}
}

/**
 文件 是否存在
 */
func FileIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	// 没有异常表示，存在
	return err == nil || os.IsExist(err)
}

/**
追加写入
 */
func AppendTxt(path, line string) (err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.WriteString(line + "\r\n")
	return
}

/**
覆盖写入
 */
func WriteTxt(path, txt string) (err error) {
	return ioutil.WriteFile(path, []byte(txt), 0777)
}

/**
读取文件
 */
func ReadTxt(path string) (result string, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return string(bytes), err
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
执行命令并输出到控制台，异步
*/
func StartCMDBase(name string, arg ...string) (cmd *exec.Cmd, reader *bufio.Reader, err error) {
	// 将cmd命令行输出从GBK改到UTF-8，防止中文乱码
	exec.Command("chcp", "65001").Run()
	cmd = exec.Command(name, arg...)
	//显示运行的命令
	log.WithFields(log.Fields{"cmd":cmd.Args}).Info("执行命令.")
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
			log.Warn("执行cmd:", cmd.Args, ",读取输出异常:", err)
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
		// 确认命令执行完成，获取异常信息
		if err2 := cmd.Wait(); err2 != nil {
			err = err2
			return
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

	// TODO 缓冲区暂时设置为32KB
	writer := bufio.NewWriterSize(logFile, 32*1024)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			writer.Flush()
			if io.EOF == err {
				break
			}
			log.Warn("执行cmd:", cmd.Args, ",读取输出异常:", err)
		}
		if _, err = writer.WriteString(line); err != nil {
			log.Warn("记录单条日志异常,日志内容:", line, ",异常:", err)
		}
	}
	writer.Flush()
	return
}

/**
批量删除
 */
func BatchDelFile(paths []string) (err error) {
	if len(paths) == 0 {
		return
	}
	for _, item := range paths {
		if err = os.RemoveAll(item); err != nil {
			return
		}
	}
	return
}


