package alg

import (
	"testing"
)

/*基础排序测试*/

/**
	测试
 */
func TestA(t *testing.T) {
	arr := generateRandomArray(100000,1,10000)
	testSort("选择排序",arr,100000,selectSort)
}





