package pipeline

import (
	"net"
	"bufio"
	"os"
	"strconv"
)

/*网络版 外部排序*/


/*
	创建网络版通道
	chunkCount：将一个文件分为几块
*/
func CreateNetworkPipeline(filename string,fileSize,chunkCount int) <-chan int {
	// 记录每次开启的地址
	var sortAddr []string

	// 初始化时间
	Init()

	// 计算每一块的大小
	chunkSize := fileSize / chunkCount
	for i:=0; i < chunkCount; i++{
		file,err := os.Open(filename)
		if err!= nil {
			panic(err)
		}
		// 设置下一次读写的位置，offset为相对偏移量，whence为相对位置，0相对文件开头；1相对当前位置；2相对结尾位置； 返回新的偏移量（相对开头）和异常
		file.Seek(int64(i * chunkSize),0)

		// 读取该文件的指定块为通道
		source := ReaderSource(bufio.NewReader(file),chunkSize)

		// 每次端口不同
		addr := ":" + strconv.Itoa(7000 +i)
		// 进行内部排序,并
		// 此时就是开启一个服务，等待别人连接后，将排序后的数据输出该它
		NetworkSink(addr,InMemSort(source))
		sortAddr = append(sortAddr,addr)
	}


	// 模拟若干台机器
	// 连接到刚才输出到网络的每个端口，读取输出
	var sortResults []<-chan int
	for _,addr := range sortAddr {
		sortResults = append(sortResults,
			NetworkSource(addr))
	}
	// 归并
	return MergeN(sortResults...)
}


/*
	网络版输出
	建立连接，包装conn为输出器，交由之前的输出方法，输出到对方连接
*/
func NetworkSink(addr string, in <-chan int) {
	listener,err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	// 异步
	go func() {
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		writer := bufio.NewWriter(conn)
		WriterSink(writer,in)
	}()
}

/*
	从网络中读取
	连接到对应地址，读取信息，包装为channel，
*/
func NetworkSource(addr string) <-chan int {
	out := make(chan int)
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}
		r := ReaderSource(bufio.NewReader(conn), -1)
		for v := range r {
			out <- v
		}
		defer close(out)
	}()
	return out
}
