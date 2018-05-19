package main

import (
	bolt "github.com/coreos/bbolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const latestBlockBucket = "l"

// Blockchain represents a chain of blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// BlockchainIterator represents an iterator over the blockchain
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// NewBlockchain builds a new Blockchain.
// If one does not currently exist on disk it creates one and inserts a genesis block into it,
// otherwise it simply reads the existing one from disk.
func NewBlockchain() (*Blockchain, error) {
	var (
		tip []byte
		b   *bolt.Bucket
	)

	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b = tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()

			serializedGenesis, err := genesis.Serialize()
			if err != nil {
				return err
			}

			b, err = tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return err
			}

			err = b.Put(genesis.Hash, serializedGenesis)
			if err != nil {
				return err
			}

			err = b.Put([]byte(latestBlockBucket), genesis.Hash)
			if err != nil {
				return err
			}

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte(latestBlockBucket))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc, err
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) error {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(latestBlockBucket))

		return nil
	})
	if err != nil {
		return err
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		serializedBlock, err := newBlock.Serialize()
		if err != nil {
			return err
		}

		err = b.Put(newBlock.Hash, serializedBlock)
		if err != nil {
			return err
		}

		err = b.Put([]byte(latestBlockBucket), newBlock.Hash)
		if err != nil {
			return err
		}

		bc.tip = newBlock.Hash

		return nil
	})

	return err
}

// Iterator builds a blockchain iterator from a blockchain
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

// Next returns the next block in the blockchain iterator
func (i *BlockchainIterator) Next() (*Block, error) {
	var (
		block *Block
		err   error
	)

	err = i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		serializedBlock := b.Get(i.currentHash)
		block, err = DeserializeBlock(serializedBlock)

		return err
	})

	i.currentHash = block.PrevBlockHash
	return block, err
}
