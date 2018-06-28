package taken

import (
	"os/exec"
	"testing"
)

// 通用测试类

/**
	测试自动配置pac脚本
 */
func TestAutoConfigPAC(t *testing.T) {
	cmd := exec.Command("reg","add","HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
		"/v","AutoConfigURL", "/t","REG_SZ", "/d","xxx", "/f")
	cmd.Run()

	cmd2 := exec.Command("ipconfig","/flushdns")
	cmd2.Run()
}
