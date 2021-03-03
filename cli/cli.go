package cli

import (
	"ege/blockchain"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("getbalance -address ADDRESS - get the balance for given address")
	fmt.Println("createblockchain -address ADDRESS creates blockchain for given address")
	fmt.Println("add -block BLOCK_DATA - add a block to the chain")
	fmt.Println("printchain - Prints all the blocks in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - Send amount of coins given amount to given address from given to")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Prev Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s \n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}

}

func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished!")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success")
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address")
	createBlockChainAddress := createBlockChainCmd.String("address", "", "BC address")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()

	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()

		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainAddress == "" {
			createBlockChainCmd.Usage()
			runtime.Goexit()

		}
		cli.createBlockChain(*createBlockChainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" && *sendTo == "" {
			sendCmd.Usage()
			runtime.Goexit()

		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}
