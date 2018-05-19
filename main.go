package main

import (
	"fmt"
	"os"
)

func exit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	var bc *Blockchain
	var err error
	defer func() {
		exit(err)
	}()

	bc, err = NewBlockchain()
	defer func() {
		err := bc.db.Close()
		exit(err)
	}()

	cli := CLI{bc}
	err = cli.Run()
}
