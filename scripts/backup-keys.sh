#!/bin/bash
# Backup SSH host keys
# Run this script ON THE VM
# Usage: ./scripts/backup-keys.sh

set -e

BACKUP_DIR="${HOME}/backups"
REPO_DIR="${HOME}/pcstyledev-ssh"
DATE=$(date +%Y%m%d)

mkdir -p "$BACKUP_DIR"

echo "Creating backup directory..."
mkdir -p "$BACKUP_DIR"

echo "Backing up keys from container..."
if docker ps | grep -q ssh-server; then
    docker cp ssh-server:/app/.ssh "$BACKUP_DIR/ssh-keys-$DATE" || {
        echo "Warning: Could not copy keys from container"
    }
else
    echo "Container not running, skipping container backup"
fi

echo "Backing up from host..."
if [ -d "$REPO_DIR/.ssh" ]; then
    tar -czf "$BACKUP_DIR/ssh-server-backup-$DATE.tar.gz" \
      "$REPO_DIR/.ssh" \
      "$REPO_DIR/go.mod" \
      "$REPO_DIR/Dockerfile" 2>/dev/null || {
        echo "Warning: Some files might not exist"
    }
else
    echo "Warning: .ssh directory not found at $REPO_DIR/.ssh"
fi

echo ""
echo "✅ Backup complete!"
echo "📦 Backup location: $BACKUP_DIR"
ls -lh "$BACKUP_DIR" | tail -5

