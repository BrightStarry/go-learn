package test

import (
	"testing"
	"fmt"
	"time"
)

/**
	缓冲chan阻塞测试
 */

func TestBufferChan(t *testing.T) {
	m := map[string][]byte{}
	m["a"] = nil
	_,ok := m["a"]
	fmt.Println(ok)

	c := make(chan int,10)

	// 循环写入
	go func() {
		for {
			c <- 1
			fmt.Println("写入成功")
		}
	}()

	// 循环消费,每消费一个sleep若干时间
	for i:= 0; i < 5;i++{
		go func() {
			for v:= range c{
				fmt.Println(v)
				time.Sleep(5 * time.Second)
			}
		}()
	}
	time.Sleep(time.Hour)

}
