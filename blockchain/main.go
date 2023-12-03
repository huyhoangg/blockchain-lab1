package main

import (
	bc "blockchain/block"
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
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

		// Khởi động server P2P của mỗi node
		go startServer(node, &chain, i)
	}
	// transaction := bc.Transaction{Data: []byte("node 0 create transaction")}
	// broadcastTransaction(nodes[0], transaction)

	select {}

}

func startServer(address string, chain *bc.BlockChain, nodeID int) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Node %d started. Listening on %s\n", nodeID, address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// Xử lý kết nối đến từ client
		go handleClient(conn, chain, nodeID)
	}
}

// func handleClient(conn net.Conn) {
// 	var buf bytes.Buffer
// 	_, err := buf.ReadFrom(conn)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	// Xử lý thông điệp từ client
// 	var receivedTransaction bc.Transaction
// 	decoder := gob.NewDecoder(&buf)
// 	err = decoder.Decode(&receivedTransaction)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	fmt.Printf("Received transaction: %s\n", string(receivedTransaction.Data))
// }

func handleClient(conn net.Conn, chain *bc.BlockChain, nodeID int) {
	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	// Sử dụng bufio.Scanner để đọc từ kết nối
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		receivedMsg := scanner.Text()
		fmt.Printf("Received message from client: %s\n", receivedMsg)

		if receivedMsg == "printchain" {
			// Gửi thông tin blockchain cho client
			for _, block := range chain.Blocks {
				blockInfo := fmt.Sprintf("Timestamp: %d\nPrev. hash: %x\n", block.Timestamp, block.PrevBlockHash)
				for _, transaction := range block.Transactions {
					blockInfo += fmt.Sprintf("Transaction: %s\n", string(transaction.Data))
				}
				blockInfo += fmt.Sprintf("Hash: %x\n", block.Hash)

				_, err := fmt.Fprintf(conn, "%s\n", blockInfo)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}

		if receivedMsg == "hello" {
			responseMsg := "Hello from node!"
			_, err := fmt.Fprintf(conn, "%s\n", responseMsg)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

}

func broadcastTransaction(nodeAddress string, transaction bc.Transaction) {
	for _, node := range nodes {
		if node != nodeAddress {
			go sendTransaction(node, transaction)
		}
	}
}

func sendTransaction(nodeAddress string, transaction bc.Transaction) {
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(transaction)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Transaction sent to %s\n", nodeAddress)
}
