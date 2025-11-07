package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/pcstyle/ssh-server/internal/ui"
)

// Config holds the server configuration
type Config struct {
	Host       string
	Port       int
	APIBaseURL string
}

// Server represents the SSH server
type Server struct {
	config Config
	ssh    *ssh.Server
}

// NewServer creates a new SSH server
func NewServer(config Config) (*Server, error) {
	s := &Server{
		config: config,
	}

	// Create the SSH server with Wish middleware
	sshServer, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", config.Host, config.Port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			// Allow all connections (public access)
			return true
		}),
		wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
			// Allow all connections (public access)
			return true
		}),
		wish.WithMiddleware(
			bubbletea.Middleware(s.teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH server: %w", err)
	}

	s.ssh = sshServer
	return s, nil
}

// teaHandler creates a Bubble Tea program for each SSH session
func (s *Server) teaHandler(sshSession ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Get terminal info
	pty, _, _ := sshSession.Pty()
	log.Info("Terminal", "type", pty.Term, "width", pty.Window.Width, "height", pty.Window.Height)

	// Set up Lip Gloss renderer with color support for this output
	renderer := lipgloss.NewRenderer(sshSession)
	renderer.SetHasDarkBackground(true)

	// Create a new app model for this session with the renderer
	model := ui.NewModel(s.config.APIBaseURL, renderer)

	// Configure the Bubble Tea program with proper I/O
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithInput(sshSession),
		tea.WithOutput(sshSession),
	}

	return model, opts
}

// Start starts the SSH server
func (s *Server) Start() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Info("Starting SSH server", "host", s.config.Host, "port", s.config.Port)
		if err := s.ssh.ListenAndServe(); err != nil {
			log.Error("SSH server error", "error", err)
		}
	}()

	// Wait for interrupt signal
	<-done

	log.Info("Shutting down SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.ssh.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	log.Info("Server stopped")
	return nil
}
