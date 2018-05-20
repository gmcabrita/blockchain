package main

import (
	"fmt"
	"os"
)

func main() {
	err := RunCLI()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
