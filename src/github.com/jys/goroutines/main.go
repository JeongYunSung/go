package main

import (
	"fmt"
	"time"
)

func main() {
	channel := make(chan string)
	people := []string{"joe", "jane", "jill", "jim"}
	for index, person := range people {
		go isSexy(person, index, channel)
	}
	for i := 0; i < len(people); i++ {
		result := <-channel
		fmt.Println(result)
	}
}

func isSexy(person string, size int, c chan string) {
	time.Sleep(time.Second * time.Duration(size+1))
	c <- person
}
