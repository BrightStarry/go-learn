package obtain

import (
	"testing"
	"zx/ipProxyPool/config"
)


/**
	 测试
 */

func TestXiciAnonymity_InitObtain(t *testing.T) {

	config.Init()
	WebObtainers[6].IncrementObtain()
}
