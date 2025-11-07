#!/bin/bash
# Setup Google Cloud monitoring agent
# Run this script ON THE VM
# Usage: ./scripts/setup-monitoring.sh

set -e

echo "Installing Google Cloud Ops Agent..."
curl -sSO https://dl.google.com/cloudagents/add-google-cloud-ops-agent-repo.sh
sudo bash add-google-cloud-ops-agent-repo.sh --also-install

echo "Cleaning up installer..."
rm -f add-google-cloud-ops-agent-repo.sh

echo ""
echo "✅ Monitoring agent installed!"
echo ""
echo "📊 Next steps:"
echo "  1. Go to Cloud Console: Monitoring > Uptime checks"
echo "  2. Create uptime check:"
echo "     - Protocol: TCP"
echo "     - Port: 22"
echo "     - Host: ssh.pcstyle.dev"

