package main

import (
	"devisions.org/go-blockchain/blockchain"
	"fmt"
	"strconv"
)

func main() {

	chain := blockchain.InitBlockchain()

	chain.AddBlock("First block")
	chain.AddBlock("Second block")
	chain.AddBlock("Third block")

	for _, block := range chain.Blocks {
		fmt.Printf("  pref hash: %x\n", block.PrevHash)
		fmt.Printf("  curr hash: %x\n", block.Hash)
		fmt.Printf("       data: %x\n", block.Data)

		pow := blockchain.NewProof(block)
		fmt.Printf("        pow: %s\n\n", strconv.FormatBool(pow.Validate()))
	}

}
