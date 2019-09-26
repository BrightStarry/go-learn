package test

import (
	"testing"
	"zx/ipProxyPool/config"
	"zx/ipProxyPool/util"
)

func TestParseToUrl(t *testing.T) {
	config.Init()
	util.GetByProxy("https://www.baidu.com",nil)
}
