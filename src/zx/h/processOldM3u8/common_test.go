package main

import (
	"testing"
	"io/ioutil"
)

func Test1(t *testing.T) {
	ioutil.WriteFile(`E:\新建文本文档.txt`,[]byte("2222"),0666)
}
