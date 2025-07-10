package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type SecurityScan struct {
	Date               string
	Hostname           string
	OpenPorts          []PortInfo
	SuspiciousFiles    []string
	HighCPUProcesses   []ProcessInfo
	NetworkConnections []ConnectionInfo
	SudoLogs           []string
	UserLogins         []string
	NewUsers           []UserInfo
	ModifiedFiles      []string
	UnusualPerms       []string
	SetuidBinaries     []string
	FirewallStatus     string
	ListeningServices  []string
}

type PortInfo struct {
	Protocol string
	Port     string
	Process  string
	PID      string
}

type SecurityScanner struct{}

func NewSecurityScanner() *SecurityScanner {
	return &SecurityScanner{}
}

func (ss *SecurityScanner) PerformSecurityScan() SecurityScan {
	scan := SecurityScan{
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Hostname: ss.getHostname(),
	}

	scan.OpenPorts = ss.getOpenPorts()
	scan.SuspiciousFiles = ss.checkSuspiciousFiles()
	scan.HighCPUProcesses = ss.getHighCPUProcesses()
	scan.NetworkConnections = ss.getNetworkConnections()
	scan.SudoLogs = ss.getSudoLogs()
	scan.UserLogins = ss.getUserLogins()
	scan.NewUsers = ss.getNewUsers()
	scan.ModifiedFiles = ss.getModifiedFiles()
	scan.UnusualPerms = ss.getUnusualPermissions()
	scan.SetuidBinaries = ss.getSetuidBinaries()
	scan.FirewallStatus = ss.getFirewallStatus()
	scan.ListeningServices = ss.getListeningServices()

	return scan
}

func (ss *SecurityScanner) getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func (ss *SecurityScanner) getOpenPorts() []PortInfo {
	var ports []PortInfo

	// Get TCP listening ports
	cmd := exec.Command("ss", "-tlnp")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "LISTEN") {
				fields := strings.Fields(line)
				if len(fields) >= 4 {
					addrParts := strings.Split(fields[3], ":")
					if len(addrParts) >= 2 {
						port := PortInfo{
							Protocol: "TCP",
							Port:     addrParts[len(addrParts)-1],
							Process:  fields[len(fields)-1],
						}
						ports = append(ports, port)
					}
				}
			}
		}
	}

	// Get UDP listening ports
	cmd = exec.Command("ss", "-ulnp")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "UNCONN") {
				fields := strings.Fields(line)
				if len(fields) >= 4 {
					addrParts := strings.Split(fields[3], ":")
					if len(addrParts) >= 2 {
						port := PortInfo{
							Protocol: "UDP",
							Port:     addrParts[len(addrParts)-1],
							Process:  fields[len(fields)-1],
						}
						ports = append(ports, port)
					}
				}
			}
		}
	}

	return ports
}

func (ss *SecurityScanner) checkSuspiciousFiles() []string {
	var suspicious []string

	suspiciousPaths := []string{
		"/tmp/.X11-unix",
		"/dev/.udev",
		"/dev/.initramfs",
		"/lib/libkeyutils.so.1",
		"/lib/libproc.so.1",
		"/lib/libproc.so.2",
	}

	for _, path := range suspiciousPaths {
		if _, err := os.Stat(path); err == nil {
			suspicious = append(suspicious, path)
		}
	}

	return suspicious
}

func (ss *SecurityScanner) getHighCPUProcesses() []ProcessInfo {
	var processes []ProcessInfo

	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				cpu, err := strconv.ParseFloat(fields[2], 64)
				if err == nil && cpu > 50.0 {
					process := ProcessInfo{
						PID:     fields[1],
						User:    fields[0],
						CPU:     cpu,
						Command: strings.Join(fields[10:], " "),
					}
					processes = append(processes, process)
				}
			}
		}
	}

	return processes
}

func (ss *SecurityScanner) getNetworkConnections() []ConnectionInfo {
	var connections []ConnectionInfo

	cmd := exec.Command("ss", "-tun")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		ipCount := make(map[string]int)

		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				addrParts := strings.Split(fields[1], ":")
				if len(addrParts) >= 1 {
					ip := addrParts[0]
					if ip != "*" && ip != "::" {
						ipCount[ip]++
					}
				}
			}
		}

		for ip, count := range ipCount {
			connection := ConnectionInfo{
				RemoteIP: ip,
				Count:    count,
			}
			connections = append(connections, connection)
		}
	}

	return connections
}

func (ss *SecurityScanner) getSudoLogs() []string {
	var logs []string

	logFiles := []string{"/var/log/auth.log", "/var/log/secure"}

	for _, logFile := range logFiles {
		if _, err := os.Stat(logFile); err == nil {
			cmd := exec.Command("grep", "sudo:", logFile)
			output, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(output), "\n")
				for i, line := range lines {
					if i < 20 && line != "" {
						logs = append(logs, line)
					}
				}
			}
		}
	}

	return logs
}

func (ss *SecurityScanner) getUserLogins() []string {
	var logins []string

	logFiles := []string{"/var/log/auth.log", "/var/log/secure"}

	for _, logFile := range logFiles {
		if _, err := os.Stat(logFile); err == nil {
			cmd := exec.Command("grep", "session opened", logFile)
			output, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(output), "\n")
				for i, line := range lines {
					if i < 10 && line != "" {
						logins = append(logins, line)
					}
				}
			}
		}
	}

	return logins
}

