package main

import (
	bc "blockchain/block"
	nn "blockchain/network"
)

var nodes = []string{"127.0.0.1:9001", "127.0.0.1:9002", "127.0.0.1:9003"}

func main() {
	chain := bc.BlockChain{}

	transactions1 := []*bc.Transaction{
		{Data: []byte("Transaction 1")},
		{Data: []byte("Transaction 2")},
	}

	transactions2 := []*bc.Transaction{
		{Data: []byte("Transaction 3")},
		{Data: []byte("Transaction 4")},
	}

	chain.AddBlock(transactions1)
	chain.AddBlock(transactions2)

	for i, node := range nodes {

		nodeChain := CopyBlockchain(&chain)

    go nn.StartServer(node, nodeChain, i)
	}

	select {}

}

func CopyBlockchain(source *bc.BlockChain) *bc.BlockChain {
	newChain := bc.BlockChain{}

	for _, block := range source.Blocks {
		newBlock := bc.Block{
			Timestamp:     block.Timestamp,
			Transactions:  block.Transactions[:], 
			PrevBlockHash: block.PrevBlockHash,
			Hash:          block.Hash,
		}

		newChain.Blocks = append(newChain.Blocks, &newBlock)
	}

	return &newChain
}
