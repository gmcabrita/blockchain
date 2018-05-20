package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// CLI represents the command line interface for the blockahin
type CLI struct{}

// Run is the entrypoint of the CLI tool
func (cli *CLI) Run() error {
	err := cli.validateArgs()
	if err != nil {
		return err
	}

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ContinueOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ContinueOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send the genesis block to")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			return errors.New("failed to parse createblockchain command")
		}

		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			return errors.New("bad address")
		}

		err = cli.createBlockchain(*createBlockchainAddress)
		if err != nil {
			return err
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			return errors.New("failed to parse printchain command")
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
	fmt.Println("  createblockchain -address ADDRESS - Create a new blockchain and generate a genesis block")
	fmt.Println("  printchain                        - print all the blocks in the blockchain")
}

func (cli *CLI) validateArgs() error {
	if len(os.Args) < 2 {
		return errors.New("insufficient args")
	}

	return nil
}

func (cli *CLI) createBlockchain(address string) error {
	bc, err := CreateBlockchain(address)
	defer func() {
		if bc != nil && bc.db != nil {
			err := bc.db.Close()

			if err != nil {
				panic(err)
			}
		}
	}()

	if err != nil {
		return errors.Wrap(err, "failed to create blockchain")
	}

	return nil
}

func (cli *CLI) printChain() error {
	bc, err := NewBlockchain("")
	defer func() {
		if bc != nil && bc.db != nil {
			err := bc.db.Close()

			if err != nil {
				panic(err)
			}
		}
	}()

	if err != nil {
		return errors.Wrap(err, "failed to read blockchain from disk")
	}

	i := bc.Iterator()

	for {
		block, err := i.Next()
		if err != nil {
			return errors.Wrap(err, "failed to read block")
		}

		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
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
