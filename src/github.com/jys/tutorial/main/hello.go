package main

import (
	"fmt"
	"jys/tutorial/main/tutorial"
	"log"
)

func main() {

	names := []string{"Gladys", "Samantha", "Darrin"}

	messages, err := tutorial.Hellos(names)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messages)
}
