package main

import (
	"testing"
	"os"
	"bufio"
	"fmt"
)

/*模拟数据*/

/*不停写入数据到指定文件*/


/*
*/
func TestMockData(t  *testing.T) {
	file,err := os.Open("C:\\Users\\97038\\Desktop\\access.log")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)


	file2,err := os.OpenFile("./access.log",os.O_APPEND,os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer file2.Close()
	writer := bufio.NewWriter(file2)

	for i:=0;i<=100;i++{
		line,_  := reader.ReadBytes('\n')
		n,err := writer.Write(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(n)

	}
	writer.Flush()


}
