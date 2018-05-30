package test

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// 一个简单的 http服务器例子
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		s := fmt.Sprintf("当前的时间是: %s", time.Now().String())
		// 将format的结果用response写入
		fmt.Fprintf(response, "%v\n", s)
		log.Printf("%v\n", s)
	})

	// 监听服务器
	// 此处执行顺序应该是，先执行ListenAndServe函数并阻塞，如果出现异常，就会继续执行并返回异常，并执行下一个判断
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("服务器异常:", err)
	}
}
