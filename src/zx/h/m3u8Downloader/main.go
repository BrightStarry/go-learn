package main

import (
	"github.com/spf13/viper"
	"os"
	"github.com/pkg/errors"
	"fmt"
	"strings"
	"path/filepath"
	"zx/h/m3u8Downloader/pool"
	"zx/h/m3u8Downloader/util"
	"zx/h/m3u8Downloader/aria2"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

/**
m3u8下载器

1.读取等待队列txt
2.运行任务开启协程至最大值
3.任务运行成功后写入成功txt，失败则写入失败队列，并终止整个程序(通常不会发生)


 */
var config Config

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("error:", err)
		}
		util.GetParam("exit")
	}()
	//start()
}

/**
初始化方法
 */
func init1() {
	// 初始化配置
	initConfig()
	// 清理临时文件
	if config.IsClean == IsCleanTrue {
		log.Info("开始清理临时文件...")
		if err := util.BatchDelFile([]string{config.taskLogPathPre,config.logPath,
			config.successTaskPath,config.errorTaskPath, config.backTaskQueuePath,
			config.downloadingTaskPath,
		}); err != nil {
			log.Panicln("清理临时文件异常:",err)
		}
		log.Info("临时文件清理完毕...")
	}

	// 初始化日志
	initLog()
	log.WithFields(log.Fields{
		"config":config,
	}).Info("当前参数.")

	// 创建目录
	if err :=os.MkdirAll(config.taskLogPathPre,0777);err != nil {
		log.Panicln("日志目录创建失败:",err)
	}

	// 备份队列
	taskQueueBytes,err := ioutil.ReadFile(config.taskQueuePath)
	if err != nil {
		log.Panicln("备份队列失败:",err)
	}
	if err = ioutil.WriteFile(config.backTaskQueuePath, taskQueueBytes, 0777);err != nil {
		log.Panicln("备份队列失败:",err)
	}

}

/**
配置日志
 */
func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	// 设置日志级别为warn以上
	log.SetLevel(log.InfoLevel)
	if config.LogType == LogTypeFile {
		file, err := os.OpenFile(config.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Panicln("日志文件创建失败:" + err.Error())
		}
		log.SetOutputFile(file)
	}else {
		log.SetOutput(os.Stdout)
	}
}

/**
 主进程
 */
func start() {
	// 读取任务队列
	taskQueue, err := readTaskQueue(config.ThreadNum)
	if err != nil {
		log.Panicln("读取任务队列失败:" + err.Error())
	}
	// 任务队列为空，结束
	if taskQueue == nil || len(taskQueue)==0 {
		log.Info("任务队列为空.")
		return
	}
	// 任务队列小于线程数
	if len(taskQueue) < config.ThreadNum {
		config.ThreadNum = len(taskQueue)
		log.Warn("当前任务队列任务数:",len(taskQueue),",小于任务线程数,将临时设置线程数等于当前任务数.")
	}
	// 线程池
	threadPool := pool.ThreadPool{}
	threadPool.Init(config.ThreadNum, func(args []interface{}) (err error) {
		defer func() {
			if err := recover(); err != nil {
				err = errors.New(fmt.Sprintf("未知异常:%v",err) )
				return
			}
		}()
		return runTask(args[0].(Task))
	})
	// 启动
	threadPool.Start()

	getNewTaskChannel := make(chan int)

	errSize := 0
	successSize :=0
	// 处理结果
	threadPool.ProcessResult(func(r pool.Result) {
		task := r.Data[0].(Task)
		// 将任务从下载中队列移除
		if err := delLineByTaskQueue(config.downloadingTaskPath,task.m3u8Path);err != nil {
			log.Panicln("操作下载中队列文件失败:",err)
			os.Exit(-1)
		}
		if r.Success {
			// 写入成功日志，失败则停止程序
			if err := util.AppendTxt(config.successTaskPath,task.m3u8Path);err != nil {
				log.WithField("taskName",task.taskName).Panicln("写入成功队列异常:",err)
				os.Exit(-1)
			}
			log.WithField("taskName",task.taskName).Info("任务成功.", )
			successSize += 1
		}else{
			// 任务失败
			if err := util.AppendTxt(config.errorTaskPath,task.m3u8Path);err != nil {
				log.WithField("taskName",task.taskName).Panicln("写入异常队列异常:",err)
				os.Exit(-1)
			}
			log.WithFields(log.Fields{
				"taskName": task.taskName,
				"message": r.Message,
			}).Error("任务失败!具体原因查看任务对应日志:",task.loaPath)
			errSize += 1
		}
		// 发送获取新任务信号
		getNewTaskChannel <- 1
	})
	taskId := 0
	// 任务入队
	for _, v := range taskQueue {
		taskId += 1
		threadPool.Put(taskId,[]interface{}{v})
	}

	// 当前已停止的线程数
	currentStopTaskSize := 0
	// 主进程等待信号
	for i := range getNewTaskChannel {
		i = i
		log.Info("线程空闲，开始读取任务队列.")
		// 读取任务队列
		taskQueue, err := readTaskQueue(1)
		if err != nil {
			log.Panicln("读取任务队列失败:" + err.Error())
		}
		// 任务队列为空
		if taskQueue == nil || len(taskQueue)==0 {
			currentStopTaskSize += 1
			log.Info("任务队列为空,当前进行中任务数:",config.ThreadNum - currentStopTaskSize,",等待进行中任务结束后，程序将会关闭...")
			// 如果所有线程都在等待
			if currentStopTaskSize == config.ThreadNum {
				log.Info("所有任务执行完毕.总任务数:",taskId,",成功任务数:",successSize,",失败任务数:",errSize)
				threadPool.CloseResultQueue()
				threadPool.CloseTaskQueue()
				return
			}
			continue
		}
		// 添加新任务
		log.Info("读取到新任务:",taskQueue[0].taskName)
		taskId += 1
		threadPool.Put(taskId,[]interface{}{taskQueue[0]})
	}

}

