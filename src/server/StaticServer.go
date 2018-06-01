package server

/*实现静态文件服务器*/

// 多路复用器
// 定义一个key为string，值为 该函数的 map
//var mux map[string]func(http.ResponseWriter, *http.Request)
//
//func main() {
//	server := http.Server{
//		Addr: "8081",
//		Handler: &myHandler{},
//		ReadTimeout: 5*time.Second,
//	}
//}
//
///*处理类*/
//type myHandler struct {}
