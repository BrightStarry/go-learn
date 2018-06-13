package alg

import (
	"testing"
	"fmt"
	"time"
	"math/rand"
	"log"
)

/*基础排序测试*/

/**
	测试
 */
func TestA(t *testing.T) {
	size := 10 * 10000
	arr1 := generateRandomArray(size, 0, size)
	arr2 := copyIntArray(arr1,size)
	arr3 := copyIntArray(arr1,size)
	arr4 := copyIntArray(arr1,size)
	testSort("选择排序", arr1, size, selectSort)
	testSort("插入排序",arr2,size,insertSort)
	testSort("冒泡排序",arr3,size,bubbleSort)
	testSort("希尔排序",arr4,size,shellSort)
}

/**
	性能测试
	会执行若干次,显示平均执行时间
	默认会执行1秒左右的测试
 */
func BenchmarkTest(b *testing.B) {
	// 多显示几个参数
	b.ReportAllocs()
	// ...假设进行若干操作
	// 重置时间统计
	b.ResetTimer()
	for i := 0; i < b.N; i++{
		size := 1 * 10000
		arr1 := generateRandomArray(size, 0, size)
		arr2 := copyIntArray(arr1,size)
		arr3 := copyIntArray(arr1,size)
		arr4 := copyIntArray(arr1,size)
		testSort("选择排序", arr1, size, selectSort)
		testSort("插入排序",arr2,size,insertSort)
		testSort("冒泡排序",arr3,size,bubbleSort)
		testSort("希尔排序",arr4,size,shellSort)
	}
}

/*
	生成随机数组
	长度为size,范围为[rangeL,rangeR]
*/
func generateRandomArray(size , rangeL , rangeR int) []int{
	if rangeL > rangeR || rangeL < 0 {
		panic("随机数生成参数有误")
	}
	arr := make([]int,size,size)
	// 随机数种子
	rand.Seed(time.Now().Unix())
	for i:=0;i<size;i++{
		// intn方法,随机生成[0,参数]的数字
		arr[i] = rand.Intn(rangeR - rangeL +1) + rangeL
	}
	return arr
}

/**
	生成近乎有序的随机数
 */
 func generateNearlyOrderedArray(size,swapTimes int) []int{
 	// 先生成一个完全有序的数组, go中,使用make构建的分片,元素会被附上默认值,int,也就是全为0的分片
 	arr := make([]int,size,size)
 	for i:=range arr{
 		arr[i] = i
	}
	 // 随机数种子
	 rand.Seed(time.Now().Unix())
 	// 然后随机选取元素进行交换,进行swapTimes次
 	for i:=0;i< swapTimes; i++{
		posX :=rand.Intn(size)
		posY := rand.Intn(size)
		arr[posX],arr[posY] = arr[posY],arr[posX]
	}
	return arr
 }

/**
	测试算法性能
 */
func testSort(sortName string, arr []int,size int,sort func([]int, int) []int) {
	if len(arr) != size {
		log.Fatalln("数组长度和传入的元素个数不等")
	}
	startTime := time.Now()
	arr = sort(arr,size)
	endTime := time.Now()
	if !isSorted(arr,size) {
		log.Println(sortName + "排序有误")
	}
	fmt.Println(sortName,",元素个数:",size,"执行时间:",endTime.Sub(startTime))
	// 打印结果
	//fmt.Println(sortName,arr)
}

/**
	测试数组是否有序
 */
func isSorted(arr []int,size int) bool {
	for i:=0;i<  size -1; i++{
		if arr[i] > arr[i+1] {
			return false
		}
	}
	return true
}

/**
	拷贝int数组
 */
func copyIntArray(arr []int , size int) []int{
	if size > len(arr) {
		log.Panicln("要拷贝的长度越界")
	}
	result := make([]int,size)
	copy(result,arr[:size])
	return result
}
