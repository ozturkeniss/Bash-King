# Host Monitoring System

A comprehensive Go-based host monitoring system that provides detailed system information, security scanning, performance monitoring, network analysis, and package management capabilities.

## Features

### 🔍 System Monitor
- Hardware information (CPU, RAM, Disk)
- Network interface details
- Temperature monitoring
- System uptime and load

### 🛡️ Security Scanner
- Open port analysis
- Suspicious file detection
- High CPU process monitoring
- User and sudo activity logs
- System integrity checks
- Firewall status

### 📊 Performance Monitor
- Real-time CPU usage
- Memory utilization
- Disk usage statistics
- Network bandwidth
- Top processes by CPU/Memory
- Load average monitoring

### 🌐 Network Analyzer
- Network interface details
- Routing table information
- DNS configuration
- Active connections
- Firewall rules
- Bandwidth usage

### 📦 Package Manager
- Installed packages list
- Available updates
- Security updates
- Repository information
- Update history

## Installation

### Prerequisites
- Go 1.19 or higher
- Linux system with standard monitoring tools
- Root access for some features

### Build

```bash
# Navigate to host-monitor directory
cd host-monitor

# Build the system
make build

# Or manually
go build -o host-monitor main.go host_monitor.go host_security.go host_performance.go host_network.go host_package.go
```

## Usage

### Interactive Mode
```bash
./host-monitor
```

This will show a menu with options:
1. System Monitor
2. Security Scanner
3. Performance Monitor
4. Network Analyzer
5. Package Manager
6. All Modules
7. Exit

### Direct Module Execution
```bash
# Run all modules at once
make all

# Or manually
echo "6" | ./host-monitor
```

### Makefile Commands
```bash
make build    # Build the system
make clean    # Clean build artifacts
make run      # Build and run interactively
make all      # Build and run all modules
make test     # Test the build
make help     # Show available commands
```

## Module Details

### System Monitor (`host_monitor.go`)
- **CPU Info**: Model, cores, threads, usage
- **RAM Info**: Total, used, free, available memory
- **Disk Usage**: Filesystem information and usage
- **Network**: Interface details and traffic statistics
- **Temperature**: CPU and disk temperature monitoring

### Security Scanner (`host_security.go`)
- **Port Analysis**: TCP/UDP listening ports
- **Process Monitoring**: High CPU usage detection
- **File Integrity**: Suspicious file detection
- **User Activity**: Login and sudo logs
- **System Checks**: Modified files, unusual permissions

### Performance Monitor (`host_performance.go`)
- **CPU Usage**: User, system, idle, I/O wait percentages
- **Memory**: Detailed memory and swap usage
- **Disk**: Usage statistics and I/O information
- **Network**: Bandwidth and packet statistics
- **Processes**: Top processes by resource usage

### Network Analyzer (`host_network.go`)
- **Interfaces**: IP, MAC, status, statistics
- **Routing**: Routing table and gateway information
- **DNS**: Nameservers and search domains
- **Connections**: Active network connections
- **Firewall**: UFW and iptables rules

### Package Manager (`host_package.go`)
- **Packages**: Installed package information
- **Updates**: Available and security updates
- **Repositories**: APT repository configuration
- **History**: Recent package update history

## System Requirements

### Required Tools
- `ss` - Socket statistics
- `ps` - Process information
- `df` - Disk usage
- `ip` - Network configuration
- `dpkg` - Package management
- `apt` - Package updates
- `sensors` - Temperature monitoring (optional)
- `hddtemp` - Disk temperature (optional)

### File Access
- `/proc/*` - System information
- `/sys/class/net/*` - Network statistics
- `/etc/resolv.conf` - DNS configuration
- `/var/log/*` - System logs
- `/etc/apt/sources.list*` - Package repositories

## Output Examples

### System Monitor Output
```
=== SYSTEM MONITORING REPORT ===
Date: 2024-01-15 10:30:45
Hostname: server01
==========================================

1. HARDWARE INFORMATION
-----------------------
CPU Info:
  Model: Intel(R) Core(TM) i7-8700K CPU @ 3.70GHz
  CPUs: 12
  Cores: 6
  Threads: 12

RAM Info:
  Total: 16384000 KB
  Used: 8192000 KB
  Free: 4096000 KB
  Available: 12288000 KB
```

### Security Scanner Output
```
=== SECURITY SCAN REPORT ===
Date: 2024-01-15 10:30:45
Hostname: server01
==========================================

1. OPEN PORTS ANALYSIS
----------------------
TCP Listening Ports:
  TCP:22 - sshd
  TCP:80 - nginx
  TCP:443 - nginx
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add your monitoring module
4. Update the main.go to include your module
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Troubleshooting

### Common Issues

1. **Permission Denied**: Some features require root access
   ```bash
   sudo ./host-monitor
   ```

2. **Missing Tools**: Install required packages
   ```bash
   sudo apt update
   sudo apt install procps iproute2 net-tools lsof
   ```

3. **Build Errors**: Ensure Go version is 1.19+
   ```bash
   go version
   ```

4. **Temperature Not Available**: Install lm-sensors
   ```bash
   sudo apt install lm-sensors
   sudo sensors-detect
   ```

### Debug Mode
For debugging, you can run individual modules directly:
```bash
go run host_monitor.go
go run host_security.go
go run host_performance.go
go run host_network.go
go run host_package.go
``` 