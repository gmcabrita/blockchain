package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/pkg/errors"
)

const subsidy = 7

// Transaction represents a blockchain transaction
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// TXInput represents the inputs of a transaction
type TXInput struct {
	Txid     []byte
	Vout     int
	ScritSig string
}

// TXOutput represents the outputs of a transaction
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// NewCoinbaseTX generates a coinbase transaction
func NewCoinbaseTX(to, data string) (*Transaction, error) {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}

	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	err := tx.SetID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to set transaction ID")
	}

	return &tx, nil
}

// SetID sets ID of a transaction
func (tx *Transaction) SetID() error {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		return errors.Wrap(err, "failed to serialize transaction struct")
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]

	return nil
}
