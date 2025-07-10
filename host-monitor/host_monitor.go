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

type SystemInfo struct {
	Hostname    string
	Date        string
	CPUInfo     CPUInfo
	RAMInfo     RAMInfo
	DiskInfo    []DiskInfo
	NetworkInfo []NetworkInfo
	Temperature TemperatureInfo
}

type CPUInfo struct {
	ModelName string
	Sockets   string
	Threads   string
	Cores     string
	CPUs      string
	MHz       string
}

type RAMInfo struct {
	Total     string
	Used      string
	Free      string
	Available string
}

type DiskInfo struct {
	Filesystem string
	Size       string
	Used       string
	Available  string
	UsePercent string
	Mounted    string
}

type NetworkInfo struct {
	Interface string
	IP        string
	MAC       string
	RXBytes   int64
	TXBytes   int64
}

type TemperatureInfo struct {
	CPUTemp  string
	DiskTemp []string
}

type HostMonitor struct{}

func NewHostMonitor() *HostMonitor {
	return &HostMonitor{}
}

func (hm *HostMonitor) GetSystemInfo() SystemInfo {
	info := SystemInfo{
		Hostname: hm.getHostname(),
		Date:     time.Now().Format("2006-01-02 15:04:05"),
	}

	info.CPUInfo = hm.getCPUInfo()
	info.RAMInfo = hm.getRAMInfo()
	info.DiskInfo = hm.getDiskInfo()
	info.NetworkInfo = hm.getNetworkInfo()
	info.Temperature = hm.getTemperatureInfo()

	return info
}

func (hm *HostMonitor) getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func (hm *HostMonitor) getCPUInfo() CPUInfo {
	info := CPUInfo{}

	// Read /proc/cpuinfo
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return info
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			info.ModelName = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	// Get CPU count
	info.CPUs = hm.getCPUCount()
	info.Cores = hm.getCPUCores()
	info.Threads = hm.getCPUThreads()

	return info
}

func (hm *HostMonitor) getCPUCount() string {
	cmd := exec.Command("nproc")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (hm *HostMonitor) getCPUCores() string {
	cmd := exec.Command("nproc", "--all")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (hm *HostMonitor) getCPUThreads() string {
	cmd := exec.Command("nproc")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (hm *HostMonitor) getRAMInfo() RAMInfo {
	info := RAMInfo{}

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return info
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				info.Total = parts[1] + " KB"
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				info.Available = parts[1] + " KB"
			}
		} else if strings.HasPrefix(line, "MemFree:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				info.Free = parts[1] + " KB"
			}
		}
	}

	// Calculate used memory
	if info.Total != "" && info.Available != "" {
		total, _ := strconv.ParseInt(strings.Fields(info.Total)[0], 10, 64)
		available, _ := strconv.ParseInt(strings.Fields(info.Available)[0], 10, 64)
		used := total - available
		info.Used = fmt.Sprintf("%d KB", used)
	}

	return info
}

func (hm *HostMonitor) getDiskInfo() []DiskInfo {
	var disks []DiskInfo

	cmd := exec.Command("df", "-hT")
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
		if len(fields) >= 7 && !strings.Contains(fields[6], "tmpfs") {
			disk := DiskInfo{
				Filesystem: fields[0],
				Size:       fields[1],
				Used:       fields[2],
				Available:  fields[3],
				UsePercent: fields[4],
				Mounted:    fields[6],
			}
			disks = append(disks, disk)
		}
	}

	return disks
}

func (hm *HostMonitor) getNetworkInfo() []NetworkInfo {
	var networks []NetworkInfo

	// Get network interfaces
	interfaces, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return networks
	}

	for _, iface := range interfaces {
		if iface.Name() == "lo" {
			continue
		}

		network := NetworkInfo{Interface: iface.Name()}

		// Get IP address
		network.IP = hm.getInterfaceIP(iface.Name())

		// Get MAC address
		network.MAC = hm.getInterfaceMAC(iface.Name())

		// Get RX/TX bytes
		network.RXBytes = hm.getInterfaceRXBytes(iface.Name())
		network.TXBytes = hm.getInterfaceTXBytes(iface.Name())

		networks = append(networks, network)
	}

	return networks
}

func (hm *HostMonitor) getInterfaceIP(iface string) string {
	cmd := exec.Command("ip", "addr", "show", iface)
	output, err := cmd.Output()
	if err != nil {
		return "N/A"
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "inet ") {
			fields := strings.Fields(line)
			for _, field := range fields {
				if strings.Contains(field, ".") {
					return strings.Split(field, "/")[0]
				}
			}
		}
	}

	return "N/A"
}

