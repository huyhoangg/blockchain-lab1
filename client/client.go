package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
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
		fmt.Print("Select a node by entering its number (0 to ", len(nodes)-1, "): ")
		input, _ := r.ReadString('\n')
		input = strings.TrimSpace(input)

		nodeIndex := 0
		_, err := fmt.Sscanf(input, "%d", &nodeIndex)
		if err != nil || nodeIndex < 0 || nodeIndex >= len(nodes) {
			fmt.Println("Invalid input. Please enter a valid node number.")
			continue
		}

		selectedNode = nodes[nodeIndex]
		break
	}

	fmt.Println("Connecting to node:", selectedNode)

	conn, err := net.Dial("tcp", selectedNode)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to node:", selectedNode)

	fmt.Println("Blockchain Client CLI - Commands: hello, create, printchain, verify, exit")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {

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
			
		case "verify":
			fmt.Println("Enter 'verify <block index> <transaction data>':")
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			verifyMessage := Message{Source: "client", Content: input}

			err := SendStructMessage(conn, verifyMessage)
			if err != nil {
				log.Fatal(err)
			}

			receivedVerifyResult, err := ReceiveStructMessage(conn)
			if err != nil {
				log.Println(err)
				continue
			}

			fmt.Println("Verification result:", receivedVerifyResult.Content)

		case "exit":
			fmt.Println("Exiting...")
			return
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
