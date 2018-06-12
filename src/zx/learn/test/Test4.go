package main

import "fmt"

func main() {

	/**
	将切片 b 的元素追加到切片 a 之后：a = append(a, b...)

	复制切片 a 的元素到新的切片 b 上：

	b = make([]T, len(a))
	copy(b, a)
	删除位于索引 i 的元素：a = append(a[:i], a[i+1:]...)

	切除切片 a 中从索引 i 至 j 位置的元素：a = append(a[:i], a[j:]...)

	为切片 a 扩展 j 个元素长度：a = append(a, make([]T, j)...)

	在索引 i 的位置插入元素 x：a = append(a[:i], append([]T{x}, a[i:]...)...)

	在索引 i 的位置插入长度为 j 的新切片：a = append(a[:i], append(make([]T, j), a[i:]...)...)

	在索引 i 的位置插入切片 b 的所有元素：a = append(a[:i], append(b, a[i:]...)...)

	取出位于切片 a 最末尾的元素 x：x, a = a[len(a)-1], a[:len(a)-1]

	将元素 x 追加到切片 a：a = append(a, x)
	 */


	// 切片（动态数组）,
	// 定义切片,此时 a == nil
	//var a []int
	// 创建切片,切片不需要指定容量，3为当前长度
	//a := make([]int,3)
	// 也可以指定容量
	//a := make([]int,3,10)
	// 初始化切片 cap=len=3
	//a := []int{1,2,3}

	// 数组
	var arr1 = []int{1, 2, 3}
	// 将arr1的引用作为 切片
	//a := arr1[:]
	// 将arr1的 下标从 0到2，取头不取尾（不包含2），作为索引
	a := arr1[0:2]
	// 从0到1
	//a := arr1[:1]
	// 从1，到末尾
	//a := arr1[1:]

	// 切片长度 2
	println(len(a))
	// 切片容量 3
	println(cap(a))
	fmt.Printf("%v", a)

	var numbers []int
	printSlice(numbers)

	/* 允许对空切片进行追加 */
	numbers = append(numbers, 0)
	printSlice(numbers)

	/* 向切片添加一个元素 */
	numbers = append(numbers, 1)
	printSlice(numbers)

	/* 同时添加多个元素 */
	numbers = append(numbers, 2, 3, 4)
	printSlice(numbers)

	/* 创建切片 numbers1 是之前切片的两倍容量*/
	numbers1 := make([]int, len(numbers), (cap(numbers))*2)

	/* 拷贝 numbers 的内容到 numbers1 */
	copy(numbers1, numbers)
	printSlice(numbers1)

	// range关键字，用于在for中迭代 array，slice，channel，map（返回key）等

	//这是我们使用range去求一个slice的和。使用数组跟这个很类似
	nums := []int{2, 3, 4}
	sum := 0
	for _, num := range nums {
		sum += num
	}
	fmt.Println("sum:", sum)
	//在数组上使用range将传入index和值两个变量。上面那个例子我们不需要使用该元素的序号，所以我们使用空白符"_"省略了。有时侯我们确实需要知道它的索引。
	for i, num := range nums {
		if num == 3 {
			fmt.Println("index:", i)
		}
	}
	//range也可以用在map的键值对上。
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}
	//range也可以用来枚举Unicode字符串。第一个参数是字符的索引，第二个是字符（Unicode的值）本身。
	for i, c := range "go" {
		fmt.Println(i, c)
	}
}

func printSlice(x []int) {
	fmt.Printf("len=%d cap=%d slice=%v\n", len(x), cap(x), x)
}
