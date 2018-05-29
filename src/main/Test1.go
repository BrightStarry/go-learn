package main

import (
	"fmt"
	"unsafe"
)

// 标识符（包括常量、变量、类型、函数名、结构字段等等）以一个大写字母开头，则它可以被外部代码导入，小写则为private

// main函数作为程序的入口，只有package名为main的包中的go文件可以包含main函数，一个可执行文件，有且仅有一个main包
// 可以通过import关键字导入其他非main包

// 为fmt起别名为fmt2
//import fmt2 "fmt"

// 调用的时候只需要Println()，而不需要fmt.Println()
//import . "fmt"
func main() {
	// 常量(此处省略了类型)
	const PI = 3.14
	// 变量
	var name = "ZhengXing"
	// 打印
	println(name)
	// 一般类型声明
	type newType int
	// 结构类型声明
	type user struct{}
	// 接口声明
	type userInterface interface {}

	// 基本类型
	var a bool = true
	// 默认8位(-128-127)
	var b int = 100
	// 无符号8位(0-255)
	var b1 uint = 1
	var b2 int32 = 100000
	var c float32 = 10.32
	// 实数和虚数
	var d complex64 = 10
	// 字节
	var e byte = 254
	// 1.9版本对于数字类型，无需定义int及float32、float64，系统会自动识别。
	var f = 10.32
	// 省略var(声明新变量的缩写语法，变量名不能为 定义过的，否则会报错)
	g := true
	//类型相同多个变量, 非全局变量
	//var vname1, vname2, vname3 bool
	//和python很像,不需要显示声明类型，自动推断
	var vname1, vname2, vname3 = true, "xx", 30.2
	// 或 vname1, vname2, vname3 := v1, v2, v3 (这种不带声明格式的只能在函数体中出现)

	// 这种因式分解关键字的写法一般用于声明全局变量
	var (
		c1 int = 3
		c2 string = "dfd"
	)

	// 类型转换
	var sum int = 17
	var count int = 5
	var mean float32

	mean = float32(sum)/float32(count)
	fmt.Printf("mean 的值为: %f\n",mean)

	// 多个变量可以再同一行赋值
	//var a, b int
	//var c string
	//a, b, c = 5, 7, "abc"


	// 上面这些基本类型都属于值类型
	// 可用如下方法获取变量的内存地址
	fmt.Println(&a)
	// 以下，是通过取得的地址，获取值，也就是一个指针，指向a的内存地址
	aPointer := &a
	fmt.Println(*aPointer)//此处输出a的值

	// 交换两个变量的值,必须同类型
	b,c1 = c1,b
	// _ 空白标识符，可以将值赋给一个被抛弃的变量
	_,c1 = 5,6

	// 并行赋值也可用于获取一个函数返回的多个返回值
	//	val, err = Func1(var1)

	fmt.Println(a,b,b1,b2,c,d,e,f,g,vname1,vname2,vname3,c1,c2)

	// 常量可用于枚举
	const (
		Unknown = 0
		Female = 1
		Male = 2
	)

	//常量可以用len(), cap(), unsafe.Sizeof()函数计算表达式的值。常量表达式中，函数必须是内置函数，否则编译不过：
	const (
		aaa = "abc"
		bbb = len(aaa)
		ccc = unsafe.Sizeof(aaa)
	)
	println(aaa,bbb,ccc)
	// ccc值为16，字符串类型在 go 里是个结构, 包含指向底层数组的指针和长度,这两部分每部分都是 8 个字节，所以字符串类型大小为 16 个字节。


	//iota，特殊常量，可以认为是一个可以被编译器修改的常量。
	//iota在const关键字出现时将被重置为0(const内部的第一行之前)，const中每新增一行常量声明将使iota计数一次(iota可理解为const语句块中的行索引)。
	// 我们可以使用下划线跳过不想要的值。

	type AudioOutput int
	const (
		OutMute AudioOutput = iota // 0
		OutMono                    // 1
		OutStereo                  // 2
		_
		_
		OutSurround                // 5
	)

	type Allergen int

	const (
		IgEggs Allergen = 1 << iota // 1 << 0 which is 00000001
		IgChocolate                         // 1 << 1 which is 00000010
		IgNuts                              // 1 << 2 which is 00000100
		IgStrawberries                      // 1 << 3 which is 00001000
		IgShellfish                         // 1 << 4 which is 00010000
	)

	type ByteSize float64

	const (
		_           = iota                   // ignore first value by assigning to blank identifier
		KB ByteSize = 1 << (10 * iota) // 1 << (10*1)
		MB                                   // 1 << (10*2)
		GB                                   // 1 << (10*3)
		TB                                   // 1 << (10*4)
		PB                                   // 1 << (10*5)
		EB                                   // 1 << (10*6)
		ZB                                   // 1 << (10*7)
		YB                                   // 1 << (10*8)
	)
}
