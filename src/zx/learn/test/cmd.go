package main

import (
	"os/exec"
	"fmt"
)



func main() {
	cmd := exec.Command("ping", "www.baidu.com")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("stdout=[%s]\n", string(out))
}
