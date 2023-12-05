package network

import (
	bc "blockchain/block"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strings"
)
type Message struct {
	Content string
	Source  string
	Transactions []*bc.Transaction
	Blocks []*bc.Block
}


var nodes = []string{"127.0.0.1:9001", "127.0.0.1:9002", "127.0.0.1:9003"}
var transactions []*bc.Transaction


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

		// Xử lý kết nối đến từ client
		go HandleClient(conn, chain, nodeID)
	}
}

func HandleClient(conn net.Conn, chain *bc.BlockChain, nodeID int) {
	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	for {
		receivedMessage, err := ReceiveStructMessage(conn)
		if err != nil {
			log.Println(err)
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
				blockInfo := ""
				for _, block := range chain.Blocks {
					blockInfo += fmt.Sprintf("Timestamp: %d\nPrev. hash: %x\n", block.Timestamp, block.PrevBlockHash)
					for _, transaction := range block.Transactions {
						blockInfo += fmt.Sprintf("Transaction: %s\n", string(transaction.Data))
					}
					blockInfo += fmt.Sprintf("Hash: %x\n", block.Hash)
					blockInfo += "\n"
				}

				responseMessage := Message{Source: "node", Content: blockInfo}
				err = SendStructMessage(conn, responseMessage)
				if err != nil {
					log.Println(err)
				}

			case strings.HasPrefix(receivedMessage.Content, "tx:"):
				fmt.Println("Node", nodeID , "Received transaction from client:", receivedMessage.Content[3:])
				transaction := &bc.Transaction{Data: []byte(receivedMessage.Content[3:])}

				transactions = append(transactions, transaction)
				
				if len(transactions) == 3 {
					message := Message{
						Source:       "node",
						Transactions: transactions,
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

			default:
				fmt.Println("Unknown content from client")
			}
		case "node":
			for _, tx := range receivedMessage.Transactions {
				fmt.Println("Node",nodeID,"received transaction broadcast from another node:", string(tx.Data))
			}
			if (nodeID == 2) {
				fmt.Println("\nNode", nodeID, "is responsibled for creating block for these transactions ...")
				chain.AddBlock(receivedMessage.Transactions)
				fmt.Println("\nNode", nodeID, "new block added to chain ...")
				
				// message := Message{
				// 	Source:       "node",
				// 	Transactions: transactions,
				// }
				// for _, node := range nodes[1:] {
				// 	connn, errr := net.Dial("tcp", node)
				// 	if errr != nil {
				// 		log.Fatal(errr)
				// 	}
				// 	defer conn.Close()

				// 	err = SendStructMessage(connn, message)
				// }
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
	buf := make([]byte, 4096) // Kích thước buffer, bạn có thể điều chỉnh tùy theo kích thước dữ liệu truyền
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
