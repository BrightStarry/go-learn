package jump

import (
	"testing"
	"net"
	"fmt"
	"log"
)

/*通用测试*/

func TestIp(t *testing.T) {
	conn,_ := net.Dial("tcp","106.14.7.29:80")
	ip,port := Ip2Bytes(conn.LocalAddr().String())
	fmt.Println(ip,"    ",port)

	data := [4]byte{12,12,12,12}
	s := net.IPv4(data[0],data[1],data[2],data[3]).String()
	log.Println(s)
}
