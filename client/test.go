package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Transaction represents a transaction in the block
type Transaction struct {
	Data []byte
}

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	MerkleRoot    []byte
}

// Blockchain represents the chain of blocks
type Blockchain struct {
	Blocks []*Block
}

// BuildMerkleRoot calculates the Merkle Root from a list of transactions
func BuildMerkleRoot(transactions []*Transaction) []byte {
	var hashes [][]byte

	for _, tx := range transactions {
		hash := sha256.Sum256(tx.Data)
		hashes = append(hashes, hash[:])
	}

	for len(hashes) > 1 {
		var newHashes [][]byte
		for i := 0; i < len(hashes); i += 2 {
			var combined []byte
			if i+1 < len(hashes) {
				combined = append(hashes[i], hashes[i+1]...)
			} else {
				// Duplicate the last hash to create a pair
				combined = append(hashes[i], hashes[i]...)
			}
			hash := sha256.Sum256(combined)
			newHashes = append(newHashes, hash[:])
		}
		hashes = newHashes
	}

	return hashes[0]
}

func main1() {
	// Sample transactions
	transactions := []*Transaction{
		{Data: []byte("Transaction 1")},
		{Data: []byte("Transaction 2")},
		{Data: []byte("Transaction 3")},
	}

	// Create a new block
	newBlock := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: []byte("PreviousBlockHash"),
	}

	// Calculate Merkle Root for the block's transactions
	newBlock.MerkleRoot = BuildMerkleRoot(transactions)

	// Print Merkle Root as hex

	// Add the block to the blockchain
	blockchain := &Blockchain{
		Blocks: []*Block{newBlock},
	}

	PrintBlockchain(blockchain)
}

func PrintBlockchain(chain *Blockchain) {
	// In ra thông tin của blockchain
	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)

		fmt.Println("Transactions:")
		for _, transaction := range block.Transactions {
			fmt.Printf("- %s\n", string(transaction.Data))
		}
		fmt.Printf("MerkleRoot: %x\n", block.MerkleRoot)
		fmt.Printf("hash: %x\n", block.Hash)
	}
}