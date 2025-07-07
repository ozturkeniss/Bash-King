#!/bin/bash

echo "=== HOST PERFORMANCE MONITORING ==="
echo "Timestamp: $(date)"
echo "Hostname: $(hostname)"
echo ""

echo "=== CPU PERFORMANCE ==="
echo "CPU Usage:"
top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1
echo ""

echo "CPU Load Average:"
uptime | awk -F'load average:' '{print $2}'
echo ""

echo "CPU Core Information:"
nproc
echo ""

echo "=== MEMORY PERFORMANCE ==="
echo "Memory Usage:"
free -h
echo ""

echo "Memory Details:"
cat /proc/meminfo | grep -E "(MemTotal|MemFree|MemAvailable|Buffers|Cached|SwapTotal|SwapFree)"
echo ""

echo "=== DISK PERFORMANCE ==="
echo "Disk Usage:"
df -h
echo ""

echo "Disk I/O Statistics:"
iostat -x 1 1 2>/dev/null || echo "iostat not available"
echo ""

echo "=== NETWORK PERFORMANCE ==="
echo "Network Interface Statistics:"
cat /proc/net/dev
echo ""

echo "Network Connections:"
ss -tuln | head -20
echo ""

echo "=== PROCESS PERFORMANCE ==="
echo "Top 10 CPU Usage Processes:"
ps aux --sort=-%cpu | head -11
echo ""

echo "Top 10 Memory Usage Processes:"
ps aux --sort=-%mem | head -11
echo ""

echo "=== SYSTEM UPTIME ==="
uptime
echo ""

echo "=== LAST SYSTEM REBOOT ==="
who -b
echo ""

echo "=== SYSTEM LOAD HISTORY ==="
echo "Load average over time:"
uptime | awk '{print "1 min: " $(NF-2) " 5 min: " $(NF-1) " 15 min: " $NF}'
echo ""

echo "=== PERFORMANCE SUMMARY ==="
echo "CPU Cores: $(nproc)"
echo "Total Memory: $(free -h | grep Mem | awk '{print $2}')"
echo "Available Memory: $(free -h | grep Mem | awk '{print $7}')"
echo "Disk Usage: $(df -h / | tail -1 | awk '{print $5}')"
echo "Uptime: $(uptime -p)"
echo ""

echo "=== END OF PERFORMANCE REPORT ===" 