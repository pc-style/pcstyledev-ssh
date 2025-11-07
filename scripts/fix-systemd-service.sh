#!/bin/bash
# Quick fix for broken systemd service
# Run this script ON THE VM to fix the existing service
# Usage: ./scripts/fix-systemd-service.sh

set -e

REPO_DIR="${HOME}/pcstyledev-ssh"

if [ ! -d "$REPO_DIR" ]; then
    echo "Error: Repository not found at $REPO_DIR"
    exit 1
fi

echo "Creating wrapper script..."
sudo tee /usr/local/bin/ssh-server-start.sh > /dev/null <<'SCRIPT'
#!/bin/bash
if docker ps -a --format '{{.Names}}' | grep -q '^ssh-server$'; then
    docker start ssh-server
else
    docker run -d --name ssh-server --restart unless-stopped -p 22:2222 ssh-server
fi
SCRIPT

sudo chmod +x /usr/local/bin/ssh-server-start.sh

echo "Updating systemd service file..."
sudo tee /etc/systemd/system/ssh-server.service > /dev/null <<EOF
[Unit]
Description=SSH Terminal Server
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=${REPO_DIR}
ExecStart=/usr/local/bin/ssh-server-start.sh
ExecStop=/usr/bin/docker stop ssh-server

[Install]
WantedBy=multi-user.target
EOF

echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

echo "Restarting service..."
sudo systemctl restart ssh-server

echo "Checking service status..."
sleep 2
sudo systemctl status ssh-server --no-pager || true

echo ""
echo "✅ Service fixed!"

