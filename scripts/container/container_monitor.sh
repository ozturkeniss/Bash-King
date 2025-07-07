#!/bin/bash

echo "=== CONTAINER SYSTEM MONITORING ==="
echo "Date: $(date)"
echo "Hostname: $(hostname)"
echo "Container ID: $(hostname)"
echo "=========================================="

# 1. BASIC SYSTEM INFO (Container-friendly)
echo "\n1. BASIC SYSTEM INFORMATION"
echo "----------------------------"

# CPU Info (basic)
echo "CPU Info:"
lscpu | grep -E 'Model name|Socket|Thread|Core|CPU\(s\):|MHz' | uniq

# RAM Info
echo "\nRAM Info:"
free -h

# Disk Info (container filesystem)
echo "\nDisk Usage (Container):"
df -h

# Network Interfaces (container can see)
echo "\nNetwork Interfaces:"
for iface in $(ls /sys/class/net | grep -v lo); do
    ip_addr=$(ip addr show $iface 2>/dev/null | grep 'inet ' | awk '{print $2}')
    mac_addr=$(cat /sys/class/net/$iface/address 2>/dev/null)
    echo "$iface: IP=$ip_addr, MAC=$mac_addr"
done

# 2. CONTAINER-SPECIFIC INFO
echo "\n2. CONTAINER INFORMATION"
echo "------------------------"

# Container environment
echo "Environment variables (first 10):"
env | head -10

# Process count
echo "\nProcess count:"
ps aux | wc -l

# Top processes by CPU
echo "\nTop 5 processes by CPU:"
ps aux --sort=-%cpu | head -6

# Top processes by memory
echo "\nTop 5 processes by memory:"
ps aux --sort=-%mem | head -6

# 3. NETWORK INFO (Container-friendly)
echo "\n3. NETWORK INFORMATION"
echo "----------------------"

# Current network usage (per interface)
echo "Network usage (bytes):"
for iface in $(ls /sys/class/net | grep -v lo); do
    if [ -f /sys/class/net/$iface/statistics/rx_bytes ]; then
        rx=$(cat /sys/class/net/$iface/statistics/rx_bytes)
        tx=$(cat /sys/class/net/$iface/statistics/tx_bytes)
        echo "$iface: RX=$rx bytes, TX=$tx bytes"
    fi
done

# Active connections
echo "\nActive network connections:"
ss -tunap | wc -l

# Listening ports
echo "\nListening ports:"
ss -tlnp

echo "\n=== CONTAINER MONITORING COMPLETED ==="
echo "Report completed at: $(date)" 