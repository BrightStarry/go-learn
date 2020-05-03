package main

import (
	"zx/h/test/util"
	"github.com/gpmgo/gopm/modules/log"
	"fmt"
	"strings"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("error:", err)
		}
		util.GetParam("exit")
	}()
	//path := `E:\新建文本文档.txt`
	//err := os.Rename(path, `E:\NNPJ-260 ナンパJAPAN検証企画！「絆を深めるためには混浴が一番って知ってましたか？」 オフィス街で声をかけた男上司と女部下が二人きりで初めての混浴体験！巨乳編！！但し用意された水着は極小マイクロビキニのみ！場所はラブホテルのジャグジー！.txt`)
	//if err != nil {
	//	panic(err)
	//}
	lenLimit:= 240
	newFileName := `NNPJ-260 ナンパJAPAN検証企画！「絆を深めるためには混浴が一番って知ってましたか？」 オフィス街で声をかけた男上司と女部下が二人きりで初めての混浴体験！巨乳編！！但し用意された水着は極小マイクロビキニのみ！場所はラブホテルのジャグジー！`
	// 截取开始索引
	startIndex := 40
	fmt.Println(len(newFileName))
	if len(newFileName) > lenLimit{
		// 差值
		diff := len(newFileName) - lenLimit +6
		fmt.Println(len( newFileName[startIndex + diff:]))
		newFileName = strings.TrimRight(newFileName[:startIndex],"�") + "……" + strings.TrimLeft(newFileName[startIndex + diff:],"�")
	}
	fmt.Println(newFileName)
	fmt.Println(len(newFileName))
	fmt.Println(len("……"))

}
