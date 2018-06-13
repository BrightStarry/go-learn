package new

import (
	"testing"
	"net/http"
	"fmt"
	"time"
	"io/ioutil"
	"strings"
	"net/http/cookiejar"
	"net/url"
)

/*
	用go发起htto请求
*/

func TestHttpClient(t *testing.T) {
	u,_ := url.Parse("http://127.0.0.1")

	// 构造请求客户端
	client := &http.Client{}
	// 设置超时时间
	client.Timeout =  5 *  time.Second


	// 设置cookie,其中New()本应传入的参数是指cookie的作用域范围
	// 该cookie会自动保存网站写入的cookie
	jar,_ := cookiejar.New(nil)
	client.Jar = jar
	// 设置cookie，需要设置cookie的作用域范围
	client.Jar.SetCookies(u,[]*http.Cookie{{Name: "a", Value: "b"},{Name: "a", Value: "b"}})


	// 构造请求
	request,err := http.NewRequest(http.MethodPost,"http://127.0.0.1/a",strings.NewReader("name=aaa"))
	if err != nil {
		panic(err)
	}

	// 设置请求头
	request.Header.Set("x","x")

	// 进行请求
	response,err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 或者直接(实际上内部用一个默认的client，进行了NewRequest，Do等同上的操作)
	//response,err := http.Get("https://www.jianshu.com/p/757d133021de")

	// 进行post form提交，否则构造post请求时需要追加content-type
	//client.PostForm()


	// 输出返回到控制台
	//if _,err=io.Copy(os.Stdout,response.Body);err != nil {
	//	panic(err)
	//}

	// 直接读取为字节
	body,_ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	// 用goquery将其解析为doc
	//doc,err := goquery.NewDocumentFromReader(response.Body)
	//if err != nil {
	//	panic(err)
	//}
	// 用selector选择dom
	//ele := doc.Find("body > div.note > div.post > div.article > div.show-content > div > h1:nth-child(16)")
	//fmt.Println(ele.Text())

	fmt.Println(client.Jar.Cookies(u))


}
