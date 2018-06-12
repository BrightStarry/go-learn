package new

import (
	"net/http"
	"log"
	"testing"
	"errors"
	"fmt"
)

/*
	统一异常处理测试
*/
func TestUnifyErrorProcess(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wrapError(w,r,handler)
	})
	if 	err := http.ListenAndServe(":8080",nil); err != nil {
		log.Panicln()
	}
}

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
func wrapError(w http.ResponseWriter, r *http.Request,handler func (http.ResponseWriter,  *http.Request) error) {
	if err := handler(w, r);err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
}
