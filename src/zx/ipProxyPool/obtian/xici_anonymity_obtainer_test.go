package obtian

import (
	"testing"
	"zx/ipProxyPool/config"
	"fmt"
)


/**
	 测试
 */

func TestXiciAnonymity_InitObtain(t *testing.T) {
	config.Init()
	fmt.Println(config.DefaultClient)
	o := &XiciAnonymity{BaseObtainer{WebInfo:&config.WebInfos[0]}}

	proxyIps := o.InitObtain()


	fmt.Println(o.WebInfo)
	fmt.Println(proxyIps)
}
