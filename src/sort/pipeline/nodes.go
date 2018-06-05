package pipeline

import (
	"sort"
	"encoding/binary"
	"math/rand"
	"bufio"
)

/*节点*/

/*将任意数组放入一个通道，返回一个单向接收通道*/
func ArraySource(a ...int) <-chan int {
	out := make(chan int)
	// 通道传输通常要另起线程，防止阻塞
	go func() {
		defer close(out)
		for _,v:= range a{
			out <- v
		}
	}()
	return out
}

/*从接收通道接收数组，进行内部排序，然后重新放入一个接收通道*/
func InMemSort(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		// 读取到内存
		var a []int
		for v := range in{
			a = append(a,v)
		}
		// 排序
		sort.Ints(a)

		// 输出
		for _,v := range a{
			out <- v
		}
	}()
	return out
}


/*归并,从两个通道中读取数据，按序输出到另一个通道中*/
func Merge( in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		// 从通道1获取数据
		v1,ok1 := <- in1
		// 从通道2获取数据
		v2,ok2 := <- in2
		// 只要不是两个通道都关闭了
		for ok1 || ok2{
			// 如果ok2被关闭了或者（ok1没关闭并且v1小于等于v2）才输出v1
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				// 然后再取下个值
				v1,ok1 = <- in1
			}else {
				out <- v2
				v2,ok2 = <- in2
			}
		}
	}()
	return out
}

/*读取元数据到一个通道*/
func ReaderSource(reader *bufio.Reader) <-chan int{
	out := make(chan int)
	go func() {
		defer close(out)
		// 缓冲区，每次读取 64bit
		buffer := make([]byte,8)
		for{
			n,err :=reader.Read(buffer)
			if n > 0 {
				// 将64bit转为int
				v := int( binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil {
				break
			}
		}
	}()
	return out
}

/*从指定通道输出数据*/
func WriterSink(writer *bufio.Writer,in <-chan int) {
	for v:= range in {
		// 缓冲区
		buffer := make([]byte,8)
		// int转[]byte
		binary.BigEndian.PutUint64(buffer,uint64(v))
		// 输出
		writer.Write(buffer)
	}
	writer.Flush()

}

/*随机数数据源*/
func RandomSource(count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i:=0;i<count;i++{
			out <- rand.Int()
		}
	}()
	return out
}