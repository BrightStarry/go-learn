package test

import (
	"testing"
	"zx/ipProxyPool/config"
	"fmt"
	"zx/ipProxyPool/obtain"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"strings"
	"log"
)

func TestIp66(t *testing.T) {
	// 初始化系统参数
	config.InitSystemConfig()
	fmt.Println("参数:",config.Config)

	// 配置初始化
	config.Init()
	obtain.WebObtainers[2].IncrementObtain()

}

func TestOtto(t *testing.T) {
	filePath := "test.js"
	//先读入文件内容
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}


	vm := otto.New()
	_,err = vm.Run(string(bytes))
	if err!=nil {
		panic(err.Error())
	}
	v, err := vm.Get("a")
	if  err != nil {
		panic(err.Error())
	}

	fmt.Println(v)


	filePath2 := "test2.js"
	//先读入文件内容
	bytes2, err := ioutil.ReadFile(filePath2)
	if err != nil {
		panic(err)
	}
	script2 := string(bytes2)
	flag1 := "Path=/;'"
	i1 := strings.Index(script2,"='__jsl")
	i2 := strings.Index(script2,flag1) + len(flag1)
	log.Println("%d , %d",i1,i2)
	script2 = "result" + script2[:]

	// 执行第二个js
	_,err = vm.Run(script2)
	result2, err := vm.Get("result")
	if  err != nil {
		panic(err.Error())
	}
	fmt.Println(result2)
}

func testTemp(t *testing.T) {

}
