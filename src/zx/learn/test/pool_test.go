package test

import (
	"testing"
	"sync"
	"fmt"
)

/**
	对象池使用测试
 */
func TestPool(t *testing.T) {
	// 创建对象池
	pool := sync.Pool{New: func() interface{} {
		return "s"
	}}
	// 可使用put方法,主动将对象放入池中
	//pool.Put()
	// 如果调用get方法时,pool中没有对象了.则调用pool的New方法创建,如果没有指定,则返回nil
	fmt.Println(pool.Get())
	// 用完之后需要将对象放回去
	pool.Put("x")
	fmt.Println(pool.Get())


}
