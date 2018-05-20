package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
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

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) (*Transaction, error) {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs, err := bc.FindSpendableOutputs(from, amount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find spendable outputs")
	}

	if acc < amount {
		return nil, errors.New("not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode txid to string")
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	err = tx.SetID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to set transaction id")
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

// IsCoinbase checks if the transaction is a coinbase transaction
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// TXInput represents the inputs of a transaction
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// TXOutput represents the outputs of a transaction
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
