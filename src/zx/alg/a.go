package main

import (
	"fmt"
	"time"
)

func main() {

	go b()

	fmt.Println("5")
	time.Sleep(time.Hour)
}

var aa = [...]func(){c}

func b()  {
	i :=0
	defer func() {
		if err:= recover();err != nil{
			fmt.Println("---",i)
		}
	}()
	i++
	for _,v := range aa{
		v()
	}
}

func c() {

	fmt.Println("1")
	panic("2")
	fmt.Println("3")
}
