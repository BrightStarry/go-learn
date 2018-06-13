package test

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"strings"
)

/**
	http测试
 */
func TestErrWrapper(t *testing.T) {
	tests := []struct{
		h httpHandler
		code int
		message string
	}{
		{handler,500,"Internal Server Error"},
	}

	for _,tt := range tests{
		// 处理方法
		f := errWrapper(tt.h)
		// 创建测试响应和请求
		response := httptest.NewRecorder()
		request := httptest.NewRequest(
			http.MethodGet,
			//随意写的url
			"http://127.0.0.1",nil)
		// 执行处理方法
		f(response,request)

		// 从response中读取
		b,_ := ioutil.ReadAll(response.Body)
		// 去除换行符
		body := strings.Trim(string(b),"\n")
		if response.Code != tt.code ||
			 body != tt.message{
			t.Errorf("预期(%d, %s); 实际(%d, %s)",tt.code,tt.message,response.Code,body)
		}
	}

	/**
		也可以如下直接启动一个测试服务器,然后发起请求
		httptest.NewServer()
	 */

}
