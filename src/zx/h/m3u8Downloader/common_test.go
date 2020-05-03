package main

import (
	"testing"

	"encoding/base64"
	"io/ioutil"
	"strings"
	"fmt"
)

/**
	后台运行aria2
E:\test\aria2c.exe --conf-path=E:\test\aria2.conf --log=E:\test\task1.log --input-file=E:\test\aria2.session --save-session=E:\test\aria2.session --rpc-listen-port=1088 --daemon=false
E:\test\aria2c.exe --conf-path=E:\test\aria2.conf --log=E:\test\task1.log  --rpc-listen-port=1088 --daemon=false -P -Z

 */
func Test1(t *testing.T) {

	//f, err := os.Open("C:\\h\\av2\\m3u8\\all\\186p00035.m3u8")
	//if err != nil {
	//	panic(err)
	//}
	//p, _, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	//if err != nil {
	//	panic(err)
	//}
	//data := p.(*m3u8.MediaPlaylist)
	//
	//
	///**
	//获取ts 的uri列表
	// */
	//uris := make([]string,0)
	//for _,x := range data.Segments {
	//	if x == nil {
	//		continue
	//	}
	//	uris = append(uris, x.URI)
	//}
	//
	//
	//option := make(map[string]string)
	//option["dir"] = "E:\\test\\2"

	//for _, v := range uris {
	//	g, err := rpc1.AddURI(v,option)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//}

}

func Test2(t *testing.T) {
	split := strings.Split("ktds-814~1~h_094ktds00814", "~")

	fmt.Println(split)

}

func TestBase64(t *testing.T) {
	base64Encoder :=base64.StdEncoding
	bytes, err := base64Encoder.DecodeString("l25hZADrbvoO1EFLwauJuQ==")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(`C:\h\ipx00367.key`, bytes, 0666)
	if err != nil {
		panic(err)
	}
}
