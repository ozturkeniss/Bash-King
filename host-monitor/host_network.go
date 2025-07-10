package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type NetworkAnalysis struct {
	Date              string
	Hostname          string
	Interfaces        []InterfaceInfo
	RoutingTable      []RouteInfo
	DNSInfo           DNSInfo
	ActiveConnections []ConnectionInfo
	NetworkStats      NetworkStats
	FirewallRules     []string
	BandwidthUsage    []BandwidthInfo
}

type InterfaceInfo struct {
	Name      string
	IP        string
	Netmask   string
	MAC       string
	Status    string
	MTU       string
	RXBytes   int64
	TXBytes   int64
	RXPackets int64
	TXPackets int64
	RXErrors  int64
	TXErrors  int64
}

type RouteInfo struct {
	Destination string
	Gateway     string
	Interface   string
	Flags       string
}

type DNSInfo struct {
	Nameservers []string
	Domain      string
	Search      []string
}

type NetworkStats struct {
	TotalConnections int
	TCPConnections   int
	UDPConnections   int
	Established      int
	Listen           int
}

type BandwidthInfo struct {
	Interface string
	RXRate    int64
	TXRate    int64
	Timestamp time.Time
}

type NetworkAnalyzer struct{}

func NewNetworkAnalyzer() *NetworkAnalyzer {
	return &NetworkAnalyzer{}
}

func (na *NetworkAnalyzer) AnalyzeNetwork() NetworkAnalysis {
	analysis := NetworkAnalysis{
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Hostname: na.getHostname(),
	}

	analysis.Interfaces = na.getInterfaces()
	analysis.RoutingTable = na.getRoutingTable()
	analysis.DNSInfo = na.getDNSInfo()
	analysis.ActiveConnections = na.getActiveConnections()
	analysis.NetworkStats = na.getNetworkStats()
	analysis.FirewallRules = na.getFirewallRules()
	analysis.BandwidthUsage = na.getBandwidthUsage()

	return analysis
}

func (na *NetworkAnalyzer) getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func (na *NetworkAnalyzer) getInterfaces() []InterfaceInfo {
	var interfaces []InterfaceInfo

	// Get network interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return interfaces
	}

	for _, iface := range ifaces {
		if iface.Name == "lo" {
			continue
		}

		info := InterfaceInfo{
			Name:   iface.Name,
			MAC:    iface.HardwareAddr.String(),
			Status: "down",
		}

		if iface.Flags&net.FlagUp != 0 {
			info.Status = "up"
		}

		// Get IP addresses
		addrs, err := iface.Addrs()
		if err == nil && len(addrs) > 0 {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok {
					if ipnet.IP.To4() != nil {
						info.IP = ipnet.IP.String()
						info.Netmask = net.IP(ipnet.Mask).String()
						break
					}
				}
			}
		}

		// Get interface statistics
		info.RXBytes = na.getInterfaceStat(iface.Name, "rx_bytes")
		info.TXBytes = na.getInterfaceStat(iface.Name, "tx_bytes")
		info.RXPackets = na.getInterfaceStat(iface.Name, "rx_packets")
		info.TXPackets = na.getInterfaceStat(iface.Name, "tx_packets")
		info.RXErrors = na.getInterfaceStat(iface.Name, "rx_errors")
		info.TXErrors = na.getInterfaceStat(iface.Name, "tx_errors")

		// Get MTU
		info.MTU = fmt.Sprintf("%d", iface.MTU)

		interfaces = append(interfaces, info)
	}

	return interfaces
}

func (na *NetworkAnalyzer) getInterfaceStat(iface, stat string) int64 {
	data, err := os.ReadFile(fmt.Sprintf("/sys/class/net/%s/statistics/%s", iface, stat))
	if err != nil {
		return 0
	}
	value, _ := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	return value
}

func (na *NetworkAnalyzer) getRoutingTable() []RouteInfo {
	var routes []RouteInfo

	cmd := exec.Command("ip", "route", "show")
	output, err := cmd.Output()
	if err != nil {
		return routes
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 3 {
			route := RouteInfo{
				Destination: fields[0],
			}

			for i, field := range fields {
				if field == "via" && i+1 < len(fields) {
					route.Gateway = fields[i+1]
				} else if field == "dev" && i+1 < len(fields) {
					route.Interface = fields[i+1]
				}
			}

			routes = append(routes, route)
		}
	}

	return routes
}

func (na *NetworkAnalyzer) getDNSInfo() DNSInfo {
	info := DNSInfo{}

	// Read /etc/resolv.conf
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return info
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				info.Nameservers = append(info.Nameservers, fields[1])
			}
		} else if strings.HasPrefix(line, "domain") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				info.Domain = fields[1]
			}
		} else if strings.HasPrefix(line, "search") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				info.Search = append(info.Search, fields[1])
			}
		}
	}

	return info
}

func (na *NetworkAnalyzer) getActiveConnections() []ConnectionInfo {
	var connections []ConnectionInfo

	cmd := exec.Command("ss", "-tunap")
	output, err := cmd.Output()
	if err != nil {
		return connections
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "Netid") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 4 {
			conn := ConnectionInfo{
				Protocol: fields[0],
			}

			// Parse local and remote addresses
			if len(fields) >= 2 {
				conn.LocalAddr = fields[3]
				if len(fields) >= 4 {
					conn.RemoteAddr = fields[4]
				}
			}

			// Get process info
			if len(fields) >= 6 {
				processInfo := fields[5]
				if strings.Contains(processInfo, "pid=") {
					parts := strings.Split(processInfo, ",")
					for _, part := range parts {
						if strings.HasPrefix(part, "pid=") {
							conn.PID = strings.TrimPrefix(part, "pid=")
						} else if strings.HasPrefix(part, "cmd=") {
							conn.Program = strings.TrimPrefix(part, "cmd=")
						}
					}
				}
			}

			connections = append(connections, conn)
		}
	}

	return connections
}

