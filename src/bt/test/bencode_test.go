package test

import (
	"bt/util"
	"fmt"
	"testing"
)

/*编解码测试*/

func TestA(t *testing.T) {
	data := util.IntEncode(2)
	fmt.Println(data)
}

/*测试*/
func TestEncode(t *testing.T) {
	var testMap = make(map[string]interface{})
	testMap["a"] = 343
	testMap["b"] = "zx"
	testList := make([]interface{},3)
	testList[0] = "3"
	testList[1] = 34
	testList[2] = 133


	testMap["c"] = testList

	result := util.Encode(testMap)
	fmt.Println(result)


	result2,err := util.Decode([]byte(result))
	if err != nil {
		fmt.Println("异常:",err)
	}
	fmt.Println(result2)

}
