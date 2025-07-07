#!/bin/bash

echo "=== HOST PACKAGE MANAGER ANALYSIS ==="
echo "Timestamp: $(date)"
echo "Hostname: $(hostname)"
echo ""

echo "=== SYSTEM PACKAGE MANAGER ==="
echo "Detecting Package Manager:"
if command -v apt >/dev/null 2>&1; then
    echo "Package Manager: APT (Debian/Ubuntu)"
    PKG_MGR="apt"
elif command -v yum >/dev/null 2>&1; then
    echo "Package Manager: YUM (RHEL/CentOS)"
    PKG_MGR="yum"
elif command -v dnf >/dev/null 2>&1; then
    echo "Package Manager: DNF (Fedora/RHEL)"
    PKG_MGR="dnf"
elif command -v pacman >/dev/null 2>&1; then
    echo "Package Manager: PACMAN (Arch Linux)"
    PKG_MGR="pacman"
elif command -v zypper >/dev/null 2>&1; then
    echo "Package Manager: ZYPPER (openSUSE)"
    PKG_MGR="zypper"
else
    echo "Package Manager: UNKNOWN"
    PKG_MGR="unknown"
fi
echo ""

echo "=== PACKAGE REPOSITORIES ==="
case $PKG_MGR in
    "apt")
        echo "APT Sources:"
        cat /etc/apt/sources.list 2>/dev/null | head -10
        echo ""
        echo "Additional Sources:"
        ls /etc/apt/sources.list.d/ 2>/dev/null | head -5
        echo ""
        ;;
    "yum"|"dnf")
        echo "YUM/DNF Repositories:"
        yum repolist 2>/dev/null || dnf repolist 2>/dev/null
        echo ""
        ;;
    "pacman")
        echo "Pacman Repositories:"
        cat /etc/pacman.conf 2>/dev/null | grep -E "^\[.*\]"
        echo ""
        ;;
esac

echo "=== PACKAGE UPDATE STATUS ==="
case $PKG_MGR in
    "apt")
        echo "Checking for APT updates:"
        apt update 2>/dev/null
        echo "Available updates:"
        apt list --upgradable 2>/dev/null | head -10
        echo ""
        ;;
    "yum")
        echo "Checking for YUM updates:"
        yum check-update 2>/dev/null | head -10
        echo ""
        ;;
    "dnf")
        echo "Checking for DNF updates:"
        dnf check-update 2>/dev/null | head -10
        echo ""
        ;;
    "pacman")
        echo "Checking for Pacman updates:"
        pacman -Qu 2>/dev/null | head -10
        echo ""
        ;;
esac

echo "=== INSTALLED PACKAGES ==="
echo "Total installed packages:"
case $PKG_MGR in
    "apt")
        dpkg -l | wc -l
        echo ""
        echo "Recently installed packages:"
        dpkg -l | tail -10
        echo ""
        ;;
    "yum"|"dnf")
        rpm -qa | wc -l
        echo ""
        echo "Recently installed packages:"
        rpm -qa --last | head -10
        echo ""
        ;;
    "pacman")
        pacman -Q | wc -l
        echo ""
        echo "Recently installed packages:"
        pacman -Q | tail -10
        echo ""
        ;;
esac

echo "=== PACKAGE DEPENDENCIES ==="
echo "Checking for broken dependencies:"
case $PKG_MGR in
    "apt")
        apt-get check 2>/dev/null || echo "No broken dependencies found"
        echo ""
        ;;
    "yum")
        package-cleanup --problems 2>/dev/null || echo "No broken dependencies found"
        echo ""
        ;;
    "dnf")
        dnf check 2>/dev/null || echo "No broken dependencies found"
        echo ""
        ;;
    "pacman")
        pacman -Dk 2>/dev/null || echo "No broken dependencies found"
        echo ""
        ;;
esac

echo "=== PACKAGE CACHE STATUS ==="
case $PKG_MGR in
    "apt")
        echo "APT cache size:"
        du -sh /var/cache/apt/archives/ 2>/dev/null || echo "Cache size not available"
        echo ""
        echo "APT cache statistics:"
        apt-cache stats 2>/dev/null
        echo ""
        ;;
    "yum"|"dnf")
        echo "YUM/DNF cache size:"
        du -sh /var/cache/yum/ 2>/dev/null || echo "Cache size not available"
        echo ""
        ;;
    "pacman")
        echo "Pacman cache size:"
        du -sh /var/cache/pacman/pkg/ 2>/dev/null || echo "Cache size not available"
        echo ""
        ;;
esac

echo "=== SECURITY UPDATES ==="
echo "Checking for security updates:"
case $PKG_MGR in
    "apt")
        apt list --upgradable 2>/dev/null | grep -i security | head -5
        echo ""
        ;;
    "yum")
        yum updateinfo list security 2>/dev/null | head -5
        echo ""
        ;;
    "dnf")
        dnf updateinfo list security 2>/dev/null | head -5
        echo ""
        ;;
esac

echo "=== PACKAGE MANAGER SUMMARY ==="
echo "Package Manager: $PKG_MGR"
case $PKG_MGR in
    "apt")
        echo "Total Packages: $(dpkg -l | wc -l)"
        echo "Available Updates: $(apt list --upgradable 2>/dev/null | wc -l)"
        ;;
    "yum"|"dnf")
        echo "Total Packages: $(rpm -qa | wc -l)"
        echo "Available Updates: $(yum check-update 2>/dev/null | wc -l)"
        ;;
    "pacman")
        echo "Total Packages: $(pacman -Q | wc -l)"
        echo "Available Updates: $(pacman -Qu 2>/dev/null | wc -l)"
        ;;
esac
echo ""

echo "=== END OF PACKAGE MANAGER ANALYSIS ===" 