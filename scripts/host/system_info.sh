#!/bin/bash

echo "=== SYSTEM INFORMATION ==="
echo "Hostname: $(hostname)"
echo "Uptime: $(uptime -p)"
echo ""

echo "=== CPU USAGE ==="
echo "CPU Load: $(uptime | awk -F'load average:' '{print $2}')"
echo "CPU Cores: $(nproc)"
echo ""

echo "=== MEMORY USAGE ==="
free -h | grep -E "Mem|Swap"
echo ""

echo "=== DISK USAGE ==="
df -h / | tail -1
echo ""

echo "=== PROCESS COUNT ==="
echo "Total Processes: $(ps aux | wc -l)"
echo ""

echo "=== NETWORK INTERFACES ==="
ip addr show | grep -E "inet.*scope global" | awk '{print $2}'
echo ""

echo "=== DOCKER CONTAINERS ==="
if command -v docker &> /dev/null; then
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
else
    echo "Docker not available"
fi
echo ""

echo "=== AGENT STATUS ==="
if pgrep -f "agent" > /dev/null; then
    echo "Agent Process: RUNNING"
    echo "Agent PID: $(pgrep -f agent)"
else
    echo "Agent Process: NOT RUNNING"
fi 