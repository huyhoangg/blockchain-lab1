package network

import (
	bc "blockchain/block"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Content      string
	Source       string
	Transactions []*bc.Transaction
}

var nodes = []string{"127.0.0.1:9001", "127.0.0.1:9002", "127.0.0.1:9003"}

var transactions0 []*bc.Transaction
var transactions1 []*bc.Transaction
var transactions2 []*bc.Transaction

func StartServer(address string, chain *bc.BlockChain, nodeID int) {
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

		go HandleClient(conn, chain, nodeID)
	}
}

func HandleClient(conn net.Conn, chain *bc.BlockChain, nodeID int) {
	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	for {
		receivedMessage, err := ReceiveStructMessage(conn)
		if err != nil {
			if strings.Contains(err.Error(), "forcibly closed") {
				fmt.Println("Connection forcibly closed by client")
				break 
			}
			if err.Error() == "EOF" {
				continue 
		}
			log.Println(err)
			continue
		}

		switch receivedMessage.Source {
		case "client":
			switch {
			case receivedMessage.Content == "hello":
				responseMessage := Message{Source: "node", Content: "helloo from node"}
				err = SendStructMessage(conn, responseMessage)
				if err != nil {
					log.Println(err)
				}

			case receivedMessage.Content == "printchain":
				blockInfo := "/n"
				for i, block := range chain.Blocks {
					blockInfo += fmt.Sprintf("block [%d]\n", i)
					blockInfo += fmt.Sprintf("Timestamp: %d\nPrev. hash: %x\n", block.Timestamp, block.PrevBlockHash)
					for _, transaction := range block.Transactions {
						blockInfo += fmt.Sprintf("Transaction: %s\n", string(transaction.Data))
					}
					blockInfo += fmt.Sprintf("MerkleRoot:  %x\n", block.MerkleRoot)
					blockInfo += fmt.Sprintf("Hash: %x\n", block.Hash)
					blockInfo += "\n"
				}

				responseMessage := Message{Source: "node", Content: blockInfo}
				err = SendStructMessage(conn, responseMessage)
				if err != nil {
					log.Println(err)
				}

			case strings.HasPrefix(receivedMessage.Content, "tx:"):
				fmt.Println("Node", nodeID, "Received transaction from client:", receivedMessage.Content[3:])
				transaction := &bc.Transaction{Data: []byte(receivedMessage.Content[3:])}

				transactions0 = append(transactions0, transaction)

				if len(transactions0) == 3 {
					message := Message{
						Source:       "node",
						Transactions: transactions0,
					}
					for _, node := range nodes[1:] {
						connn, errr := net.Dial("tcp", node)
						if errr != nil {
							log.Fatal(errr)
						}
						defer conn.Close()

						err = SendStructMessage(connn, message)
					}
				}

			case strings.HasPrefix(receivedMessage.Content, "verify"):

				components := strings.Fields(receivedMessage.Content)

				if len(components) < 3 || components[0] != "verify" {
					responseMessage := Message{Source: "node", Content: "Invalid verify request format"}
					err := SendStructMessage(conn, responseMessage)
					if err != nil {
						log.Println(err)
					}
					return
				}

				blockIndex, err := strconv.Atoi(components[1])
				if err != nil || blockIndex < 0 || blockIndex >= len(chain.Blocks) {
					responseMessage := Message{Source: "node", Content: "Invalid block index"}
					err = SendStructMessage(conn, responseMessage)
					if err != nil {
						log.Println(err)
					}
					return
				}

				transactionData := strings.Join(components[2:], " ")

				block := chain.Blocks[blockIndex]
				if bc.VerifyTransactionInBlock(block, []byte(transactionData)) {
					responseMessage := Message{Source: "node", Content: "Transaction exists in the block"}
					err = SendStructMessage(conn, responseMessage)
					if err != nil {
						log.Println(err)
					}
				} else {
					responseMessage := Message{Source: "node", Content: "Transaction does not exist in the block"}
					err = SendStructMessage(conn, responseMessage)
					if err != nil {
						log.Println(err)
					}
				}

			default:
				fmt.Println("Unknown content from client")
			}
		case "node":
			if receivedMessage.Content == "createblock" {
				if nodeID == 1 {
					fmt.Println("\nNode", nodeID, "is updating new block ...")
					chain.AddBlock(transactions1)
					fmt.Println("\nNode", nodeID, "is synced with newest chain ...")

				}
				if nodeID == 0 {
					fmt.Println("\nNode", nodeID, "is updating new block ...")
					chain.AddBlock(transactions0)
					fmt.Println("\nNode", nodeID, "is synced with newest chain ...")
				}
				break
			}

			for _, tx := range receivedMessage.Transactions {
				fmt.Println("Node", nodeID, "received transaction broadcast from another node:", string(tx.Data))
			}
			if nodeID == 2 {
				fmt.Println("\nNode", nodeID, "is responsibled for creating block for these transactions ...")
				chain.AddBlock(receivedMessage.Transactions)
				fmt.Println("\nNode", nodeID, "new block added to chain ...")

				fmt.Println("Sync chain to other nodes ...")

				message := Message{
					Source:  "node",
					Content: "createblock",
				}
				for _, node := range nodes[:2] {
					connn, errr := net.Dial("tcp", node)
					if errr != nil {
						log.Fatal(errr)
					}
					defer conn.Close()

					err = SendStructMessage(connn, message)
					time.Sleep(8 * time.Millisecond)
				}

			}

			if nodeID == 1 {
				for _, tx := range receivedMessage.Transactions {
					transactions1 = append(transactions1, tx)
				}
			}

			// switch {
			// case strings.HasPrefix(receivedMessage.Content, "tx:"):
			// 	fmt.Println("Received transaction broadcast from another node:", receivedMessage.Content[3:])
			// default:
			// 	fmt.Println("Unknown content from another node")
			// }
		default:
			fmt.Println("Unknown source")
		}
	}
}

func SendStructMessage(conn net.Conn, message Message) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(message)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func ReceiveStructMessage(conn net.Conn) (Message, error) {
	var receivedMessage Message
	buf := make([]byte, 4096) 
	n, err := conn.Read(buf)
	if err != nil {
		return receivedMessage, err
	}
	decoder := gob.NewDecoder(bytes.NewReader(buf[:n]))
	err = decoder.Decode(&receivedMessage)
	if err != nil {
		return receivedMessage, err
	}
	return receivedMessage, nil
}
