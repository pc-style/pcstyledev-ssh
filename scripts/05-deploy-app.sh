#!/bin/bash
# Step 5: Deploy application
# Run this script ON THE VM after installing dependencies
# Usage: ./scripts/05-deploy-app.sh [REPO_URL]

set -e

REPO_URL="${1:-https://github.com/pc-style/pcstyledev-ssh.git}"
REPO_DIR="$HOME/pcstyledev-ssh"

echo "Cloning repository..."
if [ -d "$REPO_DIR" ]; then
    echo "Repository already exists, pulling latest changes..."
    cd "$REPO_DIR"
    git pull
else
    cd ~
    git clone "$REPO_URL"
    cd "$REPO_DIR"
fi

echo "Building Docker image..."
docker build -t ssh-server .

echo "Stopping existing container if running..."
docker stop ssh-server 2>/dev/null || true
docker rm ssh-server 2>/dev/null || true

echo "Starting container on port 22..."
docker run -d \
  --name ssh-server \
  --restart unless-stopped \
  -p 22:2222 \
  ssh-server

echo "Waiting for container to start..."
sleep 2

echo "Checking container status..."
docker ps | grep ssh-server

echo ""
echo "✅ Application deployed!"
echo ""
echo "View logs: docker logs -f ssh-server"
echo "Test connection: ssh localhost -p 22"
echo ""
echo "Next step: Run: ./scripts/08-setup-systemd.sh"

