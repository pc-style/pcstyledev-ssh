#!/bin/bash
# Setup Fail2Ban for SSH protection
# Run this script ON THE VM
# Usage: ./scripts/setup-fail2ban.sh

set -e

echo "Installing Fail2Ban..."
sudo apt-get update
sudo apt-get install -y fail2ban

echo "Configuring Fail2Ban for SSH..."
sudo tee /etc/fail2ban/jail.local > /dev/null <<EOF
[sshd]
enabled = true
port = 22
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 3600
findtime = 600
EOF

echo "Enabling Fail2Ban service..."
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

echo "Checking Fail2Ban status..."
sudo fail2ban-client status sshd

echo ""
echo "✅ Fail2Ban configured!"
echo ""
echo "Useful commands:"
echo "  Check status: sudo fail2ban-client status sshd"
echo "  Unban IP: sudo fail2ban-client set sshd unbanip <IP>"

