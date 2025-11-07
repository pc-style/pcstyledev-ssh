#!/bin/bash
# Update server script (runs on VM)
# Usage: ./scripts/update-server.sh [update-type]
# update-type: quick (zero downtime) or rolling (default)

set -e

UPDATE_TYPE="${1:-rolling}"
REPO_DIR="${HOME}/pcstyledev-ssh"

if [ ! -d "$REPO_DIR" ]; then
    echo "Error: Repository not found at $REPO_DIR"
    exit 1
fi

cd "$REPO_DIR"

if [ "$UPDATE_TYPE" = "quick" ]; then
    echo "🔄 Quick update (zero downtime)..."
    
    echo "Pulling latest changes..."
    git pull
    
    echo "Building new image..."
    docker build -t ssh-server-new .
    
    echo "Starting new container on port 2222..."
    docker run -d --name ssh-server-new -p 2222:2222 ssh-server-new
    
    echo "Testing new container..."
    echo "⚠️  Test with: ssh localhost -p 2222"
    echo "Press Enter to continue after testing, or Ctrl+C to cancel..."
    read -r
    
    echo "Swapping containers..."
    docker stop ssh-server || true
    docker rm ssh-server || true
    docker rename ssh-server-new ssh-server
    docker stop ssh-server
    docker run -d --name ssh-server --restart unless-stopped -p 22:2222 ssh-server-new
    
    echo "Cleaning up..."
    docker rmi ssh-server-new 2>/dev/null || true
    docker image prune -f
    
    echo "✅ Quick update complete!"
else
    echo "🔄 Rolling update..."
    
    echo "Updating code..."
    git pull
    
    echo "Rebuilding image..."
    docker build -t ssh-server .
    
    echo "Restarting service..."
    sudo systemctl restart ssh-server
    
    echo "✅ Rolling update complete!"
fi

echo ""
echo "View logs: docker logs -f ssh-server"

