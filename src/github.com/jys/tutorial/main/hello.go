package main

import (
	"fmt"
	"log"
	"tutorial/main/tutorial"
)

func main() {

	names := []string{"Gladys", "Samantha", "Darrin"}

	messages, err := tutorial.Hellos(names)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messages)
}
