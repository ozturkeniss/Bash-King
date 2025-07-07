#!/bin/bash

echo "=== FILE BACKUP SCRIPT ==="
echo "Date: $(date)"
echo ""

# Create backup directory
BACKUP_DIR="/tmp/backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p $BACKUP_DIR

echo "Creating backup directory: $BACKUP_DIR"

# Backup important files
echo "Backing up important files..."

# Backup system info
echo "=== System Info Backup ===" > $BACKUP_DIR/system_info.txt
uname -a >> $BACKUP_DIR/system_info.txt
echo "" >> $BACKUP_DIR/system_info.txt
echo "=== Memory Info ===" >> $BACKUP_DIR/system_info.txt
free -h >> $BACKUP_DIR/system_info.txt
echo "" >> $BACKUP_DIR/system_info.txt
echo "=== Disk Usage ===" >> $BACKUP_DIR/system_info.txt
df -h >> $BACKUP_DIR/system_info.txt

# Backup process list
echo "=== Process List ===" > $BACKUP_DIR/processes.txt
ps aux >> $BACKUP_DIR/processes.txt

# Backup network info
echo "=== Network Info ===" > $BACKUP_DIR/network.txt
ip addr show >> $BACKUP_DIR/network.txt 2>/dev/null || echo "ip command not available" >> $BACKUP_DIR/network.txt

echo "Backup completed!"
echo "Files created:"
ls -la $BACKUP_DIR/
echo ""
echo "Backup directory size:"
du -sh $BACKUP_DIR/ 