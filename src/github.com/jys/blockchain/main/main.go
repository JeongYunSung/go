package main

import (
	"block"
	"fmt"
	"log"
	"strconv"
)

func main() {
	//createWallet()
	//createBlockchain("1L6eeATb88Db7AZYXvA7Euxp6VMQ8myrN3")
	//printChain()
	//send("1L6eeATb88Db7AZYXvA7Euxp6VMQ8myrN3", "1EYdEoeNMSx66fAA59JSj3fvKEPSXg7yZo", 1)
	//getBalance("1L6eeATb88Db7AZYXvA7Euxp6VMQ8myrN3")
	//getBalance("1EYdEoeNMSx66fAA59JSj3fvKEPSXg7yZo")
}

func getBalance(address string) {
	if !block.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := block.NewBlockchain(address)
	defer bc.Close()

	balance := 0
	pubKeyHash := block.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func send(from, to string, amount int) {
	if !block.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !block.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := block.NewBlockchain(from)
	defer bc.Close()

	tx, err := block.NewUTXOTransaction(from, to, amount, bc)
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}

	bc.MineBlock([]*block.Transaction{tx})
	fmt.Println("Success!")
}

func createWallet() {
	wallets, err := block.NewWallets()
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
	address, err := wallets.CreateWallet()
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
	if err := wallets.SaveToFile(); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}

	fmt.Printf("Your new address: %s\n", address)
}

func createBlockchain(address string) {
	if !block.ValidateAddress(address) {
		log.Fatalf("ERROR: Address is not valid\n")
	}
	bc := block.CreateBlockchain(address)
	bc.Close()
	fmt.Println("Done!")
}

func listAddresses() {
	wallets, err := block.NewWallets()
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func printChain() {
	bc := block.NewBlockchain("")
	defer bc.Close()

	bci := bc.Iterator()

	for {
		b := bci.Next()

		fmt.Printf("============ Block %x ============\n", b.Hash)
		fmt.Printf("Prev. b: %x\n", b.PrevBlockHash)
		pow := block.NewProofOfWork(b)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range b.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}
}
