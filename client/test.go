package main

import (
	"crypto/sha256"
	"bytes"
	"fmt"
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


func BuildMerkleRoot(transactions []*Transaction) []byte {
	var hashes [][]byte

	// Chuyển đổi dữ liệu giao dịch thành slice của byte và tính toán hash cho mỗi giao dịch
	for _, tx := range transactions {
			hash := sha256.Sum256(tx.Data)
			hashes = append(hashes, hash[:])
	}

	// Xây dựng Merkle Tree từ danh sách các hash
	for len(hashes) > 1 {
			if len(hashes)%2 != 0 {
					hashes = append(hashes, hashes[len(hashes)-1])
			}
			var newHashes [][]byte
			for i := 0; i < len(hashes); i += 2 {
					combined := append(hashes[i], hashes[i+1]...)
					hash := sha256.Sum256(combined)
					newHashes = append(newHashes, hash[:])
			}
			hashes = newHashes
	}

	return hashes[0]
}

func VerifyTransactionInBlock(block *Block, transactionData []byte) bool {
	// Kiểm tra nếu Merkle Root trong block khớp với Merkle Root được tạo từ danh sách giao dịch
	calculatedMerkleRoot := BuildMerkleRoot(block.Transactions)
	if !bytes.Equal(calculatedMerkleRoot, block.MerkleRoot) {
			return false
	}

	// Tìm kiếm giao dịch trong danh sách các giao dịch của block
	for _, tx := range block.Transactions {
			if bytes.Equal(tx.Data, transactionData) {
					return true // Giao dịch được tìm thấy trong block
			}
	}

	return false // Giao dịch không tồn tại trong block
}


// Hàm main để minh họa việc xác minh giao dịch trong block
func main1() {
	// Khởi tạo block với các thông tin cần thiết
	block := &Block{
			Timestamp:     1638736738, // Thời gian
			Transactions:  []*Transaction{
				{Data: []byte("Transaction 1")},
				{Data: []byte("Transaction 2")},
					// Thêm các hash của giao dịch khác nếu cần
			},
			PrevBlockHash: []byte("PreviousHash"), // Hash của block trước
			Hash:          []byte("BlockHash"),    // Hash của block hiện tại
			MerkleRoot:    nil,                   // Merkle Root sẽ được tính sau
	}
	// Tính toán và gán Merkle Root cho block
	block.MerkleRoot = BuildMerkleRoot(block.Transactions)

	// Gọi hàm VerifyTransactionInBlock để xác minh giao dịch trong block
	transactionToVerify := []byte("Transaction 1")
	if VerifyTransactionInBlock(block, transactionToVerify) {
			fmt.Println("Giao dịch tồn tại trong block")
	} else {
			fmt.Println("Giao dịch không tồn tại trong block")
	}
}
