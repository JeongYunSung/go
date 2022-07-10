package main

import (
	"fmt"
	"github.com/jys/learngo/banking"
)

func main() {
	account := banking.NewAccount("John", 1000)
	account1 := banking.NewAccount("John", 1000)
	fmt.Println(*account, *account1)
}