func (ss *SecurityScanner) getNewUsers() []UserInfo {
	var users []UserInfo

	cmd := exec.Command("find", "/home", "-maxdepth", "1", "-type", "d", "-mtime", "-30")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line != "" {
				username := filepath.Base(line)
				if username != "lost+found" {
					cmd := exec.Command("stat", "-c", "%y", line)
					statOutput, err := cmd.Output()
					if err == nil {
						created := strings.Fields(string(statOutput))[0]
						user := UserInfo{
							Username: username,
							Created:  created,
						}
						users = append(users, user)
					}
				}
			}
		}
	}

	return users
}

func (ss *SecurityScanner) getModifiedFiles() []string {
	var files []string

	cmd := exec.Command("find", "/etc", "-type", "f", "-mtime", "-7")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if i < 10 && line != "" {
				files = append(files, line)
			}
		}
	}

	return files
}

func (ss *SecurityScanner) getUnusualPermissions() []string {
	var files []string

	cmd := exec.Command("find", "/etc", "-type", "f", "-perm", "/o+w")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if i < 10 && line != "" {
				files = append(files, line)
			}
		}
	}

	return files
}

func (ss *SecurityScanner) getSetuidBinaries() []string {
	var binaries []string

	paths := []string{"/usr/bin", "/usr/sbin", "/bin", "/sbin"}

	for _, path := range paths {
		cmd := exec.Command("find", path, "-type", "f", "-perm", "-4000")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for i, line := range lines {
				if i < 10 && line != "" {
					binaries = append(binaries, line)
				}
			}
		}
	}

	return binaries
}

func (ss *SecurityScanner) getFirewallStatus() string {
	// Check UFW
	cmd := exec.Command("ufw", "status")
	output, err := cmd.Output()
	if err == nil {
		return string(output)
	}

	// Check iptables
	cmd = exec.Command("iptables", "-L")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 20 {
			return strings.Join(lines[:20], "\n")
		}
		return string(output)
	}

	return "No firewall detected"
}

func (ss *SecurityScanner) getListeningServices() []string {
	var services []string

	cmd := exec.Command("ss", "-tlnp")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "0.0.0.0") {
				services = append(services, line)
			}
		}
	}

	return services
}

func (ss *SecurityScanner) PrintSecurityReport() {
	scan := ss.PerformSecurityScan()

	fmt.Println("=== SECURITY SCAN REPORT ===")
	fmt.Printf("Date: %s\n", scan.Date)
	fmt.Printf("Hostname: %s\n", scan.Hostname)
	fmt.Println("==========================================")

	// Open Ports Analysis
	fmt.Println("1. OPEN PORTS ANALYSIS")
	fmt.Println("----------------------")
	fmt.Println("TCP Listening Ports:")
	for _, port := range scan.OpenPorts {
		if port.Protocol == "TCP" {
			fmt.Printf("  %s:%s - %s\n", port.Protocol, port.Port, port.Process)
		}
	}

	fmt.Println("\nUDP Listening Ports:")
	for _, port := range scan.OpenPorts {
		if port.Protocol == "UDP" {
			fmt.Printf("  %s:%s - %s\n", port.Protocol, port.Port, port.Process)
		}
	}

	// Malicious Software Detection
	fmt.Println("\n2. MALICIOUS SOFTWARE DETECTION")
	fmt.Println("-------------------------------")

	fmt.Println("Checking for common rootkit indicators:")
	for _, file := range scan.SuspiciousFiles {
		fmt.Printf("  ⚠️  SUSPICIOUS FILE FOUND: %s\n", file)
	}

	fmt.Println("\nHidden process detection:")
	for _, process := range scan.HighCPUProcesses {
		fmt.Printf("  ⚠️  HIGH CPU PROCESS: %s (%f%%) - %s\n",
			process.PID, process.CPU, process.Command)
	}

	fmt.Println("\nUnusual network connections:")
	for _, conn := range scan.NetworkConnections {
		fmt.Printf("  %s: %d connections\n", conn.RemoteIP, conn.Count)
	}

	// User and Sudo Activity
	fmt.Println("\n3. USER AND SUDO ACTIVITY")
	fmt.Println("-------------------------")

	fmt.Println("Last 20 sudo commands:")
	for _, log := range scan.SudoLogs {
		fmt.Printf("  %s\n", log)
	}

	fmt.Println("\nRecent user logins:")
	for _, login := range scan.UserLogins {
		fmt.Printf("  %s\n", login)
	}

	fmt.Println("\nUsers created in last 30 days:")
	for _, user := range scan.NewUsers {
		fmt.Printf("  New user: %s (created: %s)\n", user.Username, user.Created)
	}

	// System Integrity Checks
	fmt.Println("\n4. SYSTEM INTEGRITY CHECKS")
	fmt.Println("--------------------------")

	fmt.Println("Recently modified system files (last 7 days):")
	for _, file := range scan.ModifiedFiles {
		fmt.Printf("  %s\n", file)
	}

	fmt.Println("\nFiles with unusual permissions:")
	for _, file := range scan.UnusualPerms {
		fmt.Printf("  %s\n", file)
	}

	fmt.Println("\nSetuid binaries:")
	for _, binary := range scan.SetuidBinaries {
		fmt.Printf("  %s\n", binary)
	}

	// Network Security
	fmt.Println("\n5. NETWORK SECURITY")
	fmt.Println("-------------------")

	fmt.Println("Firewall status:")
	fmt.Printf("  %s\n", scan.FirewallStatus)

	fmt.Println("\nServices listening on all interfaces:")
	for _, service := range scan.ListeningServices {
		fmt.Printf("  %s\n", service)
	}

	fmt.Println("\n=== SECURITY SCAN COMPLETED ===")
	fmt.Printf("Scan completed at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
