# SSH Server Deployment Guide

This guide covers deploying the SSH server to Google Cloud Platform.

## Prerequisites

- Google Cloud account
- `gcloud` CLI installed and configured
- Docker installed locally
- Domain `ssh.pcstyle.dev` configured

## Option 1: Google Compute Engine (Recommended)

### Step 1: Create a VM Instance

```bash
# Create an e2-micro instance (free tier eligible)
gcloud compute instances create ssh-server \
  --zone=us-central1-a \
  --machine-type=e2-micro \
  --image-family=ubuntu-2204-lts \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=10GB \
  --tags=ssh-server

# Reserve a static IP
gcloud compute addresses create ssh-server-ip --region=us-central1

# Get the static IP
gcloud compute addresses describe ssh-server-ip --region=us-central1 --format="get(address)"
```

### Step 2: Configure Firewall

```bash
# Allow SSH traffic on port 22
gcloud compute firewall-rules create allow-ssh-server \
  --allow=tcp:22 \
  --target-tags=ssh-server \
  --description="Allow SSH traffic on port 22"
```

### Step 3: Deploy the Application

```bash
# SSH into the instance
gcloud compute ssh ssh-server --zone=us-central1-a

# Install Docker
sudo apt-get update
sudo apt-get install -y docker.io docker-compose
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker $USER

# Log out and back in for docker group to take effect
exit
gcloud compute ssh ssh-server --zone=us-central1-a

# Clone your repository or copy files
git clone https://github.com/pcstyle/ssh-server.git
cd ssh-server

# Build the Docker image
docker build -t ssh-server .

# Run the container
docker run -d \
  --name ssh-server \
  --restart unless-stopped \
  -p 22:2222 \
  ssh-server

# Check logs
docker logs -f ssh-server
```

### Step 4: Configure DNS

1. Go to your DNS provider (e.g., Cloudflare, Google Domains)
2. Add an A record:
   - Name: `ssh` (or `ssh.pcstyle.dev`)
   - Type: `A`
   - Value: `[YOUR_STATIC_IP]`
   - TTL: `Auto` or `300`

### Step 5: Test Connection

```bash
# From your local machine
ssh ssh.pcstyle.dev

# Or if using a custom user
ssh user@ssh.pcstyle.dev
```

### Step 6: Set Up Automatic Updates (Optional)

Create a systemd service for automatic Docker container restarts:

```bash
# Create systemd service
sudo tee /etc/systemd/system/ssh-server.service > /dev/null <<EOF
[Unit]
Description=SSH Server Docker Container
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/$USER/ssh-server
ExecStart=/usr/bin/docker run -d --name ssh-server --restart unless-stopped -p 22:2222 ssh-server
ExecStop=/usr/bin/docker stop ssh-server
ExecStopPost=/usr/bin/docker rm ssh-server

[Install]
WantedBy=multi-user.target
EOF

# Enable and start the service
sudo systemctl enable ssh-server
sudo systemctl start ssh-server

# Check status
sudo systemctl status ssh-server
```

## Option 2: Google Cloud Run (Limited SSH Support)

**Note:** Cloud Run primarily supports HTTP/HTTPS traffic. For SSH, Compute Engine is recommended.

If you want to try Cloud Run with websocket SSH:

```bash
# Build and push to Container Registry
gcloud builds submit --tag gcr.io/[PROJECT_ID]/ssh-server

# Deploy to Cloud Run (won't work for direct SSH)
gcloud run deploy ssh-server \
  --image gcr.io/[PROJECT_ID]/ssh-server \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

## Option 3: Google Kubernetes Engine (For Scale)

If you need high availability and scaling:

```bash
# Create a GKE cluster
gcloud container clusters create ssh-cluster \
  --zone=us-central1-a \
  --num-nodes=1 \
  --machine-type=e2-small

# Get credentials
gcloud container clusters get-credentials ssh-cluster --zone=us-central1-a

# Build and push image
gcloud builds submit --tag gcr.io/[PROJECT_ID]/ssh-server

# Create deployment
kubectl create deployment ssh-server --image=gcr.io/[PROJECT_ID]/ssh-server

# Expose as LoadBalancer
kubectl expose deployment ssh-server --type=LoadBalancer --port=22 --target-port=2222

# Get external IP
kubectl get service ssh-server
```

## Monitoring & Logs

### Compute Engine

```bash
# View logs
docker logs -f ssh-server

# Monitor resource usage
docker stats ssh-server
```

### Cloud Logging

```bash
# View VM logs
gcloud logging read "resource.type=gce_instance AND resource.labels.instance_id=[INSTANCE_ID]"
```

## Security Best Practices

### 1. Rate Limiting

The server includes basic rate limiting. For production, consider:

```bash
# Use fail2ban
sudo apt-get install fail2ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

### 2. Update Regularly

```bash
# Update system packages
sudo apt-get update && sudo apt-get upgrade -y

# Update Docker image
cd ~/ssh-server
git pull
docker build -t ssh-server .
docker stop ssh-server
docker rm ssh-server
docker run -d --name ssh-server --restart unless-stopped -p 22:2222 ssh-server
```

### 3. Backup SSH Keys

```bash
# Backup host keys
docker cp ssh-server:/app/.ssh ./backup-ssh-keys
```

## Cost Estimates (2025)

### Compute Engine e2-micro
- **Instance**: $0 (free tier) or ~$7-10/month
- **Storage**: $0.40/month for 10GB
- **Network**: First 1GB egress free, then ~$0.12/GB
- **Total**: ~$0-10/month depending on usage

### Static IP
- **In use**: Free
- **Reserved but unused**: $0.01/hour (~$7/month)

## Troubleshooting

### Connection Refused

```bash
# Check if container is running
docker ps

# Check firewall rules
gcloud compute firewall-rules list

# Check if port is listening
sudo netstat -tlnp | grep 22
```

### DNS Not Resolving

```bash
# Test DNS resolution
dig ssh.pcstyle.dev

# Check A record
nslookup ssh.pcstyle.dev
```

### Container Crashes

```bash
# View logs
docker logs ssh-server

# Restart container
docker restart ssh-server

# Check resource usage
docker stats ssh-server
```

## Updating the Server

### Rolling Update

```bash
# SSH into the instance
gcloud compute ssh ssh-server --zone=us-central1-a

# Pull latest changes
cd ~/ssh-server
git pull

# Rebuild image
docker build -t ssh-server .

# Stop old container
docker stop ssh-server
docker rm ssh-server

# Start new container
docker run -d \
  --name ssh-server \
  --restart unless-stopped \
  -p 22:2222 \
  ssh-server

# Verify
docker logs -f ssh-server
```

## Scaling & High Availability

For high traffic, consider:

1. **Load Balancing**: Use Google Cloud Load Balancer
2. **Multiple Regions**: Deploy to multiple zones
3. **Auto-scaling**: Use managed instance groups
4. **Health Checks**: Implement health check endpoints

## Support & Maintenance

- Monitor logs regularly
- Set up Cloud Monitoring alerts
- Keep Docker images updated
- Review security patches monthly
- Backup SSH host keys

## Quick Reference

```bash
# Start server
./bin/ssh-server

# Build Docker image
docker build -t ssh-server .

# Run Docker container
docker run -p 2222:2222 ssh-server

# Deploy to GCE
gcloud compute ssh ssh-server --zone=us-central1-a

# View logs
docker logs -f ssh-server

# Test connection
ssh localhost -p 2222  # Local
ssh ssh.pcstyle.dev     # Production
```
