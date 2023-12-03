package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Blockchain Client CLI - Commands: create, printchain, hello, exit")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		// ... (các lựa chọn khác của CLI)

		case "hello":
			fmt.Println("Sending 'hello'...")
			_, err := fmt.Fprintf(conn, "hello\n")
			if err != nil {
				log.Println(err)
				return
			}

			// Đọc và hiển thị phản hồi từ node
			responseMsg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println("Response from node:", responseMsg)

		case "printchain":
			_, err := fmt.Fprintf(conn, "printchain\n")
			if err != nil {
				log.Println(err)
				return
			}

			scanner := bufio.NewScanner(conn)
			fmt.Println("Response from node:")
			for scanner.Scan() {
				receivedMsg := scanner.Text()
				fmt.Println(receivedMsg)
			}
			// ... (phần còn lại của switch case)
		}
	}
}