/*
获取新任务
 */
 func getNewTask() {
	 defer func() {
		 if err := recover(); err != nil{
		 	log.Panicln("获取新任务线程内部异常:",err)
		 }
	 }()



 }

// 任务
type Task struct {
	m3u8Path string // m3u8文件路径
	taskName string // 任务名，默认为m3u8文件名
	loaPath  string // 任务日志文件存储位置,默认为根目录+log+任务名
	outDir   string // 任务文件保存位置,默认为根目录+ts+任务名
}



/**
读取任务队列
 */
func readTaskQueue(size int) ([]Task, error) {
	if size <= 0 {
		return nil, nil
	}
	// 读取
	taskTxt, err := util.ReadTxt(config.taskQueuePath)
	if err != nil {
		return nil, err
	}
	if taskTxt == "" {
		return nil, nil
	}
	taskTxt = strings.TrimSpace(taskTxt)
	// 按行分割任务
	m3u8Paths := strings.Split(taskTxt, "\r\n")
	// 返回对象
	var taskQueue []Task
	for _, item := range m3u8Paths {
		if item == "" {
			continue
		}
		item = strings.TrimSpace(item)
		// 文件名带后缀
		fileName := filepath.Base(item)
		// 后缀
		suf := filepath.Ext(item)
		if suf != ".m3u8" {
			return nil, errors.New("任务队列解析失败，任务路径后缀非m3u8.")
		}
		// 任务名
		taskName := strings.TrimSuffix(fileName, suf)
		taskQueue = append(taskQueue, Task{
			item,
			taskName,
			config.taskLogPathPre + string(os.PathSeparator) + taskName + ".log",
			config.taskFilePathPre + string(os.PathSeparator) + taskName,
		})
	}
	// 如果要获取的任务数大于等于当前任务数，设置任务队列为空
	if size >= len(m3u8Paths) {
		if err = util.WriteTxt(config.taskQueuePath, "");err!= nil {
			log.Panicln("写入任务队列文件异常:",err)
		}
		return taskQueue[:len(m3u8Paths)], nil
	}
	// 否则， 将剩余任务回写
	if err = util.WriteTxt(config.taskQueuePath, strings.Join(m3u8Paths[size:], "\r\n"));err!= nil {
		log.Panicln("写入任务队列文件异常:",err)
	}
	return taskQueue[:size], nil

}

/**
	单个任务
 */
