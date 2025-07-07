#!/bin/bash

echo "=== CONTAINER SECURITY SCAN ==="
echo "Date: $(date)"
echo "Hostname: $(hostname)"
echo "Container ID: $(hostname)"
echo "=========================================="

# 1. CONTAINER-SPECIFIC SECURITY CHECKS
echo "\n1. CONTAINER SECURITY CHECKS"
echo "-----------------------------"

# Check if running as root
echo "Running as root:"
if [ "$(id -u)" -eq 0 ]; then
    echo "⚠️  WARNING: Running as root"
else
    echo "✅ Running as non-root user: $(whoami)"
fi

# Check for setuid binaries in container
echo "\nSetuid binaries in container:"
find /usr/bin /usr/sbin /bin /sbin -type f -perm -4000 2>/dev/null | head -10

# Check for world-writable files
echo "\nWorld-writable files (first 10):"
find / -type f -perm /o+w 2>/dev/null | head -10

# 2. NETWORK SECURITY (Container-friendly)
echo "\n2. NETWORK SECURITY"
echo "-------------------"

# Listening ports
echo "Listening ports:"
ss -tlnp

# Active connections
echo "\nActive connections count:"
ss -tunap | wc -l

# Network interfaces
echo "\nNetwork interfaces:"
ip addr show

# 3. PROCESS SECURITY
echo "\n3. PROCESS SECURITY"
echo "-------------------"

# Running processes
echo "Process count:"
ps aux | wc -l

# Processes with network connections
echo "\nProcesses with network connections:"
ss -tunap | grep -v "127.0.0.1" | head -10

# 4. FILE SYSTEM SECURITY
echo "\n4. FILE SYSTEM SECURITY"
echo "-----------------------"

# Recently modified files (last 24 hours)
echo "Recently modified files (last 24h):"
find / -type f -mtime -1 2>/dev/null | head -10

# Files with unusual permissions
echo "\nFiles with unusual permissions:"
find / -type f -perm /o+w 2>/dev/null | head -10

# 5. ENVIRONMENT SECURITY
echo "\n5. ENVIRONMENT SECURITY"
echo "-----------------------"

# Environment variables (security-related)
echo "Security-related environment variables:"
env | grep -i -E 'pass|key|secret|token|auth' | head -10

# Mounted volumes
echo "\nMounted volumes:"
mount | grep -v proc | grep -v sysfs | head -10

echo "\n=== CONTAINER SECURITY SCAN COMPLETED ==="
echo "Scan completed at: $(date)" 