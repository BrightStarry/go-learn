package test

import (
	"net/http"
	"log"
	"testing"
	"errors"
	"fmt"
	"os"
)

/*
	统一异常处理测试
*/
func TestUnifyErrorProcess(t *testing.T) {
	http.HandleFunc("/list/", errWrapper(handler))
	if 	err := http.ListenAndServe(":8080",nil); err != nil {
		log.Panicln()
	}
}

type httpHandler func(w http.ResponseWriter, r *http.Request) error

/**
	处理函数
 */
func handler(w http.ResponseWriter, r *http.Request) error {
	return errors.New("xxxx")
	fmt.Fprintln(w,"22")
	return nil
}

/**
	统一异常处理
 */
func errWrapper(handler httpHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		// 进行异常捕获
		defer func() {
			if recover := recover(); r!= nil {
				log.Println("异常:",recover)
				http.Error(
					w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		 err:= handler(w,r)
		 if err != nil {

		 	// 如果是自定义异常,直接返回
		 	if userErr,ok := err.(userError);ok{
		 		http.Error(w,userErr.Message(),http.StatusBadRequest)
		 		return
			}

		 	code := http.StatusOK
		 	// 根据异常类型,作不同处理
			switch  {
			case os.IsNotExist(err):
				code = http.StatusNotFound

			default:
				code = http.StatusInternalServerError
			}
			 http.Error(w, http.StatusText(code), code)
		}
	}
}

/**
	自定义异常接口
 */
 type userError interface {
 	// 整合接口,让userError接口必须实现error接口的Error()方法
 	error
 	Message() string
 }

 /**
 	自定义异常类型
  */
type customError string
func (this customError) Error() string{
	return this.Message()
}
func (this customError) Message() string {
	return string(this)
}