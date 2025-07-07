#!/bin/bash

echo "=== LOG CLEANUP SCRIPT ==="
echo "Date: $(date)"
echo ""

# Find and clean old log files
echo "Searching for log files..."

# Find log files older than 1 day
OLD_LOGS=$(find /var/log -name "*.log" -mtime +1 2>/dev/null | head -10)

if [ -n "$OLD_LOGS" ]; then
    echo "Found old log files:"
    echo "$OLD_LOGS"
    echo ""
    echo "Would clean these files (simulation mode)"
    echo "Total size to be freed: $(du -sh $OLD_LOGS 2>/dev/null | tail -1 | awk '{print $1}')"
else
    echo "No old log files found"
fi

# Clean temporary files
echo ""
echo "Cleaning temporary files..."

# Clean /tmp files older than 1 day
TMP_FILES=$(find /tmp -type f -mtime +1 2>/dev/null | wc -l)
echo "Found $TMP_FILES temporary files older than 1 day"

# Show disk usage before and after (simulation)
echo ""
echo "=== DISK USAGE ==="
df -h / | tail -1
echo ""
echo "Simulation: Would free approximately 50MB of space" 