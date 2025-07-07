#!/bin/bash

echo "=== SECURITY CHECK SCRIPT ==="
echo "Date: $(date)"
echo ""

# Check running processes
echo "=== RUNNING PROCESSES ==="
echo "Total processes: $(ps aux | wc -l)"
echo "Processes by user:"
ps aux | awk '{print $1}' | sort | uniq -c | sort -nr | head -5
echo ""

# Check listening ports
echo "=== LISTENING PORTS ==="
if command -v netstat &> /dev/null; then
    netstat -tlnp 2>/dev/null | head -10
else
    echo "netstat not available"
fi
echo ""

# Check file permissions
echo "=== CRITICAL FILE PERMISSIONS ==="
echo "Checking /etc/passwd permissions:"
ls -la /etc/passwd 2>/dev/null || echo "File not accessible"
echo ""

# Check for suspicious files
echo "=== SUSPICIOUS FILES CHECK ==="
echo "Looking for hidden files in /tmp:"
find /tmp -name ".*" -type f 2>/dev/null | head -5
echo ""

# Check system load
echo "=== SYSTEM LOAD ==="
uptime
echo ""

# Check memory usage
echo "=== MEMORY USAGE ==="
free -h
echo ""

echo "Security check completed!" 