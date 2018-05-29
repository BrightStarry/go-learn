package main

import "time"

/*并发编程*/

func main() {
	// 单线程
	//loop()
	//loop()

	// 多线程
	go loop() // 此处另起线程
	loop()
	// 为防止 异步线程还未开启，主函数就结束，稍微等待
	time.Sleep(time.Second)// 睡眠1s


	// channel 通道
	// 创建通道, 其中的 string表示这个通道传输的消息类型
	//var channel1 = make(chan string )
	//或
	channel2 := make(chan string )

	// 发消息
	go func(message string) {
		channel2 <- message
	}("xxx")

	// 取消息
	println(<-channel2)

	// 默认的通道都是无缓冲的阻塞通道
	//channel2 <- "dfddf"// 也就是说，执行该语句后，如果不从该通道中取出消息，会一直阻塞在此处（此处会抛出异常，因为所有线程(goroutines)都sleep了）
	// 读取也是同理
	// 因此，也可以将这种阻塞作为一种便利，例如当某个线程返回某个值后，停止主线程

	// 开启异步写入和异步读取
	go asyncWrite(channel2)
	go asyncReaad(channel2)
	// 显式地关闭信道,关闭后仍可以读取通道中缓冲的数据，只是不能再写入
	close(channel2)
	time.Sleep(time.Hour)


	// 死锁例子 fatal error: all goroutines are asleep - deadlock!
	//ch := make(chan int)
	//<- ch // 阻塞main goroutine, 信道c被锁

	// 如下也是死锁
	// 其中主线等ch1中的数据流出，ch1等ch2的数据流出，但是ch2等待数据流入，两个goroutine都在等，也就是死锁。
	/*
	var ch1 chan int = make(chan int)
	var ch2 chan int = make(chan int)

	func say(s string) {
		fmt.Println(s)
		ch1 <- <- ch2 // ch1 等待 ch2流出的数据
	}

	func main() {
		go say("hello")
		<- ch1  // 堵塞主线
	}
	*/


	// 有缓冲的阻塞通道，当消息个数到达指定个数后，才进行读取
	//ch := make(chan int, 3)

	// 并且，无论有无缓冲，通道的读取都是有序的，线程安全的
}


// 异步写入
func asyncWrite(channel chan string) {
	// 一直写入数据
	for{
		channel <- "1"
		time.Sleep(time.Second)
	}
}

// 异步读取
func asyncReaad(channel chan string) {
	// 一直尝试读取
	//for{
	//	print(<-channel)
	//}
	// 一直尝试读取，直到通道关闭
	for message := range channel{
		print(message)
		//if len(channel) <= 0 { // 如果现有数据量为0，跳出循环
		//	break
		//}
	}
}

// 模拟循环
func loop() {
	for i := 0; i<10; i++{
		print(i)
	}
	println()
}