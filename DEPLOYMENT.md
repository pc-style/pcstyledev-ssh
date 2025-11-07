# SSH Server Deployment Guide

Complete guide for deploying the pcstyle.dev SSH server to production.

## 🎯 Quick Start (Recommended: Google Compute Engine)

### Prerequisites
- Google Cloud account with billing enabled
- `gcloud` CLI installed and configured (`gcloud auth login`)
- Git repository created and pushed
- Domain access for DNS configuration

---

## 📦 Option 1: Google Compute Engine (Best for SSH)

### Step 1: Create and Configure VM

```bash
# Set your project ID
export PROJECT_ID="your-project-id"
gcloud config set project $PROJECT_ID

# Create e2-micro instance (free tier eligible: $5-6/month credit)
gcloud compute instances create ssh-server \
  --zone=us-central1-a \
  --machine-type=e2-micro \
  --image-family=ubuntu-2404-lts-amd64 \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=10GB \
  --boot-disk-type=pd-standard \
  --tags=ssh-server \
  --metadata=enable-oslogin=false

# Reserve a static IP address
gcloud compute addresses create ssh-server-ip --region=us-central1

# Get the IP address (save this for DNS)
gcloud compute addresses describe ssh-server-ip \
  --region=us-central1 \
  --format="get(address)"
```

### Step 2: Configure Firewall

```bash
# Allow SSH traffic on port 22
gcloud compute firewall-rules create allow-ssh-terminal \
  --allow=tcp:22 \
  --target-tags=ssh-server \
  --source-ranges=0.0.0.0/0 \
  --description="Allow public SSH access to terminal server"

# Verify firewall rule
gcloud compute firewall-rules list --filter="name=allow-ssh-terminal"
```

### Step 3: Assign Static IP to VM

```bash
# Get the static IP
STATIC_IP=$(gcloud compute addresses describe ssh-server-ip \
  --region=us-central1 \
  --format="get(address)")

# Assign it to the VM
gcloud compute instances delete-access-config ssh-server \
  --zone=us-central1-a \
  --access-config-name="external-nat"

gcloud compute instances add-access-config ssh-server \
  --zone=us-central1-a \
  --access-config-name="external-nat" \
  --address=$STATIC_IP
```

### Step 4: Install Dependencies on VM

```bash
# SSH into the instance
gcloud compute ssh ssh-server --zone=us-central1-a

# Update system
sudo apt-get update && sudo apt-get upgrade -y

# Install Docker
sudo apt-get install -y docker.io docker-compose-v2
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker $USER

# Install Git
sudo apt-get install -y git

# Log out and back in for docker group
exit
gcloud compute ssh ssh-server --zone=us-central1-a
```

### Step 5: Deploy Application

```bash
# Clone your repository
cd ~
git clone https://github.com/pc-style/pcstyledev-ssh.git
cd pcstyledev-ssh

# Build Docker image
docker build -t ssh-server .

# Run container on port 22 (maps internal 2222 to external 22)
docker run -d \
  --name ssh-server \
  --restart unless-stopped \
  -p 22:2222 \
  ssh-server

# Verify it's running
docker ps
docker logs -f ssh-server
```

### Step 6: Configure DNS

1. **Go to your DNS provider** (Cloudflare, Google Domains, etc.)
2. **Add an A record**:
   - **Name**: `ssh` (creates `ssh.pcstyle.dev`)
   - **Type**: `A`
   - **Value**: `[YOUR_STATIC_IP]` (from Step 1)
   - **TTL**: `Auto` or `300` seconds
   - **Proxy**: Disable (must be DNS only for SSH)

3. **Wait for DNS propagation** (usually 1-5 minutes)

### Step 7: Test Connection

```bash
# Test from your local machine
ssh ssh.pcstyle.dev

# Or with explicit port (if not using 22)
ssh ssh.pcstyle.dev -p 22

# Verify DNS resolution
dig ssh.pcstyle.dev +short
# Should show your static IP

# Test with verbose output for debugging
ssh -vvv ssh.pcstyle.dev
```

### Step 8: Set Up Systemd Service (Recommended)

