package main

import (
	"testing"
)

func TestBlockSerialization(t *testing.T) {
	genesisBlock := NewGenesisBlock()
	block := NewBlock("dummy block", genesisBlock.Hash)

	var testCases = []*Block{
		genesisBlock,
		block,
	}

	for _, blk := range testCases {
		serializedBlk, err := blk.Serialize()
		if err != nil {
			t.Errorf("failed to serialize block")
		}

		deserializedBlk, err := DeserializeBlock(serializedBlk)
		if err != nil {
			t.Errorf("failed to deserialize block")
		}

		if !blk.Equal(deserializedBlk) {
			t.Errorf("%+v != %+v\n", blk, deserializedBlk)
		}
	}
}
