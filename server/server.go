package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type AgentResult struct {
	AgentName string
	Output    string
	Error     error
}

type Server struct {
	agents map[string]string // agent_name -> port
}

func NewServer() *Server {
	return &Server{
		agents: map[string]string{
			"agent1": "9001",
			"agent2": "9002",
			"agent3": "9003",
		},
	}
}

func (s *Server) sendCommandToAgent(agentName, port, command string) AgentResult {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		return AgentResult{
			AgentName: agentName,
			Error:     fmt.Errorf("connection failed: %v", err),
		}
	}
	defer conn.Close()

	// Set timeout
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Send command
	fmt.Fprintf(conn, command+"\n")

	// Read response
	scanner := bufio.NewScanner(conn)
	var output strings.Builder
	for scanner.Scan() {
		output.WriteString(scanner.Text() + "\n")
	}

	return AgentResult{
		AgentName: agentName,
		Output:    output.String(),
		Error:     scanner.Err(),
	}
}

func (s *Server) ExecuteCommandOnAllAgents(command string) {
	fmt.Printf("🚀 Executing command on all agents: %s\n", command)
	fmt.Println(strings.Repeat("=", 50))

	var wg sync.WaitGroup
	results := make(chan AgentResult, len(s.agents))

	// Send command to all agents concurrently
	for agentName, port := range s.agents {
		wg.Add(1)
		go func(name, p string) {
			defer wg.Done()
			result := s.sendCommandToAgent(name, p, command)
			results <- result
		}(agentName, port)
	}

	// Wait for all agents to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and display results
	for result := range results {
		fmt.Printf("\n📋 Agent: %s\n", result.AgentName)
		fmt.Println(strings.Repeat("-", 30))
		if result.Error != nil {
			fmt.Printf("❌ Error: %v\n", result.Error)
		} else {
			fmt.Printf("✅ Output:\n%s", result.Output)
		}
	}
	fmt.Println("\n" + strings.Repeat("=", 50))
}

func main() {
	server := NewServer()

	fmt.Println("🎯 Distributed Command Server Started!")
	fmt.Println("Available commands:")
	fmt.Println("  - system_info")
	fmt.Println("  - whoami")
	fmt.Println("  - pwd")
	fmt.Println("  - ls -la")
	fmt.Println("  - ps aux")
	fmt.Println("  - exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("💻 Enter command: ")
		if !scanner.Scan() {
			break
		}
		command := strings.TrimSpace(scanner.Text())

		if command == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		if command == "" {
			continue
		}

		server.ExecuteCommandOnAllAgents(command)
	}
}
