package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// RunCLI is the entrypoint of the CLI tool
func RunCLI() error {
	err := validateArgs()
	if err != nil {
		return err
	}

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ContinueOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ContinueOnError)
	sendCmd := flag.NewFlagSet("send", flag.ContinueOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ContinueOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get the balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send the genesis block to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		return handleGetBalance(getBalanceCmd, getBalanceAddress)
	case "createblockchain":
		return handleCreateBlockchain(createBlockchainCmd, createBlockchainAddress)
	case "printchain":
		return handlePrintChain(printChainCmd)
	case "send":
		return handleSend(sendCmd, sendFrom, sendTo, sendAmount)
	default:
		printUsage()
		return errors.New("bad args")
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS           - Retrieves the balance of the given address")
	fmt.Println("  createblockchain -address ADDRESS     - Create a new blockchain and generate a genesis block")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send coins from an address to another address")
	fmt.Println("  printchain                            - print all the blocks in the blockchain")
}

func validateArgs() error {
	if len(os.Args) < 2 {
		return errors.New("insufficient args")
	}

	return nil
}

func handleGetBalance(getBalanceCmd *flag.FlagSet, getBalanceAddress *string) error {
	err := getBalanceCmd.Parse(os.Args[2:])
	if err != nil {
		return errors.New("failed to parse getbalance command")
	}

	if *getBalanceAddress == "" {
		getBalanceCmd.Usage()
		return errors.New("bad address")
	}

	err = getBalance(*getBalanceAddress)
	return err
}

func handleCreateBlockchain(createBlockchainCmd *flag.FlagSet, createBlockchainAddress *string) error {
	err := createBlockchainCmd.Parse(os.Args[2:])
	if err != nil {
		return errors.New("failed to parse createblockchain command")
	}

	if *createBlockchainAddress == "" {
		createBlockchainCmd.Usage()
		return errors.New("bad address")
	}

	err = createBlockchain(*createBlockchainAddress)
	return err
}

func handleSend(sendCmd *flag.FlagSet, sendFrom, sendTo *string, sendAmount *int) error {
	err := sendCmd.Parse(os.Args[2:])
	if err != nil {
		return errors.New("failed to parse send command")
	}

	if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
		sendCmd.Usage()
		return errors.New("bad args")
	}

	err = send(*sendFrom, *sendTo, *sendAmount)
	return err
}

func handlePrintChain(printChainCmd *flag.FlagSet) error {
	err := printChainCmd.Parse(os.Args[2:])
	if err != nil {
		return errors.New("failed to parse printchain command")
	}

	err = printChain()
	return err
}

func createBlockchain(address string) error {
	bc, err := CreateBlockchain(address)
	defer func() {
		bc.Close()
	}()
	if err != nil {
		return errors.Wrap(err, "failed to create blockchain")
	}

	return nil
}

func getBalance(address string) error {
	bc, err := NewBlockchain(address)
	defer func() {
		bc.Close()
	}()
	if err != nil {
		return errors.Wrap(err, "failed to read blockchain from disk")
	}

	balance := 0
	UTXOs, err := bc.FindUTXO(address)
	if err != nil {
		return errors.Wrap(err, "failed to find unused transaction outputs")
	}

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
	return nil
}

func printChain() error {
	bc, err := NewBlockchain("")
	defer func() {
		bc.Close()
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

func send(sendFrom, sendTo string, sendAmount int) error {
	bc, err := NewBlockchain(sendFrom)
	defer func() {
		bc.Close()
	}()
	if err != nil {
		return errors.Wrap(err, "failed to read blockchain from disk")
	}

	tx, err := NewUTXOTransaction(sendFrom, sendTo, sendAmount, bc)
	if err != nil {
		return errors.Wrap(err, "failed to create new transaction")
	}

	err = bc.MineBlock([]*Transaction{tx})
	if err != nil {
		return errors.Wrap(err, "failed to mine block")
	}

	return nil
}
