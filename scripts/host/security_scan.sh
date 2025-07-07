#!/bin/bash

echo "=== SECURITY SCAN REPORT ==="
echo "Date: $(date)"
echo "Hostname: $(hostname)"
echo "=========================================="

# 1. OPEN PORT SCAN
echo "1. OPEN PORTS ANALYSIS"
echo "----------------------"
echo "TCP Listening Ports:"
if command -v ss &> /dev/null; then
    ss -tlnp | grep LISTEN
else
    netstat -tlnp | grep LISTEN
fi

echo ""
echo "UDP Listening Ports:"
if command -v ss &> /dev/null; then
    ss -ulnp | grep UNCONN
else
    netstat -ulnp | grep UNCONN
fi

echo ""
echo "Ports by process:"
lsof -i -P -n | grep LISTEN

# 2. MALICIOUS SOFTWARE SCAN
echo ""
echo "2. MALICIOUS SOFTWARE DETECTION"
echo "-------------------------------"

# Check for common rootkit tools
echo "Checking for common rootkit indicators:"
suspicious_files=(
    "/tmp/.X11-unix"
    "/dev/.udev"
    "/dev/.initramfs"
    "/lib/libkeyutils.so.1"
    "/lib/libproc.so.1"
    "/lib/libproc.so.2"
)

for file in "${suspicious_files[@]}"; do
    if [ -e "$file" ]; then
        echo "⚠️  SUSPICIOUS FILE FOUND: $file"
    fi
done

# Check for hidden processes
echo ""
echo "Hidden process detection:"
ps aux | awk '$3 > 50 {print "⚠️  HIGH CPU PROCESS: " $2 " (" $3 "%) - " $11}'

# Check for unusual network connections
echo ""
echo "Unusual network connections:"
netstat -ant | grep ESTABLISHED | awk '{print $5}' | cut -d: -f1 | sort | uniq -c | sort -nr | head -10

# 3. USER AND SUDO LOGS
echo ""
echo "3. USER AND SUDO ACTIVITY"
echo "-------------------------"

# Recent sudo commands
echo "Last 20 sudo commands:"
if [ -f /var/log/auth.log ]; then
    grep "sudo:" /var/log/auth.log | tail -20
elif [ -f /var/log/secure ]; then
    grep "sudo:" /var/log/secure | tail -20
else
    echo "No sudo log file found"
fi

# Recent user logins
echo ""
echo "Recent user logins:"
if [ -f /var/log/auth.log ]; then
    grep "session opened" /var/log/auth.log | tail -10
elif [ -f /var/log/secure ]; then
    grep "session opened" /var/log/secure | tail -10
else
    echo "No login log file found"
fi

# Check for new users
echo ""
echo "Users created in last 30 days:"
find /home -maxdepth 1 -type d -mtime -30 2>/dev/null | while read user; do
    username=$(basename "$user")
    echo "New user: $username (created: $(stat -c %y "$user" | cut -d' ' -f1))"
done

# 4. SYSTEM INTEGRITY CHECKS
echo ""
echo "4. SYSTEM INTEGRITY CHECKS"
echo "--------------------------"

# Check for modified system files
echo "Recently modified system files (last 7 days):"
find /etc -type f -mtime -7 2>/dev/null | head -10

# Check for unusual file permissions
echo ""
echo "Files with unusual permissions:"
find /etc -type f -perm /o+w 2>/dev/null | head -10

# Check for setuid binaries
echo ""
echo "Setuid binaries:"
find /usr/bin /usr/sbin /bin /sbin -type f -perm -4000 2>/dev/null | head -10

# 5. NETWORK SECURITY
echo ""
echo "5. NETWORK SECURITY"
echo "-------------------"

# Check firewall status
echo "Firewall status:"
if command -v ufw &> /dev/null; then
    ufw status
elif command -v iptables &> /dev/null; then
    iptables -L | head -20
else
    echo "No firewall detected"
fi

# Check for listening services
echo ""
echo "Services listening on all interfaces:"
ss -tlnp | grep "0.0.0.0"

echo ""
echo "=== SECURITY SCAN COMPLETED ==="
echo "Scan completed at: $(date)" 