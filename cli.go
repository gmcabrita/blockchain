package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// CLI represents the command line interface for the blockahin
type CLI struct {
	bc *Blockchain
}

// Run is the entrypoint of the CLI tool
func (cli *CLI) Run() error {
	err := cli.validateArgs()
	if err != nil {
		return err
	}

	addBlockCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ContinueOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			return errors.New("failed to parse add command")
		}

		if *addBlockData == "" {
			addBlockCmd.Usage()
			return errors.New("bad data")
		}

		err = cli.addBlock(*addBlockData)
		if err != nil {
			return err
		}

		fmt.Println("Successfully added the block!")
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			return errors.New("failed to parse print command")
		}

		err = cli.printChain()
		if err != nil {
			return err
		}
	default:
		cli.printUsage()
		return errors.New("bad args")
	}

	return nil
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  add -data DATA - add a block to the blockchain")
	fmt.Println("  print          - print all the blocks in the blockchain")
}

func (cli *CLI) validateArgs() error {
	if len(os.Args) < 2 {
		return errors.New("insufficient args")
	}

	return nil
}

func (cli *CLI) addBlock(data string) error {
	err := cli.bc.AddBlock(data)
	if err != nil {
		return errors.Wrap(err, "failed to add block")
	}

	return nil
}

func (cli *CLI) printChain() error {
	i := cli.bc.Iterator()

	for {
		block, err := i.Next()
		if err != nil {
			return errors.Wrap(err, "failed to read block")
		}

		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("Valid: %v\n", pow.Validate())

		if len(block.PrevBlockHash) == 0 {
			break
		}

		fmt.Println()
	}

	return nil
}
