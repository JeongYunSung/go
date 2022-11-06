package main

import (
	"block"
	"blockchain"
	"fmt"
	"strconv"
)

func main() {
	bc := blockchain.NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, b := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Printf("Nonce: %d\n", b.Nonce)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(block.NewProofOfWork(b).Validate()))
		fmt.Println()
	}
}
