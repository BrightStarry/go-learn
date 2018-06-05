package test

import (
	"testing"
	"fmt"
	"os"
	"bufio"
	"sort/pipeline"
)

/*测试类*/

/*
	外部排序
	1.分块读取一个文件，将每一块作为一个chanel，进行内部排序后，归并
	2. 写入文件
	3. 打印该输出文件
*/
func TestExternalSort(t *testing.T) {
	// 创建通道
	p := createPipeline("small.in",512,4)
	// 写入文件
	writeToFile(p,"small.out")
	// 打印文件
	printFile("small.out")
}

/*打印文件*/
func printFile(filename string) {
	file,err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(bufio.NewReader(file),-1)
	for v := range p {
		fmt.Println(v)
	}
}

/*
	写入文件
	两个defer是先进后出的，所以是先Flush，在close
*/
func writeToFile(p <-chan int, filename string) {
	file,err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	 defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer,p)
}

/*
	创建通道
	chunkCount：将一个文件分为几块

	TODO 将其中打开的file close，需要返回 *file
*/
func createPipeline(filename string,fileSize,chunkCount int) <-chan int {
	// 排序后的通道集合
	var sortResults []<-chan int

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
		source := pipeline.ReaderSource(bufio.NewReader(file),chunkSize)
		// 进行内部排序,并放入 排序通道集合
		sortResults = append(sortResults,pipeline.InMemSort(source))
	}
	// 归并
	return pipeline.MergeN(sortResults...)
}

/*测试随机生成数字，写入文件*/
func TestA(t *testing.T) {
	const filename = "small.in"
	const size = 64
	// 创建文件
	file,err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 随机生成
	out := pipeline.RandomSource(size)
	// 将其输出到创建出来的文件
	pipeline.WriterSink(bufio.NewWriter(file),out)

	// 读取该文件
	file2,err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file2.Close()
	in := pipeline.ReaderSource(bufio.NewReader(file2),-1)

	count := 0
	for  v:= range in {
		fmt.Println(v)
		count++
		if count >= 100 {
			break
		}
	}

}

/*测试单机归并*/
func TestMerge(t *testing.T) {
	// 将数组输入到一个接收通道
	// 接收一个接收通道，对通道数据进行排序，再输出一个接收通道
	p1 := pipeline.InMemSort(pipeline.ArraySource(124,34,434,6,34,43,1))

	// 再输出一个
	p2 := pipeline.InMemSort(pipeline.ArraySource(43,546,767,55,6,5656,56))

	// 归并两个通道的数据
	r := pipeline.Merge(p1,p2)

	// 一直读取，直到通道关闭
	for num := range r {
		fmt.Println(num)
	}
}
