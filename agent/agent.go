package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Read all data until connection closes
	var data []byte
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("[DEBUG] Error reading data: %v\n", err)
			fmt.Fprintf(conn, "Error reading command: %v\n", err)
			return
		}
		data = append(data, buffer[:n]...)
	}

	cmdStr := strings.TrimSpace(string(data))
	fmt.Printf("[DEBUG] Received command length: %d\n", len(cmdStr))
	fmt.Printf("[DEBUG] Command preview: %s...\n", cmdStr[:min(100, len(cmdStr))])

	// Check if it's a multi-line script
	if strings.Contains(cmdStr, "\n") || strings.HasPrefix(cmdStr, "#!/") {
		fmt.Printf("[DEBUG] Executing multi-line script\n")
		// Create temp file for script
		tmpFile, err := os.CreateTemp("/tmp", "agent_script_*.sh")
		if err != nil {
			fmt.Printf("[DEBUG] Error creating temp file: %v\n", err)
			fmt.Fprintf(conn, "Error creating temp file: %v\n", err)
			return
		}
		defer os.Remove(tmpFile.Name())

		tmpFile.WriteString(cmdStr)
		tmpFile.Close()
		os.Chmod(tmpFile.Name(), 0755)

		cmd := exec.Command("bash", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("[DEBUG] Script execution error: %v\n", err)
			fmt.Fprintf(conn, "Command error: %v\n", err)
		}

		fmt.Printf("[DEBUG] Script output length: %d\n", len(output))
		if len(output) > 0 {
			fmt.Printf("[DEBUG] Output preview: %s...\n", string(output)[:min(200, len(output))])
		}

		conn.Write(output)
		conn.Write([]byte("\n"))
	} else {
		fmt.Printf("[DEBUG] Executing single command\n")
		cmd := exec.Command("bash", "-c", cmdStr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("[DEBUG] Command execution error: %v\n", err)
			fmt.Fprintf(conn, "Command error: %v\n", err)
		}

		fmt.Printf("[DEBUG] Command output length: %d\n", len(output))
		if len(output) > 0 {
			fmt.Printf("[DEBUG] Output preview: %s...\n", string(output)[:min(200, len(output))])
		}

		conn.Write(output)
		conn.Write([]byte("\n"))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	port := "9001"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Agent listening on port %s...\n", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
