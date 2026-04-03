#!/usr/bin/env bash
# collect-system-info.sh — Gathers system information for the engagement.
# This demonstrates a skill invoking a binary/script behind the scenes
# to collect data that feeds into the assessment.

set -euo pipefail

echo "============================================"
echo "  System Information Collection Report"
echo "  Generated: $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
echo "============================================"
echo ""

# --- OS & Kernel ---
echo "## Operating System"
echo "Hostname:       $(hostname)"
echo "OS:             $(uname -s)"
echo "Kernel:         $(uname -r)"
echo "Architecture:   $(uname -m)"
if [[ -f /etc/os-release ]]; then
    echo "Distribution:   $(grep PRETTY_NAME /etc/os-release 2>/dev/null | cut -d= -f2 | tr -d '"')"
elif command -v sw_vers &>/dev/null; then
    echo "Distribution:   macOS $(sw_vers -productVersion)"
fi
echo ""

# --- CPU ---
echo "## CPU"
if command -v lscpu &>/dev/null; then
    lscpu | grep -E "^(Model name|CPU\(s\)|Thread|Socket|Architecture)" 2>/dev/null || true
elif command -v sysctl &>/dev/null; then
    echo "CPU Brand:      $(sysctl -n machdep.cpu.brand_string 2>/dev/null || echo 'N/A')"
    echo "CPU Cores:      $(sysctl -n hw.ncpu 2>/dev/null || echo 'N/A')"
fi
echo ""

# --- Memory ---
echo "## Memory"
if command -v free &>/dev/null; then
    free -h 2>/dev/null | head -2
elif command -v sysctl &>/dev/null; then
    mem_bytes=$(sysctl -n hw.memsize 2>/dev/null || echo 0)
    echo "Total Memory:   $((mem_bytes / 1073741824)) GB"
fi
echo ""

# --- Disk ---
echo "## Disk Usage"
df -h / 2>/dev/null | head -5
echo ""

# --- Network Interfaces ---
echo "## Network Interfaces"
if command -v ip &>/dev/null; then
    ip -brief addr 2>/dev/null || true
elif command -v ifconfig &>/dev/null; then
    ifconfig 2>/dev/null | grep -E "^[a-z]|inet " | head -20
fi
echo ""

# --- Container Runtime ---
echo "## Container Runtime"
if command -v podman &>/dev/null; then
    echo "Podman:         $(podman --version 2>/dev/null)"
elif command -v docker &>/dev/null; then
    echo "Docker:         $(docker --version 2>/dev/null)"
else
    echo "No container runtime detected"
fi
echo ""

# --- Kubernetes / OpenShift ---
echo "## Kubernetes / OpenShift"
if command -v oc &>/dev/null; then
    echo "oc CLI:         $(oc version --client 2>/dev/null | head -1)"
elif command -v kubectl &>/dev/null; then
    echo "kubectl:        $(kubectl version --client --short 2>/dev/null || kubectl version --client 2>/dev/null | head -1)"
else
    echo "No Kubernetes CLI detected"
fi
echo ""

# --- Runtime Environments ---
echo "## Runtime Environments"
for cmd in java python3 python node go ruby dotnet; do
    if command -v "$cmd" &>/dev/null; then
        case "$cmd" in
            java)    echo "Java:           $($cmd -version 2>&1 | head -1)" ;;
            python3) echo "Python 3:       $($cmd --version 2>&1)" ;;
            python)  echo "Python:         $($cmd --version 2>&1)" ;;
            node)    echo "Node.js:        $($cmd --version 2>&1)" ;;
            go)      echo "Go:             $($cmd version 2>&1)" ;;
            ruby)    echo "Ruby:           $($cmd --version 2>&1)" ;;
            dotnet)  echo ".NET:           $($cmd --version 2>&1)" ;;
        esac
    fi
done
echo ""

echo "============================================"
echo "  Collection Complete"
echo "============================================"
