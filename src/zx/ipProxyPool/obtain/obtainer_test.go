package obtain

import (
	"testing"
	"zx/ipProxyPool/config"
)


/**
	 测试
 */

func TestXiciAnonymity_InitObtain(t *testing.T) {
	config.InitSystemConfig()
	config.Init()
	WebObtainers[7].IncrementObtain()
}
