package bt2

import (
	"testing"
	"io/ioutil"
	"fmt"
)

/**
测试解码
 */
 func Test(t *testing.T) {
	 bytes, err := ioutil.ReadFile("C:\\Users\\97038\\Downloads\\[FHD]PGD-957.torrent")
	 if err != nil {
		panic(err)
	 }


 	result,_ := Decode(bytes)

 	r := result.( map[string]interface{})
	info:= r["info"].( map[string]interface{})
	files := info["files"].( []interface{})
	for _,item:= range files{
		file := item.( map[string]interface{})
		fmt.Println(file)
	}


 	fmt.Println(files)


 }
