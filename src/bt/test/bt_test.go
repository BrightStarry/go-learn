package test

import (
	"testing"
	"encoding/binary"
	"fmt"
	"bytes"
)

/*测试类*/


/*测试大端序下的 长度 转 []byte的几种方法*/
func TestLengthToBytes(t *testing.T) {
	data :=[]byte{3,5,45,34,4,233,2,4,54,5,}

	// 方式1
	//length := uint32(len(data))
	//lenBytes := make([]byte,4)
	//binary.BigEndian.PutUint32(lenBytes, length)
	//fmt.Println("1:",lenBytes)
	//
	//// 方式2
	//length2 := int32(len(data))
	//lenBytes2 := bytes.NewBuffer(nil)
	//binary.Write(lenBytes2, binary.BigEndian, length2)
	//fmt.Println("2:",lenBytes2.Bytes())


}

// 测试
func TestStringToBytes(t *testing.T) {
	str := "xxx"
	// 这样可以自动将str转为[]byte进行追加操作
	data := append([]byte{1,2},str...)
	fmt.Println(data)
}

