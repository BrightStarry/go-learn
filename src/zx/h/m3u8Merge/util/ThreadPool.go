package util

import (
	"github.com/pkg/errors"
	"zx/h/m3u8Merge/myLog"
)

/**
简单的线程池
 */
const (
	SUCCESS ="success"
)
// 线程池对象
type ThreadPool struct {
	// 获取任务
	Queue chan *Job
	// 发送结果
	Result chan Result
	//线程数
	ThreadNum int
	// 工作函数
	RunFunc func([]interface{}) error
}



// 初始化， 线程数
func (threadPool *ThreadPool) Init(threadNumber int,
	runFunc func([]interface{})error) *ThreadPool{
	threadPool.Queue = make(chan *Job,threadNumber)
	threadPool.Result= make(chan Result,threadNumber)
	threadPool.ThreadNum = threadNumber
	threadPool.RunFunc = runFunc
	return threadPool
}

// 异步启动, 只有当所有线程非空闲时才会阻塞
func (threadPool *ThreadPool) Start() error{
	if threadPool.Queue == nil || threadPool.Result == nil  {
		return errors.New("线程池未初始化,参数异常.")
	}
	for i:=0; i < threadPool.ThreadNum; i++{
		go threadPool.run(i)
	}
	return nil
}


// 单个任务， 任务id
func (threadPool *ThreadPool) run(threadId int) {
	for job := range threadPool.Queue{
		defer func() {
			if err := recover(); err != nil{
				myLog.Info("threadId:%d,queueId:%d,error:%v",threadId,job.Id,err)
				threadPool.Result <-Result{false,"未知异常.",job.FuncArgs}
			}
		}()
		myLog.Info("threadId:%d,queueId:%d,start.",threadId,job.Id)
		// 执行
		err := threadPool.RunFunc(job.FuncArgs)
		// 失败
		if err != nil {
			myLog.Info("threadId:%d,queueId:%d,error:%v",threadId,job.Id,err)
			threadPool.Result <- Result{false,err.Error(),job.FuncArgs}
			break
		}
		// 成功
		myLog.Info("threadId:%d,queueId:%d,success.",threadId,job.Id)
		threadPool.Result <- Result{true,SUCCESS,job.FuncArgs}
	}
}

// 任务入队
func (threadPool *ThreadPool) Put(id int,funcArgs []interface{}) {
	threadPool.Queue <- &Job{id,funcArgs}
}

// 获取结果
func (threadPool *ThreadPool) Take() Result {
	return <-threadPool.Result
}

// 关闭队列
func (threadPool *ThreadPool) Close()  {
	close(threadPool.Queue)
}
// 任务
type Job struct{
	// id
	Id int
	// 任务参数
	FuncArgs []interface{}
}

// 结果
type Result struct {
	Success bool
	Message string
	Data []interface{}
}