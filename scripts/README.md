# Deployment Scripts

Shell scripts to automate the deployment process for pcstyle.dev SSH server.

## 📋 Scripts Overview

### Local Machine Scripts (run from your computer)

1. **01-create-vm.sh** - Create VM instance and reserve static IP
2. **02-configure-firewall.sh** - Configure firewall rules
3. **03-assign-static-ip.sh** - Assign static IP to VM

### VM Scripts (run after SSH'ing into the VM)

4. **04-install-dependencies.sh** - Install Docker, Git, etc.
5. **05-deploy-app.sh** - Clone repo, build, and deploy container
6. **08-setup-systemd.sh** - Configure systemd service

### Maintenance Scripts (run on VM)

- **update-server.sh** - Update server code (quick or rolling)
- **backup-keys.sh** - Backup SSH host keys
- **setup-fail2ban.sh** - Configure Fail2Ban protection
- **setup-monitoring.sh** - Install Google Cloud monitoring agent

## 🚀 Quick Start

### Initial Deployment

```bash
# On your local machine
export PROJECT_ID="your-project-id"

# Step 1: Create VM
./scripts/01-create-vm.sh

# Step 2: Configure firewall
./scripts/02-configure-firewall.sh

# Step 3: Assign static IP
./scripts/03-assign-static-ip.sh

# Step 4: SSH into VM and install dependencies
gcloud compute ssh ssh-server --zone=us-central1-a
./scripts/04-install-dependencies.sh

# Log out and back in (for docker group)
exit
gcloud compute ssh ssh-server --zone=us-central1-a

# Step 5: Deploy application
./scripts/05-deploy-app.sh

# Step 6: Setup systemd service
./scripts/08-setup-systemd.sh
```

### Updating Server

```bash
# SSH into VM
gcloud compute ssh ssh-server --zone=us-central1-a

# Quick update (zero downtime)
./scripts/update-server.sh quick

# Or rolling update (simpler)
./scripts/update-server.sh rolling
```

### Backup

```bash
# On VM
./scripts/backup-keys.sh
```

## 📝 Script Details

### 01-create-vm.sh
Creates GCE instance and reserves static IP.

**Usage:**
```bash
./scripts/01-create-vm.sh [PROJECT_ID] [ZONE] [REGION]
# or
export PROJECT_ID="your-project-id"
./scripts/01-create-vm.sh
```

**Output:** Static IP address (save for DNS)

### 02-configure-firewall.sh
Creates firewall rule allowing SSH on port 22.

**Usage:**
```bash
./scripts/02-configure-firewall.sh
```

### 03-assign-static-ip.sh
Assigns reserved static IP to the VM.

**Usage:**
```bash
./scripts/03-assign-static-ip.sh [ZONE] [REGION]
```

### 04-install-dependencies.sh
Installs Docker, Git, and updates system packages.

**Usage:** Run on VM
```bash
./scripts/04-install-dependencies.sh
```

**Note:** Requires logout/login after running for docker group.

### 05-deploy-app.sh
Clones repository, builds Docker image, and starts container.

**Usage:** Run on VM
```bash
./scripts/05-deploy-app.sh [REPO_URL]
```

### 08-setup-systemd.sh
Configures systemd service for automatic container management.

**Usage:** Run on VM
```bash
./scripts/08-setup-systemd.sh
```

### update-server.sh
Updates server code with zero-downtime or rolling update.

**Usage:** Run on VM
```bash
# Quick update (zero downtime, requires manual testing)
./scripts/update-server.sh quick

# Rolling update (simpler, automatic)
./scripts/update-server.sh rolling
```

### backup-keys.sh
Backs up SSH host keys and important files.

**Usage:** Run on VM
```bash
./scripts/backup-keys.sh
```

**Output:** `~/backups/ssh-keys-YYYYMMDD` and `~/backups/ssh-server-backup-YYYYMMDD.tar.gz`

### setup-fail2ban.sh
Configures Fail2Ban for SSH brute force protection.

**Usage:** Run on VM
```bash
./scripts/setup-fail2ban.sh
```

### setup-monitoring.sh
Installs Google Cloud Ops Agent for monitoring.

**Usage:** Run on VM
```bash
./scripts/setup-monitoring.sh
```

## 🔧 Configuration

Most scripts use sensible defaults but can be customized:

- **ZONE**: Default `us-central1-a`
- **REGION**: Default `us-central1`
- **REPO_URL**: Default `https://github.com/pc-style/pcstyledev-ssh.git`
- **PROJECT_ID**: Must be set via environment variable or argument

## ⚠️ Notes

- Scripts use `set -e` to exit on errors
- Some scripts check for existing resources before creating
- VM scripts assume repository is cloned to `~/pcstyledev-ssh`
- Always test updates before deploying to production
- Backup keys regularly!

## 🐛 Troubleshooting

If a script fails:
1. Check error message for specific issue
2. Verify prerequisites are met
3. Check gcloud authentication: `gcloud auth list`
4. Verify project is set: `gcloud config get-value project`

For VM scripts:
- Ensure you're SSH'd into the VM
- Check Docker is running: `sudo systemctl status docker`
- Verify repository exists: `ls ~/pcstyledev-ssh`

