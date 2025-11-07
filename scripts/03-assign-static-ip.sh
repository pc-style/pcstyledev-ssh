#!/bin/bash
# Step 3: Assign static IP to VM
# Usage: ./scripts/03-assign-static-ip.sh [ZONE] [REGION]

set -e

ZONE="${1:-us-central1-a}"
REGION="${2:-us-central1}"

echo "Getting static IP address..."
STATIC_IP=$(gcloud compute addresses describe ssh-server-ip \
  --region="$REGION" \
  --format="get(address)")

if [ -z "$STATIC_IP" ]; then
    echo "Error: Could not retrieve static IP. Make sure Step 1 completed successfully."
    exit 1
fi

echo "Static IP: $STATIC_IP"
echo "Removing existing access config..."
gcloud compute instances delete-access-config ssh-server \
  --zone="$ZONE" \
  --access-config-name="external-nat" || {
    echo "Warning: No existing access config found, continuing..."
}

echo "Assigning static IP to VM..."
gcloud compute instances add-access-config ssh-server \
  --zone="$ZONE" \
  --access-config-name="external-nat" \
  --address="$STATIC_IP"

echo ""
echo "✅ Static IP assigned successfully!"
echo "📝 IP: $STATIC_IP"
echo ""
echo "Next steps:"
echo "  1. SSH into VM: gcloud compute ssh ssh-server --zone=$ZONE"
echo "  2. Run: ./scripts/04-install-dependencies.sh (on the VM)"