func (na *NetworkAnalyzer) getNetworkStats() NetworkStats {
	stats := NetworkStats{}

	cmd := exec.Command("ss", "-tunap")
	output, err := cmd.Output()
	if err != nil {
		return stats
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "Netid") {
			continue
		}

		stats.TotalConnections++

		fields := strings.Fields(line)
		if len(fields) >= 1 {
			if strings.HasPrefix(fields[0], "tcp") {
				stats.TCPConnections++
			} else if strings.HasPrefix(fields[0], "udp") {
				stats.UDPConnections++
			}
		}

		if strings.Contains(line, "ESTABLISHED") {
			stats.Established++
		} else if strings.Contains(line, "LISTEN") {
			stats.Listen++
		}
	}

	return stats
}

func (na *NetworkAnalyzer) getFirewallRules() []string {
	var rules []string

	// Check iptables
	cmd := exec.Command("iptables", "-L", "-n")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if i < 20 && line != "" {
				rules = append(rules, line)
			}
		}
	}

	// Check UFW
	cmd = exec.Command("ufw", "status")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if i < 10 && line != "" {
				rules = append(rules, "UFW: "+line)
			}
		}
	}

	return rules
}

func (na *NetworkAnalyzer) getBandwidthUsage() []BandwidthInfo {
	var bandwidth []BandwidthInfo

	interfaces, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return bandwidth
	}

	for _, iface := range interfaces {
		if iface.Name() == "lo" {
			continue
		}

		rxBytes := na.getInterfaceStat(iface.Name(), "rx_bytes")
		txBytes := na.getInterfaceStat(iface.Name(), "tx_bytes")

		info := BandwidthInfo{
			Interface: iface.Name(),
			RXRate:    rxBytes,
			TXRate:    txBytes,
			Timestamp: time.Now(),
		}

		bandwidth = append(bandwidth, info)
	}

	return bandwidth
}

func (na *NetworkAnalyzer) PrintNetworkReport() {
	analysis := na.AnalyzeNetwork()

	fmt.Println("=== NETWORK ANALYSIS REPORT ===")
	fmt.Printf("Date: %s\n", analysis.Date)
	fmt.Printf("Hostname: %s\n", analysis.Hostname)
	fmt.Println("==========================================")

	// Network Interfaces
	fmt.Println("\n1. NETWORK INTERFACES")
	fmt.Println("----------------------")
	for _, iface := range analysis.Interfaces {
		fmt.Printf("Interface: %s\n", iface.Name)
		fmt.Printf("  Status: %s\n", iface.Status)
		fmt.Printf("  IP: %s\n", iface.IP)
		fmt.Printf("  Netmask: %s\n", iface.Netmask)
		fmt.Printf("  MAC: %s\n", iface.MAC)
		fmt.Printf("  MTU: %s\n", iface.MTU)
		fmt.Printf("  RX: %d bytes (%d packets, %d errors)\n",
			iface.RXBytes, iface.RXPackets, iface.RXErrors)
		fmt.Printf("  TX: %d bytes (%d packets, %d errors)\n",
			iface.TXBytes, iface.TXPackets, iface.TXErrors)
		fmt.Println()
	}

	// Routing Table
	fmt.Println("2. ROUTING TABLE")
	fmt.Println("----------------")
	for _, route := range analysis.RoutingTable {
		fmt.Printf("Destination: %s, Gateway: %s, Interface: %s\n",
			route.Destination, route.Gateway, route.Interface)
	}

	// DNS Information
	fmt.Println("\n3. DNS CONFIGURATION")
	fmt.Println("-------------------")
	fmt.Printf("Domain: %s\n", analysis.DNSInfo.Domain)
	fmt.Printf("Nameservers: %v\n", analysis.DNSInfo.Nameservers)
	fmt.Printf("Search domains: %v\n", analysis.DNSInfo.Search)

	// Network Statistics
	fmt.Println("\n4. NETWORK STATISTICS")
	fmt.Println("--------------------")
	fmt.Printf("Total connections: %d\n", analysis.NetworkStats.TotalConnections)
	fmt.Printf("TCP connections: %d\n", analysis.NetworkStats.TCPConnections)
	fmt.Printf("UDP connections: %d\n", analysis.NetworkStats.UDPConnections)
	fmt.Printf("Established: %d\n", analysis.NetworkStats.Established)
	fmt.Printf("Listening: %d\n", analysis.NetworkStats.Listen)

	// Active Connections
	fmt.Println("\n5. ACTIVE CONNECTIONS")
	fmt.Println("---------------------")
	for _, conn := range analysis.ActiveConnections {
		fmt.Printf("Protocol: %s, Local: %s, Remote: %s\n",
			conn.Protocol, conn.LocalAddr, conn.RemoteAddr)
		if conn.PID != "" {
			fmt.Printf("  PID: %s, Program: %s\n", conn.PID, conn.Program)
		}
	}

	// Firewall Rules
	fmt.Println("\n6. FIREWALL RULES")
	fmt.Println("-----------------")
	for _, rule := range analysis.FirewallRules {
		fmt.Printf("  %s\n", rule)
	}

	// Bandwidth Usage
	fmt.Println("\n7. BANDWIDTH USAGE")
	fmt.Println("------------------")
	for _, bw := range analysis.BandwidthUsage {
		fmt.Printf("Interface: %s, RX: %d bytes, TX: %d bytes\n",
			bw.Interface, bw.RXRate, bw.TXRate)
	}

	fmt.Println("\n=== NETWORK ANALYSIS COMPLETED ===")
	fmt.Printf("Report completed at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
