package main

import (
	"net/http"
	"time"
	"log"
	"strings"
	"os"
	"encoding/json"
	"io/ioutil"
	"html/template"
)

/*实现静态文件服务器*/

// 多路复用器
// 定义一个key为string，值为 该函数的 map
var mux = make(map[string]func(http.ResponseWriter, *http.Request))

// 当前工作路径
var rootPath = "/"

func main() {
	server := http.Server{
		Addr:        ":8081",
		Handler:     &myHandler{},
		ReadTimeout: 5 * time.Second,
	}

	// 此处配置路由，可添加自定义函数，处理对应路由
	mux["/"] = IndexView

	//var err error
	//rootPath,err = os.Getwd()
	//if err != nil {
	//	panic(errors.New("当前路径获取失败:" + err.Error()))
	//}
	//log-monitor.Println("当前路径:",rootPath)

	log.Println("服务已启动")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("服务异常:", err)
	}

}

/*处理类，该类实现了http.Handler接口*/
type myHandler struct{}

/*处理类的处理方法*/
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("请求方法:", r.Method, ",请求url：", r.URL.String())
	// 解析，默认不解析，如果不解析，r.Form将拿不到数据
	r.ParseForm()
	log.Println("请求报文:", r)
	log.Println("请求参数:", r.Form)

	// 如果map中有key为该路由，也就是默认配置的路由
	if request, ok := mux[r.URL.String()]; ok {
		// 设置该元素的值
		request(w, r)
	} else {
		// 文件过滤器
		fileFilter(w,r)
	}
}

/*返回的JsonBean*/
type ResultDTO struct {
	// 后面的 "`json:"code"`" 是在json转换时，将该属性取别名
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

/*创建ResultDTO*/
func NewResultDTO(code int, message string, data interface{}) *ResultDTO {
	return  &ResultDTO{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

/*文件过滤*/
func fileFilter(w http.ResponseWriter, r *http.Request) {
	// 访问的路径（该路径是不包含主机名啥的，具体可以查看url.URL的源码注解）
	path := r.URL.Path
	// 判断是否有.
	if strings.Contains(path, ".") {
		// 截取出最后一个. 之后的字符串
		requestType := path[strings.LastIndex(path,"."):]
		switch requestType {
		case ".css":
			w.Header().Set("content-type","text/css; charset=utf-8")
		case ".js":
			w.Header().Set("content-type","text/javascript; charset=utf-8")
		default:
		}
	}

	// 读取要显示的文件
	srcFile,err := os.Open(rootPath + path)
	defer srcFile.Close()
	// 读取失败
	if err != nil {
		log.Println("读取文件失败：",err)
		// 返回json头
		w.Header().Set("content-type","text/json; charset=utf-8")
		result := NewResultDTO(404,"",nil)
		bytes,_ := json.Marshal(result)
		w.Write(bytes)
		log.Println("返回数据:",string(bytes))
		return
	}

	// 读取成功
	result,_ := ioutil.ReadAll(srcFile)
	w.Write(result)
}

/*默认首页访问方法*/
func IndexView(w http.ResponseWriter,r *http.Request) {
	t,err := template.ParseFiles("index.html")
	// 没有异常，直接返回
	if err == nil {
		t.Execute(w,nil)
		return
	}
	// 有异常,返回默认首页
	log.Println("未找到index.html文件，返回默认首页")
	w.Header().Set("content-type","text/html; charset=utf-8")
	w.Write([]byte(indexTemplate))
}

//首页模板
var indexTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>It Work</title>
</head>
<body>
<h1>:)</h1>
<h3>It Work</h3>
</body>
</html>
`
