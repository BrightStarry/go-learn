package util

import (
	"testing"
	"zx/ipProxyPool/config"
)

func TestParseToUrl(t *testing.T) {
	config.Init()
	GetByProxy("https://www.baidu.com",nil)
}
