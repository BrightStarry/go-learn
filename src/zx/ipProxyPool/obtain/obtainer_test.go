package obtain

import (
	"testing"
	"zx/ipProxyPool/config"
	"zx/ipProxyPool/util"
	"fmt"
)


/**
	 测试
 */

func TestXiciAnonymity_InitObtain(t *testing.T) {
	config.Init()
	Init()
	//fmt.Println(WebObtainers[0])
	//WebObtainers[0].IncrementObtain()
	//
	//go func() {
	//	for v:= range config.ObtainerOutChan{
	//		fmt.Println(v)
	//	}
	//}()
	//time.Sleep(time.Hour)
	document := util.GetOfDocument("https://www.baidu.com")
	fmt.Println(document.Html())
}
