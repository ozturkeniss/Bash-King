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

type PackageInfo struct {
	Date          string
	Hostname      string
	InstalledPkgs []Package
	AvailablePkgs []Package
	OutdatedPkgs  []Package
	SecurityPkgs  []Package
	PackageStats  PackageStats
	Repositories  []Repository
	UpdateHistory []UpdateRecord
}

type Package struct {
	Name         string
	Version      string
	Architecture string
	Size         string
	Description  string
	Status       string
	Priority     string
	Section      string
}

type PackageStats struct {
	TotalInstalled int
	TotalAvailable int
	TotalOutdated  int
	TotalSecurity  int
	TotalSize      int64
}

type Repository struct {
	Name     string
	URL      string
	Enabled  bool
	Priority int
}

type UpdateRecord struct {
	Package    string
	OldVersion string
	NewVersion string
	Date       string
}

type PackageManager struct{}

func NewPackageManager() *PackageManager {
	return &PackageManager{}
}

func (pm *PackageManager) GetPackageInfo() PackageInfo {
	info := PackageInfo{
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Hostname: pm.getHostname(),
	}

	info.InstalledPkgs = pm.getInstalledPackages()
	info.AvailablePkgs = pm.getAvailablePackages()
	info.OutdatedPkgs = pm.getOutdatedPackages()
	info.SecurityPkgs = pm.getSecurityPackages()
	info.PackageStats = pm.getPackageStats(info)
	info.Repositories = pm.getRepositories()
	info.UpdateHistory = pm.getUpdateHistory()

	return info
}

func (pm *PackageManager) getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func (pm *PackageManager) getInstalledPackages() []Package {
	var packages []Package

	cmd := exec.Command("dpkg", "-l")
	output, err := cmd.Output()
	if err != nil {
		return packages
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ii") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				pkg := Package{
					Name:         fields[1],
					Version:      fields[2],
					Architecture: fields[3],
					Description:  strings.Join(fields[4:], " "),
					Status:       "installed",
				}
				packages = append(packages, pkg)
			}
		}
	}

	return packages
}

func (pm *PackageManager) getAvailablePackages() []Package {
	var packages []Package

	cmd := exec.Command("apt", "list", "--upgradable")
	output, err := cmd.Output()
	if err != nil {
		return packages
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/") && !strings.HasPrefix(line, "WARNING") {
			parts := strings.Split(line, "/")
			if len(parts) >= 2 {
				nameParts := strings.Split(parts[0], " ")
				if len(nameParts) >= 1 {
					pkg := Package{
						Name:    nameParts[0],
						Status:  "available",
						Section: parts[1],
					}
					packages = append(packages, pkg)
				}
			}
		}
	}

	return packages
}

func (pm *PackageManager) getOutdatedPackages() []Package {
	var packages []Package

	cmd := exec.Command("apt", "list", "--upgradable")
	output, err := cmd.Output()
	if err != nil {
		return packages
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/") && !strings.HasPrefix(line, "WARNING") {
			parts := strings.Split(line, "/")
			if len(parts) >= 2 {
				nameParts := strings.Split(parts[0], " ")
				if len(nameParts) >= 1 {
					pkg := Package{
						Name:   nameParts[0],
						Status: "outdated",
					}
					packages = append(packages, pkg)
				}
			}
		}
	}

	return packages
}

func (pm *PackageManager) getSecurityPackages() []Package {
	var packages []Package

	cmd := exec.Command("apt", "list", "--upgradable")
	output, err := cmd.Output()
	if err != nil {
		return packages
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/") && !strings.HasPrefix(line, "WARNING") {
			parts := strings.Split(line, "/")
			if len(parts) >= 2 {
				nameParts := strings.Split(parts[0], " ")
				if len(nameParts) >= 1 {
					// Check if it's a security update
					if strings.Contains(strings.ToLower(line), "security") {
						pkg := Package{
							Name:   nameParts[0],
							Status: "security",
						}
						packages = append(packages, pkg)
					}
				}
			}
		}
	}

	return packages
}

func (pm *PackageManager) getPackageStats(info PackageInfo) PackageStats {
	stats := PackageStats{
		TotalInstalled: len(info.InstalledPkgs),
		TotalAvailable: len(info.AvailablePkgs),
		TotalOutdated:  len(info.OutdatedPkgs),
		TotalSecurity:  len(info.SecurityPkgs),
	}

	// Calculate total size
	cmd := exec.Command("du", "-sh", "/var/cache/apt/archives")
	output, err := cmd.Output()
	if err == nil {
		stats.TotalSize = pm.parseSize(strings.TrimSpace(string(output)))
	}

	return stats
}

