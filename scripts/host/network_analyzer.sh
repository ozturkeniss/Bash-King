#!/bin/bash

echo "=== HOST NETWORK ANALYSIS ==="
echo "Timestamp: $(date)"
echo "Hostname: $(hostname)"
echo ""

echo "=== NETWORK INTERFACES ==="
echo "Active Network Interfaces:"
ip addr show | grep -E "inet|UP|DOWN"
echo ""

echo "Interface Details:"
for interface in $(ip link show | grep -E "^[0-9]+:" | awk -F: '{print $2}' | tr -d ' '); do
    echo "Interface: $interface"
    ip addr show $interface 2>/dev/null | grep -E "inet|link"
    echo ""
done

echo "=== ROUTING TABLE ==="
echo "Current Routing Table:"
ip route show
echo ""

echo "=== DNS CONFIGURATION ==="
echo "DNS Servers:"
cat /etc/resolv.conf
echo ""

echo "=== NETWORK CONNECTIVITY ==="
echo "Testing Internet Connectivity:"
ping -c 3 8.8.8.8 2>/dev/null && echo "Internet: CONNECTED" || echo "Internet: DISCONNECTED"
echo ""

echo "Testing DNS Resolution:"
nslookup google.com 2>/dev/null && echo "DNS: WORKING" || echo "DNS: FAILED"
echo ""

echo "=== NETWORK STATISTICS ==="
echo "Network Interface Statistics:"
cat /proc/net/dev
echo ""

echo "=== ACTIVE CONNECTIONS ==="
echo "TCP Connections:"
ss -tuln | head -15
echo ""

echo "Listening Ports:"
ss -tuln | grep LISTEN | head -10
echo ""

echo "=== NETWORK PERFORMANCE ==="
echo "Bandwidth Usage (if available):"
if command -v iftop >/dev/null 2>&1; then
    echo "iftop available - run manually for real-time monitoring"
else
    echo "iftop not installed"
fi
echo ""

echo "Network Latency Test:"
ping -c 3 google.com 2>/dev/null | tail -1
echo ""

echo "=== FIREWALL STATUS ==="
echo "UFW Status:"
ufw status 2>/dev/null || echo "UFW not available"
echo ""

echo "iptables Rules:"
iptables -L 2>/dev/null | head -10 || echo "iptables not available"
echo ""

echo "=== NETWORK SERVICES ==="
echo "Running Network Services:"
netstat -tuln 2>/dev/null | head -10 || ss -tuln | head -10
echo ""

echo "=== WIRELESS NETWORK (if applicable) ==="
if command -v iwconfig >/dev/null 2>&1; then
    echo "Wireless Interfaces:"
    iwconfig 2>/dev/null | grep -E "ESSID|Frequency|Quality" || echo "No wireless interfaces found"
else
    echo "iwconfig not available"
fi
echo ""

echo "=== NETWORK SUMMARY ==="
echo "Total Network Interfaces: $(ip link show | grep -c '^[0-9]')"
echo "Active Connections: $(ss -tuln | wc -l)"
echo "Default Gateway: $(ip route | grep default | awk '{print $3}' | head -1)"
echo "Primary DNS: $(grep nameserver /etc/resolv.conf | head -1 | awk '{print $2}')"
echo ""

echo "=== END OF NETWORK ANALYSIS ===" 