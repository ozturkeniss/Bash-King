#!/bin/bash

echo "=== SYSTEM MONITORING REPORT ==="
echo "Date: $(date)"
echo "Hostname: $(hostname)"
echo "=========================================="

# 1. HARDWARE INFO
echo "\n1. HARDWARE INFORMATION"
echo "-----------------------"

# CPU Info
echo "CPU Info:"
lscpu | grep -E 'Model name|Socket|Thread|Core|CPU\(s\):|MHz' | uniq

# RAM Info
echo "\nRAM Info:"
free -h

# Disk Info
echo "\nDisk Usage:"
df -hT | grep -v tmpfs

# Network Interfaces
echo "\nNetwork Interfaces:"
for iface in $(ls /sys/class/net | grep -v lo); do
    ip_addr=$(ip addr show $iface | grep 'inet ' | awk '{print $2}')
    mac_addr=$(cat /sys/class/net/$iface/address)
    echo "$iface: IP=$ip_addr, MAC=$mac_addr"
done

# CPU Temperature (if available)
echo "\nCPU Temperature:"
if command -v sensors &> /dev/null; then
    sensors | grep -E 'Core|temp' || echo "No temperature info from sensors."
else
    echo "sensors command not found. Install lm-sensors for temperature info."
fi

# Disk Temperature (if available)
echo "\nDisk Temperature:"
if command -v hddtemp &> /dev/null; then
    for disk in $(lsblk -dno NAME | grep -E 'sd|nvme'); do
        sudo hddtemp /dev/$disk 2>/dev/null
    done
else
    echo "hddtemp command not found. Install hddtemp for disk temperature info."
fi

# 2. NETWORK TRAFFIC REPORT
echo "\n2. NETWORK TRAFFIC REPORT"
echo "--------------------------"

# Top 5 IPs by traffic (last 5 min, if nethogs or iftop available)
if command -v nethogs &> /dev/null; then
    echo "Top 5 IPs by traffic (nethogs, last 5 min):"
    echo "(nethogs output requires root, skipping in script mode)"
elif command -v iftop &> /dev/null; then
    echo "Top 5 connections by traffic (iftop, last 5 min):"
    echo "(iftop output requires root, skipping in script mode)"
else
    echo "Neither nethogs nor iftop found. Showing current RX/TX per interface:"
    for iface in $(ls /sys/class/net | grep -v lo); do
        rx=$(cat /sys/class/net/$iface/statistics/rx_bytes)
        tx=$(cat /sys/class/net/$iface/statistics/tx_bytes)
        echo "$iface: RX=$((rx/1024/1024))MB, TX=$((tx/1024/1024))MB"
    done
fi

# Current network usage (per interface)
echo "\nCurrent network usage (bytes):"
for iface in $(ls /sys/class/net | grep -v lo); do
    rx=$(cat /sys/class/net/$iface/statistics/rx_bytes)
    tx=$(cat /sys/class/net/$iface/statistics/tx_bytes)
    echo "$iface: RX=$rx bytes, TX=$tx bytes"
done

# Active connections
echo "\nActive network connections:"
ss -tunap | wc -l

# Top 5 remote IPs by connection count
echo "\nTop 5 remote IPs by connection count:"
ss -tun | awk 'NR>1{print $5}' | cut -d: -f1 | sort | uniq -c | sort -nr | head -5

echo "\n=== SYSTEM MONITORING COMPLETED ==="
echo "Report completed at: $(date)" 