func (pm *PackageManager) parseSize(sizeStr string) int64 {
	// Remove non-numeric characters and convert to bytes
	sizeStr = strings.TrimSpace(sizeStr)
	if strings.HasSuffix(sizeStr, "K") {
		sizeStr = strings.TrimSuffix(sizeStr, "K")
		if val, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
			return val * 1024
		}
	} else if strings.HasSuffix(sizeStr, "M") {
		sizeStr = strings.TrimSuffix(sizeStr, "M")
		if val, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
			return val * 1024 * 1024
		}
	} else if strings.HasSuffix(sizeStr, "G") {
		sizeStr = strings.TrimSuffix(sizeStr, "G")
		if val, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
			return val * 1024 * 1024 * 1024
		}
	}
	return 0
}

func (pm *PackageManager) getRepositories() []Repository {
	var repos []Repository

	// Read sources.list
	file, err := os.Open("/etc/apt/sources.list")
	if err != nil {
		return repos
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				repo := Repository{
					Name:    fields[2],
					URL:     fields[1],
					Enabled: true,
				}
				repos = append(repos, repo)
			}
		}
	}

	// Read sources.list.d
	dir, err := os.ReadDir("/etc/apt/sources.list.d")
	if err == nil {
		for _, file := range dir {
			if strings.HasSuffix(file.Name(), ".list") {
				filePath := fmt.Sprintf("/etc/apt/sources.list.d/%s", file.Name())
				content, err := os.ReadFile(filePath)
				if err == nil {
					lines := strings.Split(string(content), "\n")
					for _, line := range lines {
						if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
							fields := strings.Fields(line)
							if len(fields) >= 3 {
								repo := Repository{
									Name:    fields[2],
									URL:     fields[1],
									Enabled: true,
								}
								repos = append(repos, repo)
							}
						}
					}
				}
			}
		}
	}

	return repos
}

func (pm *PackageManager) getUpdateHistory() []UpdateRecord {
	var history []UpdateRecord

	// Read dpkg log
	file, err := os.Open("/var/log/dpkg.log")
	if err != nil {
		return history
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "upgrade") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				record := UpdateRecord{
					Package: fields[3],
					Date:    fields[0] + " " + fields[1],
				}
				history = append(history, record)
			}
		}
	}

	return history
}

func (pm *PackageManager) PrintPackageReport() {
	info := pm.GetPackageInfo()

	fmt.Println("=== PACKAGE MANAGER REPORT ===")
	fmt.Printf("Date: %s\n", info.Date)
	fmt.Printf("Hostname: %s\n", info.Hostname)
	fmt.Println("==========================================")

	// Package Statistics
	fmt.Println("\n1. PACKAGE STATISTICS")
	fmt.Println("----------------------")
	fmt.Printf("Total installed packages: %d\n", info.PackageStats.TotalInstalled)
	fmt.Printf("Available updates: %d\n", info.PackageStats.TotalAvailable)
	fmt.Printf("Outdated packages: %d\n", info.PackageStats.TotalOutdated)
	fmt.Printf("Security updates: %d\n", info.PackageStats.TotalSecurity)
	fmt.Printf("Cache size: %d bytes\n", info.PackageStats.TotalSize)

	// Installed Packages
	fmt.Println("\n2. INSTALLED PACKAGES (top 20)")
	fmt.Println("--------------------------------")
	for i, pkg := range info.InstalledPkgs {
		if i >= 20 {
			break
		}
		fmt.Printf("%s (%s) - %s\n", pkg.Name, pkg.Version, pkg.Description)
	}

	// Available Updates
	fmt.Println("\n3. AVAILABLE UPDATES")
	fmt.Println("--------------------")
	for _, pkg := range info.AvailablePkgs {
		fmt.Printf("%s (%s)\n", pkg.Name, pkg.Section)
	}

	// Outdated Packages
	fmt.Println("\n4. OUTDATED PACKAGES")
	fmt.Println("--------------------")
	for _, pkg := range info.OutdatedPkgs {
		fmt.Printf("%s\n", pkg.Name)
	}

	// Security Updates
	fmt.Println("\n5. SECURITY UPDATES")
	fmt.Println("-------------------")
	for _, pkg := range info.SecurityPkgs {
		fmt.Printf("%s\n", pkg.Name)
	}

	// Repositories
	fmt.Println("\n6. REPOSITORIES")
	fmt.Println("---------------")
	for _, repo := range info.Repositories {
		fmt.Printf("%s: %s (enabled: %t)\n", repo.Name, repo.URL, repo.Enabled)
	}

	// Update History
	fmt.Println("\n7. RECENT UPDATE HISTORY")
	fmt.Println("------------------------")
	for i, record := range info.UpdateHistory {
		if i >= 10 {
			break
		}
		fmt.Printf("%s: %s\n", record.Date, record.Package)
	}

	fmt.Println("\n=== PACKAGE MANAGER REPORT COMPLETED ===")
	fmt.Printf("Report completed at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