```bash
# SSH into the VM
gcloud compute ssh ssh-server --zone=us-central1-a

# Create systemd service file
sudo tee /etc/systemd/system/ssh-server.service > /dev/null <<'EOF'
[Unit]
Description=SSH Terminal Server
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/YOUR_USERNAME/pcstyledev-ssh
ExecStart=/usr/bin/docker start ssh-server || /usr/bin/docker run -d --name ssh-server --restart unless-stopped -p 22:2222 ssh-server
ExecStop=/usr/bin/docker stop ssh-server

[Install]
WantedBy=multi-user.target
EOF

# Replace YOUR_USERNAME with your actual username
sudo sed -i "s/YOUR_USERNAME/$USER/g" /etc/systemd/system/ssh-server.service

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable ssh-server
sudo systemctl start ssh-server

# Check status
sudo systemctl status ssh-server
```

---

## 🔄 Updating the Server

### Quick Update (Zero Downtime)

```bash
# SSH into VM
gcloud compute ssh ssh-server --zone=us-central1-a

# Pull latest changes
cd ~/pcstyledev-ssh
git pull

# Rebuild image
docker build -t ssh-server-new .

# Start new container on different port temporarily
docker run -d --name ssh-server-new -p 2222:2222 ssh-server-new

# Test new container
ssh localhost -p 2222

# If working, swap containers
docker stop ssh-server
docker rm ssh-server
docker rename ssh-server-new ssh-server
docker stop ssh-server
docker run -d --name ssh-server --restart unless-stopped -p 22:2222 ssh-server-new

# Cleanup
docker rmi ssh-server-new
docker image prune -f
```

### Rolling Update with Systemd

```bash
# Update code
cd ~/pcstyledev-ssh && git pull && docker build -t ssh-server .

# Restart service (handles container restart automatically)
sudo systemctl restart ssh-server

# Check logs
docker logs -f ssh-server
```

---

## 📊 Monitoring & Maintenance

### View Logs

```bash
# Docker logs
docker logs -f ssh-server

# Last 100 lines
docker logs --tail=100 ssh-server

# Systemd logs
sudo journalctl -u ssh-server -f

# Cloud Logging (from local machine)
gcloud logging read "resource.type=gce_instance" --limit=50
```

### Monitor Resources

```bash
# Real-time container stats
docker stats ssh-server

# VM resource usage
htop

# Disk usage
df -h
docker system df
```

### Set Up Monitoring Alerts

```bash
# Install monitoring agent (optional but recommended)
curl -sSO https://dl.google.com/cloudagents/add-google-cloud-ops-agent-repo.sh
sudo bash add-google-cloud-ops-agent-repo.sh --also-install

# Create uptime check in Cloud Console
# Navigation: Monitoring > Uptime checks > Create uptime check
# Protocol: TCP
# Port: 22
# Host: ssh.pcstyle.dev
```

---

## 🔒 Security Best Practices

### 1. Enable Fail2Ban (Recommended)

```bash
sudo apt-get install -y fail2ban

# Configure for SSH
sudo tee /etc/fail2ban/jail.local > /dev/null <<EOF
[sshd]
enabled = true
port = 22
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 3600
findtime = 600
EOF

sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# Check status
sudo fail2ban-client status sshd
```

### 2. Regular Updates

```bash
# Create update script
cat > ~/update-server.sh <<'EOF'
#!/bin/bash
set -e

echo "Updating system packages..."
sudo apt-get update && sudo apt-get upgrade -y

echo "Updating SSH server..."
cd ~/pcstyledev-ssh
git pull
docker build -t ssh-server .
sudo systemctl restart ssh-server

echo "Cleaning up..."
docker image prune -f

echo "Update complete!"
docker logs --tail=20 ssh-server
EOF

chmod +x ~/update-server.sh

# Run monthly via cron
(crontab -l 2>/dev/null; echo "0 2 1 * * /home/$USER/update-server.sh >> /home/$USER/update.log 2>&1") | crontab -
```

### 3. Backup SSH Host Keys

```bash
# Backup keys from container
mkdir -p ~/backups
docker cp ssh-server:/app/.ssh ~/backups/ssh-keys-$(date +%Y%m%d)

# Or backup from host
tar -czf ~/backups/ssh-server-backup-$(date +%Y%m%d).tar.gz \
  ~/pcstyledev-ssh/.ssh \
  ~/pcstyledev-ssh/go.mod \
  ~/pcstyledev-ssh/Dockerfile
```

### 4. Firewall Configuration

```bash
# Install UFW (alternative to GCP firewall)
sudo apt-get install -y ufw

# Allow SSH on port 22
sudo ufw allow 22/tcp

# Enable firewall
sudo ufw --force enable

# Check status
sudo ufw status
```

