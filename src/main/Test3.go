package main

import "fmt"

func main() {
	// 声明数组
	//var arr1 [10] string

	// 初始化数组
	var arr2  = [3]string{"a","b","c"}
	// 输出数组
	fmt.Printf("arr2:%v \n",arr2)

	// 根据{}中元素个数，自动设置数组长度
	var arr3 = [...]string{"c","b","a"}
	fmt.Printf("arr2:%v \n",arr3)

	// 二维数组
	var arr4 = [3][4]int{
		{0, 1, 2, 3} ,   /*  第一行索引为 0 */
		{4, 5, 6, 7} ,   /*  第二行索引为 1 */
		{8, 9, 10, 11},  /*  第三行索引为 2 */
	}
	fmt.Printf("arr4:%v\n",arr4)


	// 创建结构体对象
	user1 := user{"郑星",22}
	user1.username = "zx"
	user1.age = 1
	// 调用结构体的方法（结构体不能写在函数中）
	user1.toString()

	// 测试将结构体作为参数的 普通方法，结果依旧是值传递
	testUpdateUser(user1)

	//结构体指针
	toStringByPointer(&user1)

}

// 定义结构体
type user struct {
	username string
	age int
}

//该 method 属于 user 类型对象中的方法
func (u user) toString() {
	println(u.username,u.age)
}
// 结构体，使用指针
func toStringByPointer(u *user) {
	println("指针:",u.username,u.age)
}
func testUpdateUser(u user) user{
	u.username = "结构体默认方法"
	return u
}


// go函数默认为值传递，如下修改函数，可将其作为引用传递
// 可这样调用 ： swap2(&a, &b)
// x和y形参，接收的是x和y指针，然后将两个指针指向的内存地址值直接修改，即为引用传递

func swap2(x *int, y *int) {
	var temp int
	temp = *x    /* 保存 x 地址上的值 */
	*x = *y      /* 将 y 值赋给 x */
	*y = temp    /* 将 temp 值赋给 y */
}


