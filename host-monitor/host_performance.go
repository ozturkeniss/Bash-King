package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type PerformanceInfo struct {
	Date         string
	Hostname     string
	CPUUsage     CPUUsage
	MemoryUsage  MemoryUsage
	DiskUsage    []DiskUsage
	NetworkUsage []NetworkUsage
	ProcessInfo  []ProcessInfo
	LoadAverage  LoadAverage
	Uptime       string
	Temperature  string
}

type CPUUsage struct {
	User   float64
	System float64
	Idle   float64
	IOWait float64
}

type MemoryUsage struct {
	Total     int64
	Used      int64
	Free      int64
	Available int64
	SwapTotal int64
	SwapUsed  int64
}

type DiskUsage struct {
	Device     string
	MountPoint string
	Total      int64
	Used       int64
	Available  int64
	UsePercent float64
}

type NetworkUsage struct {
	Interface string
	RXBytes   int64
	TXBytes   int64
	RXPackets int64
	TXPackets int64
}

type LoadAverage struct {
	OneMin     float64
	FiveMin    float64
	FifteenMin float64
}

type PerformanceMonitor struct{}

func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{}
}

func (pm *PerformanceMonitor) GetPerformanceInfo() PerformanceInfo {
	info := PerformanceInfo{
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Hostname: pm.getHostname(),
	}

	info.CPUUsage = pm.getCPUUsage()
	info.MemoryUsage = pm.getMemoryUsage()
	info.DiskUsage = pm.getDiskUsage()
	info.NetworkUsage = pm.getNetworkUsage()
	info.ProcessInfo = pm.getTopProcesses()
	info.LoadAverage = pm.getLoadAverage()
	info.Uptime = pm.getUptime()
	info.Temperature = pm.getTemperature()

	return info
}

func (pm *PerformanceMonitor) getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func (pm *PerformanceMonitor) getCPUUsage() CPUUsage {
	usage := CPUUsage{}

	file, err := os.Open("/proc/stat")
	if err != nil {
		return usage
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				user, _ := strconv.ParseFloat(fields[1], 64)
				system, _ := strconv.ParseFloat(fields[3], 64)
				idle, _ := strconv.ParseFloat(fields[4], 64)
				iowait, _ := strconv.ParseFloat(fields[5], 64)

				total := user + system + idle + iowait
				if total > 0 {
					usage.User = (user / total) * 100
					usage.System = (system / total) * 100
					usage.Idle = (idle / total) * 100
					usage.IOWait = (iowait / total) * 100
				}
			}
		}
	}

	return usage
}

func (pm *PerformanceMonitor) getMemoryUsage() MemoryUsage {
	usage := MemoryUsage{}

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return usage
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				usage.Total, _ = strconv.ParseInt(parts[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				usage.Available, _ = strconv.ParseInt(parts[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "MemFree:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				usage.Free, _ = strconv.ParseInt(parts[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "SwapTotal:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				usage.SwapTotal, _ = strconv.ParseInt(parts[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "SwapFree:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				swapFree, _ := strconv.ParseInt(parts[1], 10, 64)
				usage.SwapUsed = usage.SwapTotal - swapFree
			}
		}
	}

	usage.Used = usage.Total - usage.Available

	return usage
}

func (pm *PerformanceMonitor) getDiskUsage() []DiskUsage {
	var disks []DiskUsage

	cmd := exec.Command("df", "-B1")
	output, err := cmd.Output()
	if err != nil {
		return disks
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 6 && !strings.Contains(fields[5], "tmpfs") {
			total, _ := strconv.ParseInt(fields[1], 10, 64)
			used, _ := strconv.ParseInt(fields[2], 10, 64)
			available, _ := strconv.ParseInt(fields[3], 10, 64)
			usePercent, _ := strconv.ParseFloat(strings.TrimSuffix(fields[4], "%"), 64)

			disk := DiskUsage{
				Device:     fields[0],
				MountPoint: fields[5],
				Total:      total,
				Used:       used,
				Available:  available,
				UsePercent: usePercent,
			}
			disks = append(disks, disk)
		}
	}

	return disks
}

func (pm *PerformanceMonitor) getNetworkUsage() []NetworkUsage {
	var networks []NetworkUsage

	interfaces, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return networks
	}

	for _, iface := range interfaces {
		if iface.Name() == "lo" {
			continue
		}

		network := NetworkUsage{Interface: iface.Name()}

		// Get RX bytes
		network.RXBytes = pm.readNetworkStat(iface.Name(), "rx_bytes")

		// Get TX bytes
		network.TXBytes = pm.readNetworkStat(iface.Name(), "tx_bytes")

		// Get RX packets
		network.RXPackets = pm.readNetworkStat(iface.Name(), "rx_packets")

		// Get TX packets
		network.TXPackets = pm.readNetworkStat(iface.Name(), "tx_packets")

		networks = append(networks, network)
	}

	return networks
}

func (pm *PerformanceMonitor) readNetworkStat(iface, stat string) int64 {
	data, err := os.ReadFile(fmt.Sprintf("/sys/class/net/%s/statistics/%s", iface, stat))
	if err != nil {
		return 0
	}
	value, _ := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	return value
}

func (pm *PerformanceMonitor) getTopProcesses() []ProcessInfo {
	var processes []ProcessInfo

	cmd := exec.Command("ps", "aux", "--sort=-%cpu")
	output, err := cmd.Output()
	if err != nil {
		return processes
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "USER") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 11 {
			cpu, _ := strconv.ParseFloat(fields[2], 64)
			memory, _ := strconv.ParseFloat(fields[3], 64)

			process := ProcessInfo{
				PID:     fields[1],
				User:    fields[0],
				CPU:     cpu,
				Memory:  memory,
				Command: strings.Join(fields[10:], " "),
			}
			processes = append(processes, process)

			if len(processes) >= 10 {
				break
			}
		}
	}

	return processes
}

func (pm *PerformanceMonitor) getLoadAverage() LoadAverage {
	load := LoadAverage{}

	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return load
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			load.OneMin, _ = strconv.ParseFloat(fields[0], 64)
			load.FiveMin, _ = strconv.ParseFloat(fields[1], 64)
			load.FifteenMin, _ = strconv.ParseFloat(fields[2], 64)
		}
	}

	return load
}

