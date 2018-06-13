package main

import (
	"net/http"
	"fmt"
	"sort"
)

/* hello world 接口*/
func main() {


	// 对切片进行排序
	data := []int{4,3,1,2,8,7,9}
	sort.Ints(data)
	for i,v := range data {
		fmt.Println(i,v)
	}


	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// 响应 语句，和传入的key为"key"的参数值
		fmt.Fprintln(writer,"hello world！ 参数:",request.FormValue("key"))
	})
	http.ListenAndServe(":8080",nil)
}