# Quick Start Guide

Get your SSH server running in 2 minutes!

## Test It Now (Local)

Open two terminal windows:

### Terminal 1: Start the Server

```bash
cd /Users/pcstyle/projects/pcstyledev-ssh
./bin/ssh-server
```

You should see:
```
╔════════════════════════════════════════════════════════════╗
║  SSH Server starting on 0.0.0.0:2222
║                                                            ║
║  Connect with: ssh localhost -p 2222
║                                                            ║
║  Press Ctrl+C to stop                                      ║
╚════════════════════════════════════════════════════════════╝
```

### Terminal 2: Connect to the Server

```bash
ssh localhost -p 2222
```

**Note**: You might see a host key warning. Type `yes` to continue.

## What You'll See

1. **Welcome Screen**:
   - ASCII art banner
   - Navigation menu with three options:
     - Contact
     - About
     - Exit

2. **Navigation**:
   - Use ↑/↓ or j/k to move
   - Press Enter to select

3. **Contact Form**:
   - Message field (required)
   - Optional fields: Name, Email, Discord, Phone
   - Tab to move between fields
   - Enter on "Submit" to send

4. **About Page**:
   - Information about the project
   - Press Enter or Esc to go back

## Try These Commands

```bash
# Connect and skip host key verification (for testing)
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost -p 2222

# Connect with verbose output (debugging)
ssh -vvv localhost -p 2222

# Test from Docker
docker run -p 2222:2222 ssh-server
ssh localhost -p 2222
```

## Testing the Contact Form

1. Connect to the server
2. Navigate to "Contact" and press Enter
3. Type a test message: `Testing SSH contact form!`
4. (Optional) Add your name and email
5. Tab to "Submit" and press Enter
6. You should see: `✓ Message sent successfully! Thanks for reaching out.`
7. Check your Discord webhook to see the message!

## Keyboard Shortcuts Reference

| Key | Action |
|-----|--------|
| ↑/↓ or j/k | Navigate menu |
| Enter | Select/Submit |
| Tab | Next field |
| Shift+Tab | Previous field |
| Esc | Go back |
| q | Quit (from home) |
| Ctrl+C | Force quit |

## Troubleshooting

### "Connection refused"
- Make sure the server is running in Terminal 1
- Check port 2222 isn't in use: `lsof -i :2222`

### "Host key verification failed"
- Remove old host key: `ssh-keygen -R "[localhost]:2222"`
- Or use `-o StrictHostKeyChecking=no` flag

### "Permission denied"
- The server allows all connections
- Try: `ssh user@localhost -p 2222` (any username works)

### Form submission fails
- Check server logs in Terminal 1
- Verify API endpoint is accessible: `curl https://pcstyle.dev/api/contact`
- Check DISCORD_WEBHOOK_URL is set in production

## Next Steps

Once local testing works:

1. **Deploy to Google Cloud**: See [DEPLOYMENT.md](./DEPLOYMENT.md)
2. **Configure DNS**: Point `ssh.pcstyle.dev` to your server IP
3. **Test from anywhere**: `ssh ssh.pcstyle.dev`
4. **Share with friends**: They can connect from their terminal!

## Build from Source

If you make changes:

```bash
# Rebuild
go build -o bin/ssh-server ./cmd/server

# Or with Docker
docker build -t ssh-server .
docker run -p 2222:2222 ssh-server
```

## Mobile Testing

Yes, it works on phones!

### iOS (using Termius or Prompt)
1. Install Termius or Prompt from App Store
2. Add new host: `ssh.pcstyle.dev` (or `localhost:2222` for testing)
3. Connect
4. Enjoy the full TUI experience on your phone!

### Android (using Termux or JuiceSSH)
1. Install Termux or JuiceSSH
2. Connect: `ssh ssh.pcstyle.dev`
3. Works perfectly!

## Production Checklist

Before deploying to production:

- [ ] Test all form fields
- [ ] Verify Discord webhook receives messages
- [ ] Test on different terminal emulators (macOS, Linux, Windows)
- [ ] Test on mobile devices
- [ ] Check error handling (try submitting empty form)
- [ ] Monitor server logs for issues
- [ ] Set up DNS for `ssh.pcstyle.dev`
- [ ] Configure firewall rules
- [ ] Set up monitoring and alerts
- [ ] Document access for your team

---

**Enjoy your beautiful SSH interface!** 🎨✨
