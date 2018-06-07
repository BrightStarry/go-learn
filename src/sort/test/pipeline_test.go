package test

import (
	"testing"
	"fmt"
	"os"
	"bufio"
	"sort/pipeline"
)

/*测试类*/




/**
	测试网络排序
 */
func TestNetworkSort(t *testing.T) {
	// 创建通道
	p := pipeline.CreateNetworkPipeline("small.in",512,4)
	//// 写入文件
	pipeline.WriteToFile(p,"small.out")
	//// 打印文件
	pipeline.PrintFile("small.out")
}

/*
	外部排序
	1.分块读取一个文件，将每一块作为一个chanel，进行内部排序后，归并
	2. 写入文件
	3. 打印该输出文件
*/
func TestExternalSort(t *testing.T) {
	// 创建通道
	p := pipeline.CreatePipeline("small.in",800000000,100)
	// 写入文件
	pipeline.WriteToFile(p,"small.out")
	// 打印文件
	pipeline.PrintFile("small.out")
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
