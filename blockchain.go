package main

// Blockchain represents a chain of Block structs
type Blockchain struct {
	blocks []*Block
}

// NewBlockchain creates a new Blockchain, inserting a genesis block into it
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewBlock("init", []byte{})}}
}

// AddBlock adds a new block to the Blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
