package test

import (
	"testing"
	"net"
	"jump/util"
	"fmt"
)

/*通用测试*/

func TestIp(t *testing.T) {
	conn,_ := net.Dial("tcp","106.14.7.29:80")
	ip,port := util.Ip2Bytes(conn.LocalAddr().String())
	fmt.Println(ip,"    ",port)
}
