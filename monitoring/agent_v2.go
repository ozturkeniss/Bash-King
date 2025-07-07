package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

type AgentV2 struct {
	port     string
	hostname string
}

func NewAgentV2(port string) *AgentV2 {
	hostname, _ := exec.Command("hostname").Output()
	return &AgentV2{
		port:     port,
		hostname: strings.TrimSpace(string(hostname)),
	}
}

func (a *AgentV2) getSystemInfo() string {
	cmd := exec.Command("bash", "-c", "free -h | grep Mem | awk '{print \"Memory: \" $3 \"/\" $2}'")
	memory, _ := cmd.Output()

	cmd = exec.Command("bash", "-c", "df -h / | tail -1 | awk '{print \"Disk: \" $5}'")
	disk, _ := cmd.Output()

	cmd = exec.Command("bash", "-c", "uptime | awk -F'load average:' '{print \"Load: \" $2}'")
	load, _ := cmd.Output()

	return fmt.Sprintf("=== %s (Port: %s) ===\n%s%s%s",
		a.hostname, a.port, memory, disk, load)
}

func (a *AgentV2) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Timeout ile veri oku
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	var data []byte
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(conn, "Error reading command: %v\n", err)
			return
		}
		data = append(data, buffer[:n]...)
	}

	cmdStr := strings.TrimSpace(string(data))
	fmt.Printf("[DEBUG] Received data length: %d\n", len(cmdStr))

	// Special command for system info
	if cmdStr == "system_info" {
		conn.Write([]byte(a.getSystemInfo()))
		conn.Write([]byte("\n"))
		return
	}

	// Eğer çok satırlı ise script olarak çalıştır
	if strings.Contains(cmdStr, "\n") || strings.HasPrefix(cmdStr, "#!/") {
		fmt.Printf("[DEBUG] Executing multi-line script, length: %d\n", len(cmdStr))
		tmpFile, err := os.CreateTemp("/tmp", "remote_script_*.sh")
		if err != nil {
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
			fmt.Fprintf(conn, "Command error: %v\n", err)
		}
		fmt.Printf("[DEBUG] Script output length: %d\n", len(output))
		conn.Write(output)
		conn.Write([]byte("\n"))
		return
	}

	// Tek satırlık komut ise eskisi gibi çalıştır
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(conn, "Command error: %v\n", err)
	}
	conn.Write(output)
	conn.Write([]byte("\n"))
}

func (a *AgentV2) startMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Printf("[%s] Auto-monitoring: %s", time.Now().Format("15:04:05"), a.getSystemInfo())
		}
	}()
}

func main() {
	port := "9001"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	agent := NewAgentV2(port)

	// Start auto-monitoring
	agent.startMonitoring()

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Agent V2 listening on port %s...\n", port)
	fmt.Printf("Hostname: %s\n", agent.hostname)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go agent.handleConnection(conn)
	}
}
