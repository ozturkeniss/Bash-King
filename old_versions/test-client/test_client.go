package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run test_client.go <port> <command>")
		fmt.Println("Example: go run test_client.go 9001 'ls -la'")
		return
	}

	port := os.Args[1]
	command := os.Args[2]

	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		fmt.Printf("Error connecting to localhost:%s: %v\n", port, err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, command+"\n")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
