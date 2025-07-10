package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("🎯 HOST MONITORING SYSTEM")
	fmt.Println("Available monitoring modules:")
	fmt.Println("  1. System Monitor")
	fmt.Println("  2. Security Scanner")
	fmt.Println("  3. Performance Monitor")
	fmt.Println("  4. Network Analyzer")
	fmt.Println("  5. Package Manager")
	fmt.Println("  6. All Modules")
	fmt.Println("  7. Exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("💻 Select module (1-7): ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "1":
			fmt.Println("\n" + strings.Repeat("=", 50))
			monitor := NewHostMonitor()
			monitor.PrintSystemReport()
			fmt.Println(strings.Repeat("=", 50) + "\n")

		case "2":
			fmt.Println("\n" + strings.Repeat("=", 50))
			scanner := NewSecurityScanner()
			scanner.PrintSecurityReport()
			fmt.Println(strings.Repeat("=", 50) + "\n")

		case "3":
			fmt.Println("\n" + strings.Repeat("=", 50))
			perfMonitor := NewPerformanceMonitor()
			perfMonitor.PrintPerformanceReport()
			fmt.Println(strings.Repeat("=", 50) + "\n")

		case "4":
			fmt.Println("\n" + strings.Repeat("=", 50))
			netAnalyzer := NewNetworkAnalyzer()
			netAnalyzer.PrintNetworkReport()
			fmt.Println(strings.Repeat("=", 50) + "\n")

		case "5":
			fmt.Println("\n" + strings.Repeat("=", 50))
			pkgManager := NewPackageManager()
			pkgManager.PrintPackageReport()
			fmt.Println(strings.Repeat("=", 50) + "\n")

		case "6":
			fmt.Println("\n" + strings.Repeat("=", 50))
			fmt.Println("RUNNING ALL MODULES")
			fmt.Println(strings.Repeat("=", 50))

			// System Monitor
			fmt.Println("\n1. SYSTEM MONITOR")
			monitor := NewHostMonitor()
			monitor.PrintSystemReport()

			// Security Scanner
			fmt.Println("\n2. SECURITY SCANNER")
			scanner := NewSecurityScanner()
			scanner.PrintSecurityReport()

			// Performance Monitor
			fmt.Println("\n3. PERFORMANCE MONITOR")
			perfMonitor := NewPerformanceMonitor()
			perfMonitor.PrintPerformanceReport()

			// Network Analyzer
			fmt.Println("\n4. NETWORK ANALYZER")
			netAnalyzer := NewNetworkAnalyzer()
			netAnalyzer.PrintNetworkReport()

			// Package Manager
			fmt.Println("\n5. PACKAGE MANAGER")
			pkgManager := NewPackageManager()
			pkgManager.PrintPackageReport()

			fmt.Println("\n" + strings.Repeat("=", 50))
			fmt.Println("ALL MODULES COMPLETED")
			fmt.Println(strings.Repeat("=", 50) + "\n")

		case "7":
			fmt.Println("👋 Goodbye!")
			return

		default:
			fmt.Println("❌ Invalid option. Please select 1-7.")
		}
	}
}