---

## 💰 Cost Estimates (2025)

### Google Compute Engine e2-micro
| Resource | Cost | Notes |
|----------|------|-------|
| VM Instance | $0-7/month | Free tier: ~$5-6 credit/month |
| Static IP (in use) | $0 | Free when attached |
| Storage (10GB) | $0.40/month | Standard persistent disk |
| Network egress | First 1GB free | Then $0.12/GB (NA) |
| **Total** | **$0-8/month** | Usually free tier covers it |

### Cost Optimization Tips
- Use e2-micro (free tier eligible)
- Keep disk at 10GB minimum
- Use standard persistent disk (not SSD)
- Monitor bandwidth usage
- Delete unused snapshots

---

## 🐛 Troubleshooting

### Connection Refused

```bash
# Check if container is running
docker ps

# Check container logs
docker logs ssh-server

# Check if port 22 is listening
sudo netstat -tlnp | grep :22

# Check firewall rules
gcloud compute firewall-rules list --filter="name=allow-ssh-terminal"

# Test locally on VM
ssh localhost -p 22
```

### DNS Not Resolving

```bash
# Check DNS propagation
dig ssh.pcstyle.dev +short
nslookup ssh.pcstyle.dev

# Use alternative DNS servers
dig @8.8.8.8 ssh.pcstyle.dev +short
dig @1.1.1.1 ssh.pcstyle.dev +short

# Clear local DNS cache (Mac)
sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder

# Test direct IP connection
ssh [STATIC_IP]
```

### Container Won't Start

```bash
# Check Docker logs
docker logs ssh-server

# Check Docker daemon
sudo systemctl status docker

# Rebuild image
cd ~/pcstyledev-ssh
docker build --no-cache -t ssh-server .

# Remove old container and start fresh
docker rm -f ssh-server
docker run -d --name ssh-server --restart unless-stopped -p 22:2222 ssh-server
```

### Port Already in Use

```bash
# Find what's using port 22
sudo lsof -i :22

# If it's OpenSSH server, disable it (we're replacing it)
sudo systemctl stop ssh
sudo systemctl disable ssh

# Or change OpenSSH to different port
sudo sed -i 's/#Port 22/Port 2200/' /etc/ssh/sshd_config
sudo systemctl restart ssh
```

### High Memory Usage

```bash
# Check container resources
docker stats ssh-server

# Restart container
docker restart ssh-server

# If needed, upgrade to e2-small
gcloud compute instances stop ssh-server --zone=us-central1-a
gcloud compute instances set-machine-type ssh-server \
  --machine-type=e2-small \
  --zone=us-central1-a
gcloud compute instances start ssh-server --zone=us-central1-a
```

---

## 🚀 Alternative Deployment Options

### Option 2: Railway (Easiest, but limited SSH support)

**Note**: Railway is HTTP-focused. SSH may require workarounds.

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login
railway login

# Initialize project
railway init

# Deploy
railway up
```

### Option 3: DigitalOcean Droplet

```bash
# Create droplet (via dashboard or doctl CLI)
doctl compute droplet create ssh-server \
  --size s-1vcpu-1gb \
  --image ubuntu-24-04-x64 \
  --region nyc1

# SSH in and follow GCE deployment steps above
```

### Option 4: Fly.io (Good for global distribution)

```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Login
flyctl auth login

# Initialize and deploy
flyctl launch
flyctl deploy
```

---

## 📝 Post-Deployment Checklist

- [ ] SSH server accessible at `ssh.pcstyle.dev`
- [ ] Colors working properly
- [ ] Navigation working (↑/↓, Enter)
- [ ] Contact form submits successfully
- [ ] Discord webhook configured in Vercel
- [ ] DNS A record configured
- [ ] Firewall rules in place
- [ ] Systemd service enabled
- [ ] Monitoring alerts configured
- [ ] Backup script created
- [ ] Update script scheduled
- [ ] Fail2ban installed and configured
- [ ] SSH host keys backed up
- [ ] Documentation updated

---

## 📞 Support

If you encounter issues:

1. **Check logs**: `docker logs ssh-server`
2. **Test locally**: `ssh localhost -p 22` (on VM)
3. **Verify DNS**: `dig ssh.pcstyle.dev +short`
4. **Check firewall**: `gcloud compute firewall-rules list`
5. **Review**: GitHub Issues or contact form

---

**Made with ❤️ using Go + Charm's Bubble Tea**
