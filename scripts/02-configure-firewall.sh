#!/bin/bash
# Step 2: Configure firewall rules
# Usage: ./scripts/02-configure-firewall.sh

set -e

echo "Creating firewall rule for SSH on port 22..."
gcloud compute firewall-rules create allow-ssh-terminal \
  --allow=tcp:22 \
  --target-tags=ssh-server \
  --source-ranges=0.0.0.0/0 \
  --description="Allow public SSH access to terminal server" || {
    echo "Warning: Firewall rule might already exist, continuing..."
}

echo "Verifying firewall rule..."
gcloud compute firewall-rules list --filter="name=allow-ssh-terminal"

echo ""
echo "✅ Firewall configured successfully!"
echo ""
echo "Next step: Run: ./scripts/03-assign-static-ip.sh"

