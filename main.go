package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

// Block defines a block in the chain
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// NewBlock creates a new Block, given some `data` and a `prevBlockHash`
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()

	return block
}

// SetHash sets the hash of the block based on its `Timestamp`, `Data`, and `PrevBlockHash``
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// Blockchain defines a chain of `Block` structs
type Blockchain struct {
	blocks []*Block
}

// NewBlockchain creates a new Blockchain, inserting into it its genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewBlock("init", []byte{})}}
}

// AddBlock adds a new block to the Blockchain, given some `data`
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 ZEN to Pony")
	bc.AddBlock("Send 2 ZEN to Const")

	for _, block := range bc.blocks {
		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
