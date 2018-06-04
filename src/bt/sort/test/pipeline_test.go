package test

import (
	"testing"
	"bt/sort/pipeline"
	"fmt"
	"os"
	"bufio"
)

/*测试类*/

/*测试*/
func TestA(t *testing.T) {
	const filename = "small.in"
	const size = 100000000
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
	in := pipeline.ReaderSource(bufio.NewReader(file2))

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
