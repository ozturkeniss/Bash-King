package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func sendCommand(port, command string) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		fmt.Printf("❌ Error connecting to port %s: %v\n", port, err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, command+"\n")

	scanner := bufio.NewScanner(conn)
	fmt.Printf("📋 Agent (port %s):\n", port)
	for scanner.Scan() {
		fmt.Printf("  %s\n", scanner.Text())
	}
}

func main() {
	ports := []string{"9001", "9002", "9003"}

	fmt.Println("🎯 Simple Distributed Command Server")
	fmt.Println("Available commands: whoami, pwd, ls -la, ps aux")
	fmt.Println()

	for {
		fmt.Print("💻 Enter command (or 'exit'): ")
		var command string
		fmt.Scanln(&command)

		if command == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		if command == "" {
			continue
		}

		fmt.Printf("\n🚀 Executing '%s' on all agents:\n", command)
		fmt.Println(strings.Repeat("=", 50))

		for _, port := range ports {
			sendCommand(port, command)
			fmt.Println()
		}

		fmt.Println(strings.Repeat("=", 50))
	}
}
