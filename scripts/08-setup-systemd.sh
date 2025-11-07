#!/bin/bash
# Step 8: Set up systemd service
# Run this script ON THE VM
# Usage: ./scripts/08-setup-systemd.sh

set -e

REPO_DIR="${HOME}/pcstyledev-ssh"
USERNAME=$(whoami)

if [ ! -d "$REPO_DIR" ]; then
    echo "Error: Repository not found at $REPO_DIR"
    echo "Make sure Step 5 completed successfully."
    exit 1
fi

echo "Creating systemd service file..."
sudo tee /etc/systemd/system/ssh-server.service > /dev/null <<EOF
[Unit]
Description=SSH Terminal Server
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=${REPO_DIR}
ExecStart=/usr/bin/docker start ssh-server || /usr/bin/docker run -d --name ssh-server --restart unless-stopped -p 22:2222 ssh-server
ExecStop=/usr/bin/docker stop ssh-server

[Install]
WantedBy=multi-user.target
EOF

echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

echo "Enabling service..."
sudo systemctl enable ssh-server

echo "Starting service..."
sudo systemctl start ssh-server

echo "Checking service status..."
sudo systemctl status ssh-server --no-pager || true

echo ""
echo "✅ Systemd service configured!"
echo ""
echo "Useful commands:"
echo "  Check status: sudo systemctl status ssh-server"
echo "  View logs: sudo journalctl -u ssh-server -f"
echo "  Restart: sudo systemctl restart ssh-server"

