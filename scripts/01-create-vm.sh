#!/bin/bash
# Step 1: Create VM and reserve static IP
# Usage: ./scripts/01-create-vm.sh [PROJECT_ID] [ZONE] [REGION]

set -e

PROJECT_ID="${1:-${PROJECT_ID}}"
ZONE="${2:-us-central1-a}"
REGION="${3:-us-central1}"

if [ -z "$PROJECT_ID" ]; then
    echo "Error: PROJECT_ID is required"
    echo "Usage: $0 <PROJECT_ID> [ZONE] [REGION]"
    echo "   or: export PROJECT_ID=your-project-id && $0"
    exit 1
fi

echo "Setting project to $PROJECT_ID..."
gcloud config set project "$PROJECT_ID"

echo "Creating e2-micro instance..."
gcloud compute instances create ssh-server \
  --zone="$ZONE" \
  --machine-type=e2-micro \
  --image-family=ubuntu-2404-lts-amd64 \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=10GB \
  --boot-disk-type=pd-standard \
  --tags=ssh-server \
  --metadata=enable-oslogin=false

echo "Reserving static IP address..."
gcloud compute addresses create ssh-server-ip --region="$REGION" || {
    echo "Warning: Static IP might already exist, continuing..."
}

echo "Getting static IP address..."
STATIC_IP=$(gcloud compute addresses describe ssh-server-ip \
  --region="$REGION" \
  --format="get(address)")

echo ""
echo "✅ VM created successfully!"
echo "📝 Static IP: $STATIC_IP"
echo "💾 Save this IP for DNS configuration (Step 6)"
echo ""
echo "Next steps:"
echo "  1. Run: ./scripts/02-configure-firewall.sh"
echo "  2. Run: ./scripts/03-assign-static-ip.sh"

