package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"encoding/gob"
	"bytes"
	"math/rand"
	"time"
)

type Message struct {
	Content string
	Source string
}

var wordList = []string{"apple", "banana", "orange", "grape", "melon", "pineapple"}
var nodes = []string{"127.0.0.1:9001", "127.0.0.1:9002", "127.0.0.1:9003"}

func main() {
	var selectedNode string

	fmt.Println("Available nodes:")
	for i, node := range nodes {
		fmt.Printf("[%d] %s\n", i, node)
	}

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Select a node by entering its number ( to ", len(nodes), "): ")
		input, _ := r.ReadString('\n')
		input = strings.TrimSpace(input)

		nodeIndex := 0
		_, err := fmt.Sscanf(input, "%d", &nodeIndex)
		if err != nil || nodeIndex < 1 || nodeIndex > len(nodes) {
			fmt.Println("Invalid input. Please enter a valid node number.")
			continue
		}

		selectedNode = nodes[nodeIndex-1]
		break
	}

	fmt.Println("Connecting to node:", selectedNode)

	conn, err := net.Dial("tcp", selectedNode)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to node:", selectedNode)

	fmt.Println("Blockchain Client CLI - Commands: create, printchain, hello, exit")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		// ... (các lựa chọn khác của CLI)

		case "hello":
			helloMess := Message{Source: "client", Content: "hello"}
			err := SendStructMessage(conn, helloMess)
			if err != nil {
				log.Fatal(err)
			}

			receivedMessage, err := ReceiveStructMessage(conn)
			if err != nil {
				log.Println(err)
				continue
			}
	
			fmt.Println("Received message from client:", receivedMessage.Content)

		
		case "create":
			fmt.Println("creating transactions...")
			for i := 0; i < 3; i++ {
				randomIndex1 := rand.Intn(len(wordList))
				randomIndex2 := rand.Intn(len(wordList))

				randomData := wordList[randomIndex1] + " " + wordList[randomIndex2]
				transactionMessage := Message{Source: "client", Content: "tx:" + randomData}

				err = SendStructMessage(conn, transactionMessage)
				if err != nil {
					log.Fatal(err)
				}
		
				fmt.Println("Transaction created and sent to node:", randomData)
				time.Sleep(800 * time.Millisecond)
			}
	

		case "printchain":
			req := Message{Source: "client", Content: "printchain"}
			err := SendStructMessage(conn, req)
			if err != nil {
				log.Fatal(err)
			}

			receivedMessage, err := ReceiveStructMessage(conn)
			if err != nil {
				log.Println(err)
				continue
			}
	
			fmt.Println("Received message from client:", receivedMessage.Content)
			
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


// // Function to generate a random transaction
// func CreateRandomTransaction() *bc.Transaction {
// 	// Generate a random index to select a word from the list
	

// 	// Create a transaction with random data
// 	fmt.Println(randomData)
// 	return &bc.Transaction{
// 		Data: []byte(randomData),
// 	}
// }