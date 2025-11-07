# SSH Server - Terminal Contact Interface

A beautiful terminal-based SSH interface for pcstyle.dev, allowing users to submit contact messages directly from their terminal.

Inspired by [Terminal Shop](https://github.com/charmbracelet/ssh-apps), this project demonstrates the power of SSH as an application delivery platform.

## Features

- **Universal Access**: Connect from any device with SSH (phones, tablets, desktops)
- **Beautiful TUI**: Built with Charm's Bubble Tea framework
- **Contact Form**: Submit messages with validation and error handling
- **Brand Styling**: Cyan and magenta color scheme matching pcstyle.dev
- **Navigation**: Intuitive keyboard-driven interface
- **API Integration**: Connects to existing `/api/contact` endpoint
- **No Browser Required**: Works on any terminal emulator

## Quick Start

### Local Testing

```bash
# Clone the repository
git clone https://github.com/pcstyle/ssh-server.git
cd ssh-server

# Build the server
go build -o bin/ssh-server ./cmd/server

# Start the server
./bin/ssh-server

# Connect from another terminal
ssh localhost -p 2222
```

### Using Docker

```bash
# Build the image
docker build -t ssh-server .

# Run the container
docker run -p 2222:2222 ssh-server

# Connect
ssh localhost -p 2222
```

## Architecture

```
ssh-server/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── server/
│   │   └── ssh.go            # Wish SSH server setup
│   ├── ui/
│   │   ├── app.go            # Main Bubble Tea app
│   │   ├── home.go           # Home page with navbar
│   │   ├── contact.go        # Contact form view
│   │   └── styles.go         # Lip Gloss styles
│   └── api/
│       └── client.go         # HTTP client for API
├── Dockerfile
├── DEPLOYMENT.md             # Google Cloud deployment guide
└── README.md
```

## Technology Stack

- **Language**: Go 1.21+
- **SSH Server**: [Wish](https://github.com/charmbracelet/wish)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles)

## Usage

### Navigation

- **Arrow Keys** or **j/k**: Navigate menu items
- **Enter**: Select menu item or submit form
- **Tab**: Move between form fields
- **Esc**: Go back to home
- **q** or **Ctrl+C**: Quit (from home screen)

### Contact Form

1. Connect to the server: `ssh ssh.pcstyle.dev`
2. Navigate to "Contact" and press Enter
3. Fill in the message (required) and optional fields:
   - Name
   - Email
   - Discord username
   - Phone number
4. Tab to "Submit" and press Enter
5. Wait for confirmation message

### Views

- **Home**: Welcome screen with navigation menu
- **Contact**: Contact form for sending messages
- **About**: Information about the project

## Configuration

The server accepts the following command-line flags:

```bash
./bin/ssh-server --help

Flags:
  -host string
        Host to bind to (default "0.0.0.0")
  -port int
        Port to listen on (default 2222)
  -api string
        API base URL (default "https://pcstyle.dev")
```

### Example

```bash
# Run on custom port with different API
./bin/ssh-server -port 3000 -api https://staging.pcstyle.dev
```

## Development

### Prerequisites

- Go 1.21 or higher
- SSH client for testing

### Building from Source

```bash
# Install dependencies
go mod download

# Build
go build -o bin/ssh-server ./cmd/server

# Generate SSH host keys (if not exists)
mkdir -p .ssh
ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N "" -C "ssh-server-host-key"

# Run
./bin/ssh-server
```

### Testing Locally

```bash
# Terminal 1: Start server
./bin/ssh-server

# Terminal 2: Connect
ssh localhost -p 2222

# Or with explicit options
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost -p 2222
```

## Deployment

See [DEPLOYMENT.md](./DEPLOYMENT.md) for comprehensive deployment instructions to:

- Google Compute Engine (Recommended)
- Google Cloud Run (Limited support)
- Google Kubernetes Engine (For scale)

### Production Deployment Quick Start

```bash
# 1. Create GCE instance
gcloud compute instances create ssh-server \
  --zone=us-central1-a \
  --machine-type=e2-micro \
  --image-family=ubuntu-2204-lts \
  --image-project=ubuntu-os-cloud

# 2. Configure firewall
gcloud compute firewall-rules create allow-ssh-server \
  --allow=tcp:22 \
  --target-tags=ssh-server

# 3. SSH and deploy
gcloud compute ssh ssh-server --zone=us-central1-a
# ... follow deployment steps in DEPLOYMENT.md

# 4. Configure DNS
# Add A record: ssh.pcstyle.dev → [STATIC_IP]

# 5. Connect
ssh ssh.pcstyle.dev
```

## API Integration

The contact form integrates with the existing `/api/contact` endpoint documented in [API_CONTACT_ENDPOINT.md](./API_CONTACT_ENDPOINT.md).

### Request Format

```json
{
  "message": "Hello from SSH!",
  "name": "John Doe",
  "email": "john@example.com",
  "discord": "@johndoe",
  "phone": "+1234567890",
  "source": "ssh"
}
```

### Response Format

```json
{
  "success": true,
  "message": "Message sent successfully! Thanks for reaching out."
}
```

## Security

- **Anonymous Access**: The server allows all connections (public access)
- **Rate Limiting**: Basic rate limiting included (can be enhanced)
- **Input Validation**: All form inputs are validated before submission
- **HTTPS API**: Uses HTTPS for API communication
- **SSH Encryption**: All traffic encrypted via SSH protocol

### Recommendations for Production

1. Implement fail2ban for brute force protection
2. Use firewall rules to limit connection rates
3. Monitor logs for suspicious activity
4. Keep dependencies updated
5. Use static IPs with DNS
6. Set up Cloud Monitoring alerts

## Troubleshooting

### Connection Issues

```bash
# Check if server is running
ps aux | grep ssh-server

# Check if port is listening
lsof -i :2222

# Test with verbose SSH
ssh -vvv localhost -p 2222
```

### Build Issues

```bash
# Clean and rebuild
go clean
go mod tidy
go build -o bin/ssh-server ./cmd/server
```

### Docker Issues

```bash
# View logs
docker logs ssh-server

# Restart container
docker restart ssh-server

# Rebuild image
docker build --no-cache -t ssh-server .
```

## Future Enhancements

Potential features to add:

- **Portfolio Viewer**: Browse projects and work samples
- **Interactive Resume**: Navigate through experience and skills
- **Visitor Guestbook**: Leave messages for others to see
- **Easter Eggs**: Hidden games (Snake, Tetris)
- **Real-time Chat**: IRC-style chat with other visitors
- **System Stats**: Live server statistics and visitor counter
- **Code Snippets**: Browse code examples and tutorials
- **ASCII Art Gallery**: Showcase beautiful terminal art

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

## License

[Your License Here]

## Acknowledgments

- [Charm](https://charm.sh) for the amazing TUI libraries
- [Terminal Shop](https://github.com/charmbracelet/ssh-apps) for inspiration
- The Go SSH community for excellent documentation

## Contact

- **Website**: https://pcstyle.dev
- **SSH**: `ssh ssh.pcstyle.dev` (when deployed)
- **GitHub**: https://github.com/pcstyle/ssh-server

---

**Made with ❤️ using Charm's Bubble Tea**
