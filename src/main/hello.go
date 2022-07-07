package main

import (
	"fmt"
	"jys/main/tutorial"
	"log"
)

func main() {
	message, err := tutorial.Hello("")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
}
