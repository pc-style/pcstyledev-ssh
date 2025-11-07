#!/bin/bash
# Step 4: Install dependencies on VM
# Run this script ON THE VM after SSH'ing in
# Usage: ./scripts/04-install-dependencies.sh

set -e

echo "Updating system packages..."
sudo apt-get update && sudo apt-get upgrade -y

echo "Installing Docker..."
sudo apt-get install -y docker.io docker-compose-v2

echo "Starting Docker service..."
sudo systemctl enable docker
sudo systemctl start docker

echo "Adding user to docker group..."
sudo usermod -aG docker "$USER"

echo "Installing Git..."
sudo apt-get install -y git

echo ""
echo "✅ Dependencies installed!"
echo ""
echo "⚠️  IMPORTANT: Log out and back in for docker group to take effect"
echo "   Then run: ./scripts/05-deploy-app.sh"