func (pm *PerformanceMonitor) getUptime() string {
	file, err := os.Open("/proc/uptime")
	if err != nil {
		return "unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			uptime, _ := strconv.ParseFloat(fields[0], 64)
			days := int(uptime / 86400)
			hours := int((uptime - float64(days*86400)) / 3600)
			minutes := int((uptime - float64(days*86400) - float64(hours*3600)) / 60)
			return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
		}
	}

	return "unknown"
}

func (pm *PerformanceMonitor) getTemperature() string {
	cmd := exec.Command("sensors")
	output, err := cmd.Output()
	if err != nil {
		return "No temperature info available"
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Core") || strings.Contains(line, "temp") {
			return strings.TrimSpace(line)
		}
	}

	return "No temperature info available"
}

func (pm *PerformanceMonitor) PrintPerformanceReport() {
	info := pm.GetPerformanceInfo()

	fmt.Println("=== PERFORMANCE MONITORING REPORT ===")
	fmt.Printf("Date: %s\n", info.Date)
	fmt.Printf("Hostname: %s\n", info.Hostname)
	fmt.Printf("Uptime: %s\n", info.Uptime)
	fmt.Println("==========================================")

	// CPU Information
	fmt.Println("\n1. CPU USAGE")
	fmt.Println("-------------")
	fmt.Printf("User: %.2f%%\n", info.CPUUsage.User)
	fmt.Printf("System: %.2f%%\n", info.CPUUsage.System)
	fmt.Printf("Idle: %.2f%%\n", info.CPUUsage.Idle)
	fmt.Printf("I/O Wait: %.2f%%\n", info.CPUUsage.IOWait)

	// Load Average
	fmt.Printf("\nLoad Average: %.2f, %.2f, %.2f\n",
		info.LoadAverage.OneMin, info.LoadAverage.FiveMin, info.LoadAverage.FifteenMin)

	// Memory Information
	fmt.Println("\n2. MEMORY USAGE")
	fmt.Println("---------------")
	fmt.Printf("Total: %d KB\n", info.MemoryUsage.Total)
	fmt.Printf("Used: %d KB\n", info.MemoryUsage.Used)
	fmt.Printf("Free: %d KB\n", info.MemoryUsage.Free)
	fmt.Printf("Available: %d KB\n", info.MemoryUsage.Available)
	fmt.Printf("Swap Total: %d KB\n", info.MemoryUsage.SwapTotal)
	fmt.Printf("Swap Used: %d KB\n", info.MemoryUsage.SwapUsed)

	// Disk Usage
	fmt.Println("\n3. DISK USAGE")
	fmt.Println("-------------")
	for _, disk := range info.DiskUsage {
		fmt.Printf("%s (%s): %.1f%% used (%d/%d bytes)\n",
			disk.Device, disk.MountPoint, disk.UsePercent, disk.Used, disk.Total)
	}

	// Network Usage
	fmt.Println("\n4. NETWORK USAGE")
	fmt.Println("----------------")
	for _, network := range info.NetworkUsage {
		fmt.Printf("%s: RX=%d bytes (%d packets), TX=%d bytes (%d packets)\n",
			network.Interface, network.RXBytes, network.RXPackets, network.TXBytes, network.TXPackets)
	}

	// Top Processes
	fmt.Println("\n5. TOP PROCESSES BY CPU")
	fmt.Println("------------------------")
	for i, process := range info.ProcessInfo {
		fmt.Printf("%d. PID=%s, User=%s, CPU=%.1f%%, Memory=%.1f%%, Command=%s\n",
			i+1, process.PID, process.User, process.CPU, process.Memory, process.Command)
	}

	// Temperature
	fmt.Println("\n6. TEMPERATURE")
	fmt.Println("---------------")
	fmt.Printf("%s\n", info.Temperature)

	fmt.Println("\n=== PERFORMANCE MONITORING COMPLETED ===")
	fmt.Printf("Report completed at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
