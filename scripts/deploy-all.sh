#!/bin/bash
# Master deployment script - runs all steps in sequence
# Usage: ./scripts/deploy-all.sh [PROJECT_ID] [ZONE] [REGION]

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

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🚀 Starting deployment for project: $PROJECT_ID"
echo "📍 Zone: $ZONE, Region: $REGION"
echo ""

# Step 1: Create VM
echo "📦 Step 1/6: Creating VM and static IP..."
"$SCRIPT_DIR/01-create-vm.sh" "$PROJECT_ID" "$ZONE" "$REGION"
echo ""

# Step 2: Configure firewall
echo "🔥 Step 2/6: Configuring firewall..."
"$SCRIPT_DIR/02-configure-firewall.sh"
echo ""

# Step 3: Assign static IP
echo "🌐 Step 3/6: Assigning static IP..."
"$SCRIPT_DIR/03-assign-static-ip.sh" "$ZONE" "$REGION"
echo ""

# Get static IP for display
STATIC_IP=$(gcloud compute addresses describe ssh-server-ip \
  --region="$REGION" \
  --format="get(address)" 2>/dev/null || echo "N/A")

echo "✅ Infrastructure setup complete!"
echo ""
echo "📝 Static IP: $STATIC_IP"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Next steps (run these ON THE VM):"
echo ""
echo "1. SSH into VM:"
echo "   gcloud compute ssh ssh-server --zone=$ZONE"
echo ""
echo "2. Clone scripts to VM (or copy them):"
echo "   git clone https://github.com/pc-style/pcstyledev-ssh.git"
echo "   cd pcstyledev-ssh"
echo ""
echo "3. Install dependencies:"
echo "   ./scripts/04-install-dependencies.sh"
echo ""
echo "4. Log out and back in (for docker group)"
echo ""
echo "5. Deploy application:"
echo "   ./scripts/05-deploy-app.sh"
echo ""
echo "6. Setup systemd service:"
echo "   ./scripts/08-setup-systemd.sh"
echo ""
echo "7. Configure DNS (manual step):"
echo "   Add A record: ssh.pcstyle.dev -> $STATIC_IP"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

