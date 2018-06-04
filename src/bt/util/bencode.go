package util

import (
	"strings"
	"strconv"
	"errors"
	"bytes"
	"unicode/utf8"
	"unicode"
)

/*bencode编解码*/

/*
	解码相关
*/



/*从指定偏移量后开始查找目标值*/
func find(data []byte,offset int,target int32) (index int) {
	index = bytes.IndexRune(data[offset:],target)
	if index != -1 {
		return offset + index
	}
	return index
}

/*string解码*/
func StringDecode(data []byte,offset int)(result interface{},index int, err error) {
	if offset >= len(data) || data[offset] < '0' || data[offset] > '9' {
		err = errors.New("string解码:无效")
		return
	}
	i := find(data,offset,':')
	if i == -1{
		err = errors.New("string解码:找不到':'")
		return
	}

	length,err := strconv.Atoi(string(data[offset:i]))
	if err != nil {
		return
	}
	if length < 0 {
		err = errors.New("string解码:	无效长度")
		return
	}

	index = i +  length + 1
	if index > len(data) || index < i+1{
		err = errors.New("string解码:	下标越界")
		return
	}
	result = string(data[i+1:index])
	return
}

/*int解码*/
func IntDecode(data []byte,offset int) (result interface{},index int,err error){
	if offset >= len(data) || data[offset] != 'i' {
		err = errors.New("int解码:无效")
		return
	}
	index = find(data,offset+1,'e')
	if index == -1 {
		err = errors.New("int解码:找不到':'")
		return
	}
	result,err = strconv.Atoi(string(data[offset+1:index]))
	if err != nil {
		return
	}
	index++
	return
}

/*任意类型解码*/
func anyDecode(data []byte,offset int) (result interface{},index int,err error){
	// 该语句声明为全局，会循环引用
	var decodeFunc = []func([]byte, int) (interface{}, int, error){
		StringDecode, IntDecode, ListDecode, DictDecode,
	}
	for _,f := range decodeFunc{
		result,index,err = f(data, offset)
		if err == nil {
			return
		}
	}
	err = errors.New("解码：无效")
	return
}

/*list解码*/
func ListDecode(data []byte,offset int) (result interface{},index int,err error) {
	if offset > len(data) || data[offset] != 'l' {
		err = errors.New("list解码:无效")
		return
	}

	var item interface{}
	r := make([]interface{},0,8)

	index = offset + 1
	for index < len(data){
		char,_ := utf8.DecodeRune(data[index:])
		if char == 'e' {
			break
		}
		item,index,err = anyDecode(data,index)
		if err != nil {
			return
		}
		r = append(r,item)
	}

	if index == len(data) {
		err = errors.New("list解码:没有找到结束符")
		return
	}
	index++
	result = r
	return
}

/*字典解码*/
func DictDecode(data []byte, offset int) (result interface{},index int,err error) {
	if offset > len(data) || data[offset] != 'd' {
		err = errors.New("dict解码:无效")
		return
	}

	var item,key interface{}
	r := make(map[string]interface{})

	index = offset + 1
	for index < len(data){
		char,_ := utf8.DecodeRune(data[index:])
		if char == 'e' {
			break
		}

		if !unicode.IsDigit(char) {
			err = errors.New("dict解码:无效")
			return
		}

		key,index,err = StringDecode(data,index)
		if err != nil {
			return
		}

		if index > len(data) {
			err = errors.New("dict解码:越界")
			return
		}

		item,index,err = anyDecode(data,index)
		if err != nil {
			return
		}

		r[key.(string)] = item
	}
	if index == len(data) {
		err = errors.New("dict解码:没有找到结束符")
		return
	}
	index++

	result = r
	return
}

/*解码*/
func Decode(data []byte) (result interface{},err error) {
	result ,_,err = anyDecode(data,0)
	return
}

/*
	编码相关
*/

/*string编码*/
func StringEncode(data string) string {
	return strings.Join([]string{strconv.Itoa(len(data)),data},":")
}

/*int编码*/
func IntEncode(data int) string {
	return "i" + strconv.Itoa(data) + "e"
}

/*任意类型编码*/
func anyEncode(data interface{})(item string) {
	switch value := data.(type){
	case string:
		item = StringEncode(value)
	case int:
		item = IntEncode(value)
	case []interface{}:
		item = ListEncode(value)
	case map[string]interface{}:
		item = DictEncode(value)
	default:
		panic("任意编码:无效类型")
	}
	return
}

/*list编码*/
func ListEncode(data []interface{}) string{
	result := make([]string,len(data))
	for i,item := range data {
		result[i] = anyEncode(item)
	}
	return strings.Join([]string{"l",strings.Join(result,""),"e"},"")
}

/*字典编码*/
func DictEncode(data map[string]interface{}) string {
	result,i := make([]string,len(data)),0
	for k,v := range data {
		result[i] = StringEncode(k) + anyEncode(v)
		i++
	}
	return strings.Join([]string{"d",strings.Join(result,""),"e"},"")
}


func Encode(data interface{})string {
	switch value := data.(type){
	case string:
		return StringEncode(value)
	case int:
		return IntEncode(value)
	case []interface{}:
		return ListEncode(value)
	case map[string]interface{}:
		return DictEncode(value)
	default:
		panic("编码:无效类型")
	}
}



