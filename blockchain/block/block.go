package block

import (
	"crypto/sha256"
	"time"
	"fmt"
	"bytes"
)

type Block struct {
	Timestamp int64
	Transactions []*Transaction
	PrevBlockHash []byte
	Hash []byte
	MerkleRoot []byte
}

type Transaction struct {
	Data []byte
}

type BlockChain struct {
	Blocks []*Block
}

func (block *Block) setHash() {
	headers := []byte(string(block.PrevBlockHash) + string(HashTransactions(block.Transactions)) + string(block.Timestamp) + string(block.MerkleRoot))
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

func HashTransactions(transactions []*Transaction) []byte {
	var hashTransactions []byte

	for _, transaction := range transactions {
		hashTransaction := sha256.Sum256(transaction.Data)
		hashTransactions = append(hashTransactions, hashTransaction[:]...)
	}

	finalHash := sha256.Sum256(hashTransactions)

	return finalHash[:]
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var preBlockHash []byte

	chain_size := len(chain.Blocks)
	if (chain_size > 0 ) {
		preBlockHash = chain.Blocks[chain_size - 1].Hash
	}

	newBlock := &Block{
		Timestamp: time.Now().Unix(),
		PrevBlockHash: preBlockHash,
		Transactions: transactions,
		MerkleRoot:    nil, 
	}

	newBlock.MerkleRoot = BuildMerkleRoot(transactions) 
	newBlock.setHash()
	chain.Blocks = append(chain.Blocks, newBlock)
}

func PrintBlockchain(chain *BlockChain) {
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

func VerifyTransactionInBlock(block *Block, transactionData []byte) bool {
	calculatedMerkleRoot := BuildMerkleRoot(block.Transactions)
	if !bytes.Equal(calculatedMerkleRoot, block.MerkleRoot) {
		return false
	}

	for _, tx := range block.Transactions {
		if bytes.Equal(tx.Data, transactionData) {
			return true
		}
	}

	return false
}
