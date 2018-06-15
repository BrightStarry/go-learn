package main

import (
	"io/ioutil"
	"errors"
	"fmt"
	"os"
)

/*异常 和 异常处理*/

func main() {

	_,err:=ioutil.ReadFile("")
	if err != nil {
		//.. 处理异常
	}


	/*
	异常接口
	type error interface {
		Error() string
	}
	*/
	// 构造抛出异常，一般是将这个构造出来的error作为返回值返回，此处不会有任何打印
	errors.New("发生异常")


	// 结构体toString
	fmt.Printf("Data结构体:%v \n",Data{})
	// fmt
	// 输出为一个string
	_ = fmt.Sprintf("dfdf %v", 3.2323)

	// 输出到一个输出流中
	fmt.Fprintln(os.Stdout, "dfdfd")




	// 如果在panic被调用前，或前一个方法处，增加一个类似捕获的机制recover()，则可以处理该恐慌，让程序继续运行
	// 并且该recover()必须在defer中，执行。因为一旦抛出panic，则程序停止，只有defer关键字中的代码会最后被执行到
	// 并且在恢复后，该方法内，后续的代码都不会被执行，只有上一层，调用该方法的方法的后续代码可以被继续执行
	defer func(){
		if p:=recover(); p!= nil{
			fmt.Println("Fatal error:",p)
		}
	}()

	// 抛出异常
	throw()

	println("222")



}

// 结构体
type Data struct{

}
// 定义结构体的String()方法,就类似java中的toString，可以在该结构体被输出时，转为string
func (self  Data)String()string{
	return "Date"
}

func throw() {
	// 恐慌：直接抛出异常，如果没有上文的 recover(), 程序直接停止
	panic(errors.New("程序崩溃"))
}