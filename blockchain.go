package main

import (
	"os"

	"github.com/pkg/errors"

	bolt "github.com/coreos/bbolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const latestBlockBucket = "l"
const genesisCoinbaseData = "init"

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

func dbExists() bool {
	_, err := os.Stat(dbFile)
	return !os.IsNotExist(err)
}

// NewBlockchain loads a blockchain from disk
func NewBlockchain(address string) (*Blockchain, error) {
	if !dbExists() {
		return nil, errors.New("no blockchain found")
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open db file")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte(latestBlockBucket))

		return nil
	})
	if err != nil {
		return nil, err
	}

	bc := Blockchain{tip, db}

	return &bc, nil
}

// CreateBlockchain creates a blockchain db on disk and generates a genesis block
func CreateBlockchain(address string) (*Blockchain, error) {
	if dbExists() {
		return nil, errors.New("blockchain already exists")
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open db file")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx, err := NewCoinbaseTX(address, genesisCoinbaseData)
		if err != nil {
			return errors.Wrap(err, "failed to generate coinbase transaction")
		}
		genesis := NewGenesisBlock(cbtx)
		serializedGenesis, err := genesis.Serialize()
		if err != nil {
			return errors.Wrap(err, "failed to serialize block")
		}

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return errors.Wrap(err, "failed to create bucket")
		}

		err = b.Put(genesis.Hash, serializedGenesis)
		if err != nil {
			return errors.Wrap(err, "failed to insert value into bucket")
		}

		err = b.Put([]byte(latestBlockBucket), genesis.Hash)
		if err != nil {
			return errors.Wrap(err, "failed to insert value into bucket")
		}

		tip = genesis.Hash

		return nil
	})
	if err != nil {
		return nil, err
	}

	bc := Blockchain{tip, db}

	return &bc, nil
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) error {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(latestBlockBucket))

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to mine block")
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		serializedBlock, err := newBlock.Serialize()
		if err != nil {
			return errors.Wrap(err, "failed to serialize block")
		}

		err = b.Put(newBlock.Hash, serializedBlock)
		if err != nil {
			return errors.Wrap(err, "failed to insert value into bucket")
		}

		err = b.Put([]byte(latestBlockBucket), newBlock.Hash)
		if err != nil {
			return errors.Wrap(err, "failed to insert value into bucket")
		}

		bc.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to mine block")
	}

	return nil
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

		return errors.Wrap(err, "failed to deserialize block")
	})

	i.currentHash = block.PrevBlockHash
	return block, err
}
