package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type ScriptResult struct {
	AgentName string
	Output    string
	Success   bool
	Duration  time.Duration
}

type ScriptManager struct {
	agents []string
	ports  []int
}

func NewScriptManager() *ScriptManager {
	return &ScriptManager{
		agents: []string{"agent1", "agent2", "agent3"},
		ports:  []int{9001, 9002, 9003},
	}
}

func (sm *ScriptManager) ExecuteScript(scriptPath string) []ScriptResult {
	results := make([]ScriptResult, 0)

	fmt.Printf("🚀 Executing script: %s\n", scriptPath)
	fmt.Println(strings.Repeat("=", 50))

	// Script dosyasının içeriğini oku
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		fmt.Printf("❌ Error reading script: %v\n", err)
		return results
	}

	// Tüm agent'lara script içeriğini gönder
	resultChan := make(chan ScriptResult, len(sm.agents))
	for i, agent := range sm.agents {
		go func(agentName string, port int, script string) {
			result := sm.executeOnAgent(agentName, port, script)
			resultChan <- result
		}(agent, sm.ports[i], string(scriptContent))
	}

	for i := 0; i < len(sm.agents); i++ {
		result := <-resultChan
		results = append(results, result)
	}

	return results
}

func (sm *ScriptManager) executeOnAgent(agentName string, port int, script string) ScriptResult {
	start := time.Now()

	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return ScriptResult{
			AgentName: agentName,
			Output:    fmt.Sprintf("❌ Connection failed: %v", err),
			Success:   false,
			Duration:  time.Since(start),
		}
	}
	defer conn.Close()

	// Send script content
	fmt.Printf("[DEBUG] Sending script to %s, length: %d\n", agentName, len(script))
	_, err = conn.Write([]byte(script))
	if err != nil {
		return ScriptResult{
			AgentName: agentName,
			Output:    fmt.Sprintf("❌ Failed to send script: %v", err),
			Success:   false,
			Duration:  time.Since(start),
		}
	}

	// Close write side to signal end of data
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}

	// Read response
	output, err := io.ReadAll(conn)
	if err != nil {
		return ScriptResult{
			AgentName: agentName,
			Output:    fmt.Sprintf("❌ Failed to read response: %v", err),
			Success:   false,
			Duration:  time.Since(start),
		}
	}

	return ScriptResult{
		AgentName: agentName,
		Output:    string(output),
		Success:   true,
		Duration:  time.Since(start),
	}
}

func (sm *ScriptManager) PrintResults(results []ScriptResult) {
	fmt.Println("\n📊 SCRIPT EXECUTION RESULTS")
	fmt.Println(strings.Repeat("=", 50))

	for _, result := range results {
		fmt.Printf("\n📋 Agent: %s\n", result.AgentName)
		fmt.Println(strings.Repeat("-", 30))

		if result.Success {
			fmt.Printf("✅ Success (Duration: %v)\n", result.Duration)
		} else {
			fmt.Printf("❌ Failed (Duration: %v)\n", result.Duration)
		}

		fmt.Printf("📄 Output:\n%s\n", result.Output)
	}

	// Summary
	successCount := 0
	totalDuration := time.Duration(0)

	for _, result := range results {
		if result.Success {
			successCount++
		}
		totalDuration += result.Duration
	}

	fmt.Printf("\n📈 SUMMARY:\n")
	fmt.Printf("✅ Successful: %d/%d\n", successCount, len(results))
	fmt.Printf("⏱️  Total Duration: %v\n", totalDuration)
	fmt.Printf("📊 Average Duration: %v\n", totalDuration/time.Duration(len(results)))
}

func main() {
	// Change to parent directory to access scripts folder
	if err := os.Chdir(".."); err != nil {
		fmt.Printf("❌ Error changing directory: %v\n", err)
		return
	}

	sm := NewScriptManager()

	fmt.Println("🎯 Advanced Script Manager")
	fmt.Println("Available scripts:")
	fmt.Println("  HOST SCRIPTS (for physical machine):")
	fmt.Println("    - scripts/host/system_monitor.sh")
	fmt.Println("    - scripts/host/security_scan.sh")
	fmt.Println("    - scripts/host/system_info.sh")
	fmt.Println("  CONTAINER SCRIPTS (for Docker containers):")
	fmt.Println("    - scripts/container/container_monitor.sh")
	fmt.Println("    - scripts/container/container_security.sh")
	fmt.Println("    - scripts/container/backup_files.sh")
	fmt.Println("    - scripts/container/cleanup_logs.sh")
	fmt.Println("    - scripts/container/security_check.sh")
	fmt.Println("  - exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("💻 Enter script path: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		if input == "" {
			continue
		}

		// Check if file exists
		if _, err := os.Stat(input); os.IsNotExist(err) {
			fmt.Printf("❌ Script not found: %s\n", input)
			continue
		}

		// Execute script
		results := sm.ExecuteScript(input)
		sm.PrintResults(results)
	}
}
