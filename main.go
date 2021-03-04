package main

import (
	"ege/cli"
	"ege/wallet"
	"os"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
	w := wallet.MakeWallet()
	w.Address()
}
