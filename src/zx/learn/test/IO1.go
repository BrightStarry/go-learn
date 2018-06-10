package main

import (
	"io"
	"strings"
	"fmt"
	"os"
	"bufio"
	"errors"
	"encoding/binary"
)

/*
	io相关
*/

func main() {
	// 读取字符读取器例子
	//sampleReadFromString()

	// 读取标准输入流例子
	//sampleReadStdin()


	// 从文件中读取
	//sampleReadFile()


	// 缓冲读取例子
	//sampleBuffer()

	// 统计文件行数
	//countFileLine()


	// 读取bmp图片的头信息
	readHeaderFromBmp()
}



/*从输入器中读取指定长度的字节*/
func ReadFrom(reader io.Reader,num int)([]byte , error) {
	// 缓冲切片
	buf := make([]byte,num)

	n,err := reader.Read(buf)
	if n>0{
		return buf[:n],nil
	}

	// 如果不大于0，则可能发生异常
	return buf,err
}

/*从字符读取器中 读取字节的例子*/
func sampleReadFromString() {
	// 通过传入的字符串，返回一个读取器(指针)
	reader := strings.NewReader("from string")
	bytes,_:=ReadFrom(reader,12)
	fmt.Println("读取到的数据:",bytes)
}

/*从标准输入中读取字节*/
func sampleReadStdin() {

	fmt.Println("请输入：")
	bytes,_ := ReadFrom(os.Stdin,11)

	fmt.Println("读取到:",bytes)
}

/*从文件中读取字节*/
func sampleReadFile() {
	file,_ := os.Open("H:\\goSpace\\go-learn\\bin\\jump.exe")
	defer file.Close()
	bytes,_ := ReadFrom(file,100)
	fmt.Println("读取到:",bytes)

}

/*缓冲读取例子*/
func sampleBuffer() {
	// 创建读取器
	strReader := strings.NewReader("hello world")
	// 包裹为缓冲读取器,可以传入缓冲大小，默认为4096
	bufReader := bufio.NewReader(strReader)

	// 查看5个字节,不会读取出来，只是查看，后续读取时，还可以读取到这x个字节
	bytes,_ := bufReader.Peek(5)
	fmt.Println(string(bytes))

	// 查看缓冲读取器中缓冲的字符数， 11
	fmt.Println(bufReader.Buffered())

	// 读取到某个字符为止（包含这个字符）
	str,_ := bufReader.ReadString(' ')
	println(str,bufReader.Buffered())

	// 写入器
	wtiter := bufio.NewWriter(os.Stdout)
	// 将字符写入 写入器
	fmt.Fprint(wtiter,"hello ")
	fmt.Fprint(wtiter,"world")
	// 将字符 冲出
	wtiter.Flush()

}

/*计算文件行数*/
func countFileLine() {
	// 判断启动该命令时传入的参数个数
	if len(os.Args) <2{
		panic(errors.New("参数数目小于2，非法"))
	}

	// 取第二个参数作为文件名
	fileName := os.Args[1]

	file,err := os.Open(fileName)
	if err != nil{
		panic(errors.New("打开该文件异常:" + err.Error()))
	}
	defer file.Close()

	// 创建缓冲读取器
	reader := bufio.NewReader(file)

	var count int

	for {
		// 该读取方法，当 一行数据长度超过 缓冲大小（默认4096）时，第二个标识参数会为true，下一次调用，读取的还是同一行数据
		_,isPrefix, err := reader.ReadLine()
		// 读取结束
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(errors.New("读取该文件时异常:" + err.Error()))
		}

		if !isPrefix{
			count++
		}
	}

	fmt.Println("该文件行数：",count)
}

/*读取bmp文件头信息,依次读取*/
func readHeaderFromBmp() {
	file,_ := os.Open("C:\\Users\\97038\\Desktop\\1.bmp")
	defer file.Close()


	// 读取二进制
	// 第二个参数是字节序，windows和linux应该是小端的，mac是大端的。
	// 读取到headA的内存地址中
	var headA,headB byte
	binary.Read(file,binary.LittleEndian,&headA)
	binary.Read(file,binary.LittleEndian,&headB)

	// 会输出 BM
	fmt.Printf("%c%c \n",headA,headB)

	//文件大小
	var size uint32
	binary.Read(file,binary.LittleEndian,&size)
	fmt.Printf("%d \n",size)

	// 两个保留字节
	var reserveA,reserveB byte
	binary.Read(file,binary.LittleEndian,&reserveA)
	binary.Read(file,binary.LittleEndian,&reserveB)
	fmt.Println(reserveA,reserveB)

	// 偏移量，表示图片内容从哪开始
	var offBits uint32
	binary.Read(file,binary.LittleEndian,&offBits)
	fmt.Println(offBits)

	// 使用结构体继续读取
	headerInfo := new(BmpHeader)
	binary.Read(file,binary.LittleEndian,headerInfo)
	fmt.Println(headerInfo)

}



/*bmp头信息*/
type BmpHeader struct {
	Size  uint32
	Width int32
	Height int32
	Places uint16
	BitCount uint16
	Compression uint32
	SizeImage uint32
	XperlsPerMeter int32
	YperlsPerMeter int32
	ClsrUsed uint32
	ClrImportant uint32
}