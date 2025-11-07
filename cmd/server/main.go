package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/pcstyle/ssh-server/internal/server"
)

func main() {
	// Parse command-line flags
	host := flag.String("host", "0.0.0.0", "Host to bind to")
	port := flag.Int("port", 2222, "Port to listen on")
	apiURL := flag.String("api", "https://pcstyle.dev", "API base URL")
	flag.Parse()

	// Configure logger
	log.SetLevel(log.InfoLevel)
	log.SetReportTimestamp(true)

	// Create server configuration
	config := server.Config{
		Host:       *host,
		Port:       *port,
		APIBaseURL: *apiURL,
	}

	// Create and start the server
	srv, err := server.NewServer(config)
	if err != nil {
		log.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	// Print connection info
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Printf("║  SSH Server starting on %s:%d                    \n", *host, *port)
	fmt.Println("║                                                            ║")
	fmt.Printf("║  Connect with: ssh localhost -p %d                      \n", *port)
	fmt.Println("║                                                            ║")
	fmt.Println("║  Press Ctrl+C to stop                                      ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Start the server (blocks until interrupted)
	if err := srv.Start(); err != nil {
		log.Error("Server error", "error", err)
		os.Exit(1)
	}
}
