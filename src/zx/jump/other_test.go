package jump

import (
	"testing"
	"net"
	"fmt"
	"log"
	"zx/jump/util"
)

/*通用测试*/

func TestIp(t *testing.T) {
	conn,_ := net.Dial("tcp","106.14.7.29:80")
	ip,port := util.Ip2Bytes(conn.LocalAddr().String())
	fmt.Println(ip,"    ",port)

	data := [4]byte{12,12,12,12}
	s := net.IPv4(data[0],data[1],data[2],data[3]).String()
	log.Println(s)
}

func TestA(t *testing.T) {
	arr := [...]int{1,2,3,4,5}
	s := arr[0:3]
	a(s)

	fmt.Println(arr)
	fmt.Println(s)


}

func a(s []int) {
	s[0] = 100
}


