package main

import (
	"fmt"
)

/*并发模拟*/

func main() {
	channel := make(chan string)
	for i := 0; i < 5000; i++ {
		go printHelloWorld(i,channel)
	}

	for{
		message := <- channel
		fmt.Println(message)
	}
}

func printHelloWorld(i int, channel chan string) {
	for  {
		channel <- fmt.Sprintln("Hello World" ,i)
	}
}
