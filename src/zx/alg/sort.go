package alg

import (
	"math/rand"
	"time"
	"fmt"
	"log"
)

/*基础排序*/

/*
	生成随机数组
	长度为size,范围为[rangeL,rangeR]
*/
func generateRandomArray(size , rangeL , rangeR int) *[]int{
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
	return &arr
}

/**
	测试算法性能
 */
func testSort(sortName string, arr *[]int,size int,sort func(*[]int, int),) {
	if len(*arr) != size {
		log.Fatalln("数组长度和传入的元素个数不等")
	}
	startTime := time.Now()
	sort(arr,size)
	endTime := time.Now()
	if !isSorted(arr,size) {
		log.Println("该数组排序有误")
	}
	fmt.Println(sortName,",元素个数:",size,"执行时间:",endTime.Sub(startTime))
}

/**
	测试数组是否有序
 */
func isSorted(arr *[]int,size int) bool {
	for i:=0;i<  size -1; i++{
		if (*arr)[i] > (*arr)[i+1] {
			return false
		}
	}
	return true
}

/**
	选择排序  O(n^2)
	从数组中依次找出最小的元素,放到前面(第一个,第二个第三个...)
 */
func selectSort(arr *[]int,size int) {
	// 遍历数组, 第一次遍历[0:],第二次[1:]
	for i:= 0; i< size; i++{
		// 寻找[i,n)区间内的最小值
		// 子循环当前的最小值索引
		minIndex := i
		for j := i+1; j < size; j++{
			// 如果该数小于目前的最小数,则记录该数索引
			if (*arr)[j] < (*arr)[minIndex]{
				minIndex = j
			}
		}
		// 交换:将该最小值放到前面
		(*arr)[i],(*arr)[minIndex]  = (*arr)[minIndex],(*arr)[i]
	}
}
