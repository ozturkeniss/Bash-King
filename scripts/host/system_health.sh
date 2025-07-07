#!/bin/bash

echo "=== HOST SYSTEM HEALTH CHECK ==="
echo "Timestamp: $(date)"
echo "Hostname: $(hostname)"
echo ""

echo "=== SYSTEM STATUS ==="
echo "System Uptime:"
uptime
echo ""

echo "Last System Boot:"
who -b
echo ""

echo "=== MEMORY HEALTH ==="
echo "Memory Status:"
free -h
echo ""

echo "Memory Pressure:"
cat /proc/vmstat | grep -E "(pgpgin|pgpgout|pswpin|pswpout)" 2>/dev/null || echo "Memory pressure info not available"
echo ""

echo "=== DISK HEALTH ==="
echo "Disk Usage:"
df -h
echo ""

echo "Disk Health (if SMART available):"
if command -v smartctl >/dev/null 2>&1; then
    for disk in $(lsblk -d -n -o NAME | grep -v loop); do
        echo "Checking $disk:"
        smartctl -H /dev/$disk 2>/dev/null | grep -E "(SMART|Health)" || echo "SMART not available for $disk"
    done
else
    echo "smartctl not available"
fi
echo ""

echo "=== CPU HEALTH ==="
echo "CPU Temperature (if available):"
if [ -f /sys/class/thermal/thermal_zone0/temp ]; then
    temp=$(cat /sys/class/thermal/thermal_zone0/temp)
    echo "CPU Temperature: $((temp/1000))°C"
else
    echo "Temperature sensors not available"
fi
echo ""

echo "CPU Load:"
uptime | awk -F'load average:' '{print $2}'
echo ""

echo "=== PROCESS HEALTH ==="
echo "Zombie Processes:"
ps aux | grep -w Z | grep -v grep || echo "No zombie processes found"
echo ""

echo "High CPU Usage Processes:"
ps aux --sort=-%cpu | head -6
echo ""

echo "High Memory Usage Processes:"
ps aux --sort=-%mem | head -6
echo ""

echo "=== SYSTEM LOGS ==="
echo "Recent System Errors:"
journalctl -p err --since "1 hour ago" 2>/dev/null | tail -10 || echo "journalctl not available"
echo ""

echo "Recent Kernel Messages:"
dmesg | tail -10
echo ""

echo "=== SERVICE STATUS ==="
echo "Critical Services Status:"
systemctl is-active ssh 2>/dev/null && echo "SSH: ACTIVE" || echo "SSH: INACTIVE"
systemctl is-active network-manager 2>/dev/null && echo "Network Manager: ACTIVE" || echo "Network Manager: INACTIVE"
systemctl is-active cron 2>/dev/null && echo "Cron: ACTIVE" || echo "Cron: INACTIVE"
echo ""

echo "=== NETWORK HEALTH ==="
echo "Network Interface Status:"
ip link show | grep -E "UP|DOWN"
echo ""

echo "Network Connectivity:"
ping -c 1 8.8.8.8 >/dev/null 2>&1 && echo "Internet: CONNECTED" || echo "Internet: DISCONNECTED"
echo ""

echo "=== SECURITY HEALTH ==="
echo "Failed Login Attempts:"
grep "Failed password" /var/log/auth.log 2>/dev/null | tail -5 || echo "No failed login attempts found"
echo ""

echo "Sudo Usage:"
grep sudo /var/log/auth.log 2>/dev/null | tail -5 || echo "No sudo usage found"
echo ""

echo "=== FILE SYSTEM HEALTH ==="
echo "File System Errors:"
dmesg | grep -i "error\|fail" | tail -5
echo ""

echo "Inode Usage:"
df -i
echo ""

echo "=== SYSTEM RESOURCES ==="
echo "Open File Descriptors:"
lsof | wc -l
echo ""

echo "Active Network Connections:"
ss -tuln | wc -l
echo ""

echo "=== HEALTH SUMMARY ==="
echo "System Health Score:"

# Calculate health score
score=100

# Check memory usage
mem_usage=$(free | grep Mem | awk '{printf "%.0f", $3/$2 * 100}')
if [ $mem_usage -gt 90 ]; then
    score=$((score - 20))
    echo "WARNING: High memory usage ($mem_usage%)"
fi

# Check disk usage
disk_usage=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
if [ $disk_usage -gt 90 ]; then
    score=$((score - 20))
    echo "WARNING: High disk usage ($disk_usage%)"
fi

# Check load average
load=$(uptime | awk -F'load average:' '{print $2}' | awk '{print $1}' | sed 's/,//')
if (( $(echo "$load > 2" | bc -l) )); then
    score=$((score - 15))
    echo "WARNING: High system load ($load)"
fi

echo "Overall Health Score: $score/100"

if [ $score -ge 80 ]; then
    echo "STATUS: HEALTHY"
elif [ $score -ge 60 ]; then
    echo "STATUS: WARNING"
else
    echo "STATUS: CRITICAL"
fi
echo ""

echo "=== END OF HEALTH CHECK ===" 