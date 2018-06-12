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
	size := 1 * 10000
	arr1 := generateRandomArray(size, 0, size)
	arr2 := copyIntArray(arr1,size)
	testSort("选择排序1", arr1, len(arr1), selectSort)
	testSort("插入排序",arr2,len(arr2),insertSort)
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
		log.Println("该数组排序有误")
	}
	fmt.Println(sortName,",元素个数:",size,"执行时间:",endTime.Sub(startTime))
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
