package main

import "fmt"

func main() {
	// if switch 语法和java一样，只不过条件表达式，需要去掉 ()

	//if(ptr != nil)     /* ptr 不是空指针 */
	//if(ptr == nil)    /* ptr 是空指针 */

	// 可以用switch判断对象类型
	var o interface{}
	switch i := o.(type) {
	case nil:
		fmt.Printf(" o 的类型 :%T",i)
	case int:
		fmt.Printf("o 是 int 型")
	case float64:
		fmt.Printf("o 是 float64 型")
	case func(int) float64:
		fmt.Printf("o 是 func(int) 型")
	case bool, string:
		fmt.Printf("o 是 bool 或 string 型" )
	default:
		fmt.Printf("未知型")
	}

	// select随机执行一个可运行的case。如果没有case可运行，它将阻塞，直到有case可运行。一个默认的子句应该总是可运行的。
	// 监听io的channel操作
	/**
	每个case都必须是一个通信
	所有channel表达式都会被求值
	所有被发送的表达式都会被求值
	如果任意某个通信可以进行，它就执行；其他被忽略。
	如果有多个case都可以运行，Select会随机公平地选出一个执行。其他不会执行。
	否则：
	如果有default子句，则执行该语句。
	如果没有default字句，select将阻塞，直到某个通信可以运行；Go不会重新对channel或值进行求值。
	**/

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	select {
	case e1 := <-ch1:
		//如果ch1通道成功读取数据，则执行该case处理语句
		fmt.Printf("1th case is selected. e1=%v", e1)
	case e2 := <-ch2:
		//如果ch2通道成功读取数据，则执行该case处理语句
		fmt.Printf("2th case is selected. e2=%v", e2)
	default:
		//如果上面case都没有成功，则进入default处理流程
		fmt.Println("default!.")
	}


	// for循环
	//for init; condition; post { } 普通的for
	//for condition { }  相当于while
	//for { } 相当于for(;;) 无限循环

	// 函数定义
	/*
	func function_name( [parameter list] ) [return_types] {
		函数体
	}
	*/

	//多个返回值函数
	x, y := swap("zzz", "xxx")
	println(x,y)

	// 闭包，将函数作为对象（有权访问另一个函数作用域内变量的函数都是闭包）
	addFunc := add(1,2)
	i1,i2 := addFunc()
	println(i1,i2)


	//全局变量与局部变量名称可以相同，但是函数内的局部变量会被优先考虑

}

// 闭包使用方法
func add(x1, x2 int) func()(int,int)  {
	i := 0
	return func() (int,int){
		i++
		return i,x1+x2
	}
}

/*多个返回值函数*/
func swap(x, y string) (string, string) {
	return y, x
}

/* 函数返回两个数的最大值 */
func max(num1, num2 int) int {
	/* 声明局部变量 */
	var result int

	if (num1 > num2) {
		result = num1
	} else {
		result = num2
	}
	return result
}
