package main

import (
	"fmt"
	"log"
)

func main() {
	bc, err := NewBlockchain()
	if err != nil {
		panic("Failed to initialize blockchain")
	}

	err = bc.AddBlock("Send 1 ZEN to Pony")
	if err != nil {
		log.Println("Failed send 1 ZEN to Pony")
	}

	err = bc.AddBlock("Send 2 ZEN to Const")
	if err != nil {
		log.Println("Failed send 2 ZEN to Const")
	}

	i := bc.Iterator()

	for {
		block, err := i.Next()
		if err != nil {
			log.Println("Failed to read block from blockchain iterator")
			break
		}

		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("Valid: %v\n", pow.Validate())
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