func (hm *HostMonitor) getInterfaceMAC(iface string) string {
	data, err := os.ReadFile(fmt.Sprintf("/sys/class/net/%s/address", iface))
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(data))
}

func (hm *HostMonitor) getInterfaceRXBytes(iface string) int64 {
	data, err := os.ReadFile(fmt.Sprintf("/sys/class/net/%s/statistics/rx_bytes", iface))
	if err != nil {
		return 0
	}
	bytes, _ := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	return bytes
}

func (hm *HostMonitor) getInterfaceTXBytes(iface string) int64 {
	data, err := os.ReadFile(fmt.Sprintf("/sys/class/net/%s/statistics/tx_bytes", iface))
	if err != nil {
		return 0
	}
	bytes, _ := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	return bytes
}

func (hm *HostMonitor) getTemperatureInfo() TemperatureInfo {
	info := TemperatureInfo{}

	// Try to get CPU temperature using sensors
	cmd := exec.Command("sensors")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Core") || strings.Contains(line, "temp") {
				info.CPUTemp = strings.TrimSpace(line)
				break
			}
		}
	}

	// Try to get disk temperature using hddtemp
	cmd = exec.Command("hddtemp", "/dev/sda")
	output, err = cmd.Output()
	if err == nil {
		info.DiskTemp = append(info.DiskTemp, strings.TrimSpace(string(output)))
	}

	return info
}

func (hm *HostMonitor) PrintSystemReport() {
	info := hm.GetSystemInfo()

	fmt.Println("=== SYSTEM MONITORING REPORT ===")
	fmt.Printf("Date: %s\n", info.Date)
	fmt.Printf("Hostname: %s\n", info.Hostname)
	fmt.Println("==========================================")

	// Hardware Information
	fmt.Println("\n1. HARDWARE INFORMATION")
	fmt.Println("-----------------------")

	fmt.Println("CPU Info:")
	fmt.Printf("  Model: %s\n", info.CPUInfo.ModelName)
	fmt.Printf("  CPUs: %s\n", info.CPUInfo.CPUs)
	fmt.Printf("  Cores: %s\n", info.CPUInfo.Cores)
	fmt.Printf("  Threads: %s\n", info.CPUInfo.Threads)

	fmt.Println("\nRAM Info:")
	fmt.Printf("  Total: %s\n", info.RAMInfo.Total)
	fmt.Printf("  Used: %s\n", info.RAMInfo.Used)
	fmt.Printf("  Free: %s\n", info.RAMInfo.Free)
	fmt.Printf("  Available: %s\n", info.RAMInfo.Available)

	fmt.Println("\nDisk Usage:")
	for _, disk := range info.DiskInfo {
		fmt.Printf("  %s: %s/%s (%s) mounted on %s\n",
			disk.Filesystem, disk.Used, disk.Size, disk.UsePercent, disk.Mounted)
	}

	fmt.Println("\nNetwork Interfaces:")
	for _, network := range info.NetworkInfo {
		fmt.Printf("  %s: IP=%s, MAC=%s\n", network.Interface, network.IP, network.MAC)
		fmt.Printf("    RX: %d bytes, TX: %d bytes\n", network.RXBytes, network.TXBytes)
	}

	fmt.Println("\nCPU Temperature:")
	if info.Temperature.CPUTemp != "" {
		fmt.Printf("  %s\n", info.Temperature.CPUTemp)
	} else {
		fmt.Println("  No temperature info available")
	}

	fmt.Println("\nDisk Temperature:")
	if len(info.Temperature.DiskTemp) > 0 {
		for _, temp := range info.Temperature.DiskTemp {
			fmt.Printf("  %s\n", temp)
		}
	} else {
		fmt.Println("  No disk temperature info available")
	}

	// Network Traffic Report
	fmt.Println("\n2. NETWORK TRAFFIC REPORT")
	fmt.Println("--------------------------")

	fmt.Println("Current network usage (bytes):")
	for _, network := range info.NetworkInfo {
		fmt.Printf("  %s: RX=%d bytes, TX=%d bytes\n", network.Interface, network.RXBytes, network.TXBytes)
	}

	// Active connections
	cmd := exec.Command("ss", "-tunap")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		fmt.Printf("\nActive network connections: %d\n", len(lines)-1)
	}

	fmt.Println("\n=== SYSTEM MONITORING COMPLETED ===")
	fmt.Printf("Report completed at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
