package pool

import (
	"github.com/pkg/errors"
	"fmt"
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
	TaskQueue chan *Job
	// 发送结果
	ResultQueue chan Result
	//线程数
	ThreadNum int
	// 工作函数
	RunFunc func([]interface{}) error
}



// 初始化， 线程数
func (threadPool *ThreadPool) Init(threadNumber int,
	runFunc func([]interface{})error) *ThreadPool{
	threadPool.TaskQueue = make(chan *Job,threadNumber)
	threadPool.ResultQueue = make(chan Result,threadNumber)
	threadPool.ThreadNum = threadNumber
	threadPool.RunFunc = runFunc
	return threadPool
}

// 异步启动, 只有当所有线程非空闲时才会阻塞
func (threadPool *ThreadPool) Start() error{
	if threadPool.TaskQueue == nil || threadPool.ResultQueue == nil  {
		return errors.New("线程池未初始化,参数异常.")
	}
	for i:=0; i < threadPool.ThreadNum; i++{
		go threadPool.run(i)
	}
	return nil
}


// 单个任务， 任务id
func (threadPool *ThreadPool) run(threadId int) {
	for job := range threadPool.TaskQueue {
		threadPool.run1(threadId,job)
	}
}

/**
单个任务，内循环
 */
func (threadPool *ThreadPool) run1(threadId int,job *Job) {
	defer func() {
		if err := recover(); err != nil{
			threadPool.ResultQueue <-Result{false,fmt.Sprint("未知异常: ",err),job.FuncArgs}
		}
	}()
	// 执行
	if err := threadPool.RunFunc(job.FuncArgs);err != nil {
		threadPool.ResultQueue <- Result{false,err.Error(),job.FuncArgs}
		return
	}
	// 成功
	threadPool.ResultQueue <- Result{true,SUCCESS,job.FuncArgs}
}

// 任务入队
func (threadPool *ThreadPool) Put(id int,funcArgs []interface{}) {
	threadPool.TaskQueue <- &Job{id,funcArgs}
}

// 获取结果
func (threadPool *ThreadPool) Take() Result {
	return <-threadPool.ResultQueue
}

// 关闭队列
func (threadPool *ThreadPool) CloseTaskQueue()  {
	close(threadPool.TaskQueue)
}

// 关闭结果获取队列
func  (threadPool *ThreadPool)CloseResultQueue() {
	close(threadPool.ResultQueue)
}

// 结果处理线程，异步
func  (threadPool *ThreadPool) ProcessResult(processResultFun func(r Result)) {
	go func() {
		for result := range threadPool.ResultQueue {
			processResultFun(result)
		}
	}()
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