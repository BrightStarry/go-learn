package new

import (
	"testing"
	"fmt"
	"unicode/utf8"
)

func TestRune(T *testing.T) {
	s := "a人间正道是沧桑"

	// 字符串转[]byte,遍历
	for _,b := range []byte(s) {
		// 每个英文1个字节,每个中文3个字节 UTF-8,可变长编码
		//61 E4 BA BA E9 97 B4 E6 AD A3 E9 81 93 E6 98 AF E6 B2 A7 E6 A1 91
		fmt.Printf("%X ",b)
	}
	fmt.Println()

	// 直接遍历字符串
	for i,ch := range s {
		// i是索引, ch实际上就是rune
		// (0 61)(1 4EBA)(4 95F4)(7 6B63)(10 9053)(13 662F)(16 6CA7)(19 6851)
		fmt.Printf("(%d %X)",i,ch)
	}
	fmt.Println()

	// 获取字符数量
	fmt.Println("Rune count:",
		utf8.RuneCountInString(s))

	bytes := []byte(s)
	for len(bytes) > 0{
		// 该函数返回bytes中第一个utf8编码的值,和该值的字节数
		ch, size := utf8.DecodeRune(bytes)
		bytes = bytes[size:]
		// a 人 间 正 道 是 沧 桑
		fmt.Printf("%c ",ch)
	}
	fmt.Println()


	for i, ch := range []rune(s) {
		//(0 a)(1 人)(2 间)(3 正)(4 道)(5 是)(6 沧)(7 桑)
		fmt.Printf("(%d %c)",i,ch)
	}
	fmt.Println()


}



/**
 	寻找字符串中最大的不含重复字符的子串
	例如 "abcabcb" -> "abc"
	"bbbb" -> "b"

	思路,
	有方法lastOccurred(X),可以查找字符串最后出现的位置
	一个start,记录
	对每个字符X,
	如果lastOccurred(X)==nil ||  < start 无需操作
	如果lastOccurred(X) >= start,start
	更新lastOccurred(X),更新maxLength


	将字符串转为rune,而不是转为byte,则可以支持中文

 */
func TestA(T *testing.T) {
	s := "bbbbb"
	// 记录每个字符最后出现的位置
	lastOccurred := make(map[rune]int)
	start := 0
	maxLength := 0
	//遍历string
	for i,ch := range []rune(s) {
		if lastI,ok := lastOccurred[ch];
			ok && lastI >= start{
			start = lastOccurred[ch]+1
		}
		if i - start + 1 > maxLength {
			maxLength = i -start + 1
		}
		lastOccurred[ch] = i
	}

	fmt.Println(maxLength)
}