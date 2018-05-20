package main

import (
	"fmt"
	"os"
)

func main() {
	cli := CLI{}
	err := cli.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
