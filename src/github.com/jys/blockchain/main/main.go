package main

import (
	"block"
	"fmt"
)

func main() {
	fmt.Println("=== Before ===")
	getBalance("jys")
	getBalance("jys2")

	send("jys", "jys2", 10)

	fmt.Println("=== After ===")
	getBalance("jys")
	getBalance("jys2")
}

func getBalance(address string) {
	bc := block.NewBlockchain(address)

	defer bc.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func send(from, to string, amount int) {
	bc := block.NewBlockchain(from)

	defer bc.Close()

	tx := block.NewUTXOTransaction(from, to, amount, bc)

	bc.MineBlock([]*block.Transaction{tx})

	fmt.Println("Success!")
}
