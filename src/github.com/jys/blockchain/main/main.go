package main

import (
	"block"
	"blockchain"
	"fmt"
	"strconv"
)

func main() {
	bc := blockchain.NewBlockchain()

	//bc.AddBlock("Send 1 BTC to Ivan")
	//bc.AddBlock("Last2")

	iterator := bc.Iterator()

	for b := iterator.Next(); iterator.HasNext(); b = iterator.Next() {
		fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)
		pow := block.NewProofOfWork(b)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