func runTask(task Task) error {
	// 写入下载中队列
	if err := util.AppendTxt(config.downloadingTaskPath,task.m3u8Path);err != nil {
		log.Fatalln("写入下载中队列失败:",err)
		os.Exit(-1)
	}
	thisLog := log.WithFields(log.Fields{
		"taskName":task.taskName,
	})
	thisLog.Info("开始下载.")
	// 执行命令
	err := aria2.RunAria2(task.loaPath, task.m3u8Path, task.outDir, config.Aria2Path,config.Aria2Args)
	if err != nil {
		return errors.New(fmt.Sprintf("任务:%s,运行异常:%v", task.taskName, err))
	}
	thisLog.Info("下载完毕，准备验证是否下载成功.")
	success, err := aria2.IsSuccessByLog(task.loaPath)
	// 有异常或者失败
	if err != nil || !success {
		return errors.New(fmt.Sprintf("任务:%s,任务失败,异常:%v", task.taskName, err))
	}
	// 成功
	return nil
}

/**
 配置
 */
type Config struct {
	ThreadNum       int    // 同时下载数
	RootDir         string //根目录
	Aria2Path       string //aria2路径
	LogType         int // 日志类型:  0:标准输出; 1:文件
	IsClean         int //是否清理临时文件,包括 上次运行日志;上次成功队列;上次失败队列 0:不清理; 1:清理
	Aria2Args		[]string // aria2下载参数设置
	successTaskPath string // 成功任务文件路径
	errorTaskPath   string // 失败任务文件路径
	taskQueuePath   string // 任务队列路径
	downloadingTaskPath string // 正在下载任务队列路径
	backTaskQueuePath string // 备份队列路径
	taskLogPathPre  string // 任务日志前缀
	taskFilePathPre string // 任务文件前缀
	logPath         string // 日志路径

}

const(
	// 日志类型
	LogTypeStdout = 0
	LogTypeFile   = 1

	// 是否清理
	IsCleanTrue = 1
	IsCleanFalse = 0
)

func initConfig() {
	configReader := viper.New()
	configReader.SetConfigName("m3u8Downloader")
	configReader.AddConfigPath("./")
	configReader.SetConfigType("yaml")
	if err := configReader.ReadInConfig(); err != nil {
		log.Panicln("读取配置异常:" + err.Error())
	}
	if err := configReader.Unmarshal(&config); err != nil {
		log.Panicln("读取配置异常:" + err.Error())
	}
	// TODO
	rootDir, _ := os.Getwd()// 获取当前绝对路径
	config.RootDir = rootDir
	//config.RootDir = `E:\m3u8Downloader`

	config.successTaskPath = config.RootDir + string(os.PathSeparator) + "successTask.txt"
	config.errorTaskPath = config.RootDir + string(os.PathSeparator) + "errorTask.txt"
	config.taskQueuePath = config.RootDir + string(os.PathSeparator) + "taskQueue.txt"
	config.downloadingTaskPath = config.RootDir + string(os.PathSeparator) + "taskQueue_downloading.txt"
	config.backTaskQueuePath = config.RootDir + string(os.PathSeparator) + "taskQueue_back.txt"
	config.taskLogPathPre = config.RootDir + string(os.PathSeparator) + "log"
	config.taskFilePathPre = config.RootDir + string(os.PathSeparator) + "ts"
	config.logPath = config.RootDir + string(os.PathSeparator) + "m3u8Downloader.log"

}
/**
	从文件中移除某一行
 */
func delLineByTaskQueue(path,m3u8Path string) (err error) {
	content, err := util.ReadTxt(path)
	if err != nil {
		return
	}
	if content == "" {
		return errors.New("下载中队列数据异常丢失")
	}
	m3u8Paths := strings.Split(content, "\r\n")// 切割成行
	flag := -1
	// 找到元素下标
	for i, item := range m3u8Paths {
		if item == "" {
			continue
		}
		if item == m3u8Path {
			flag = i
			break
		}
	}
	if flag == -1 {
		return errors.New("下载中队列数据异常,未找到要删除的元素")
	}
	// 删除元素
	m3u8Paths = append(m3u8Paths [:flag], m3u8Paths [flag+1:]...)
	// 回写数据
	if err = util.WriteTxt(path, strings.Join(m3u8Paths, "\r\n"));err!= nil {
		return errors.New(fmt.Sprintf("写入下载中任务队列文件异常:%v",err))
	}
	return nil
}