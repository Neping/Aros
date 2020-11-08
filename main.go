package main

import (
	"chos/utils"
	"chos/blockchain"
)

func main() {
	utils.TestAES()
	cli := blockchain.CLI{}
	cli.Run()
}
