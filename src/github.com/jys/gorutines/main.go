package main

import (
	"fmt"
	"time"
)

func main() {
	go runTask("TaskA")
	go runTask("TaskB")
	time.Sleep(time.Second * 5)
}

func runTask(input string) {
	for i := 0; i < 10; i++ {
		fmt.Println(input, "is a running", i)
		time.Sleep(time.Second)
	}
}
