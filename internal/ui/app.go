package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pcstyle/ssh-server/internal/api"
)

// View represents different screens in the app
type View int

const (
	ViewHome View = iota
	ViewContact
	ViewAbout
	ViewExit
)

// Model is the main application model
type Model struct {
	currentView  View
	homeModel    HomeModel
	contactModel ContactModel
	width        int
	height       int
	quitting     bool
	renderer     *lipgloss.Renderer
}

// NewModel creates a new application model
func NewModel(apiBaseURL string, renderer *lipgloss.Renderer) Model {
	apiClient := api.NewClient(apiBaseURL)

	m := Model{
		currentView:  ViewHome,
		homeModel:    NewHomeModel(),
		contactModel: NewContactModel(apiClient),
		renderer:     renderer,
	}

	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.currentView == ViewHome {
				m.quitting = true
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case NavigateMsg:
		// Handle navigation from home screen
		switch int(msg) {
		case 0: // Contact
			m.currentView = ViewContact
		case 1: // About
			m.currentView = ViewAbout
		case 2: // Exit
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil

	case BackMsg:
		// Handle back navigation
		m.currentView = ViewHome
		return m, nil
	}

	// Route updates to the appropriate view
	var cmd tea.Cmd
	switch m.currentView {
	case ViewHome:
		m.homeModel, cmd = m.homeModel.Update(msg)
	case ViewContact:
		m.contactModel, cmd = m.contactModel.Update(msg)
	case ViewAbout:
		// Handle about view
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "esc" || key.String() == "enter" {
				m.currentView = ViewHome
			}
		}
	case ViewExit:
		m.quitting = true
		return m, tea.Quit
	}

	if m.quitting {
		return m, tea.Quit
	}

	return m, cmd
}

// View renders the current view
func (m Model) View() string {
	if m.quitting {
		return GoodbyeView()
	}

	switch m.currentView {
	case ViewHome:
		return m.homeModel.View()
	case ViewContact:
		return m.contactModel.View()
	case ViewAbout:
		return AboutView()
	default:
		return ""
	}
}

// AboutView renders the about page
func AboutView() string {
	var b strings.Builder

	// Title
	b.WriteString(TitleStyle.Render("About pcstyle.dev"))
	b.WriteString("\n\n")

	// Name section
	b.WriteString(NavItemSelectedStyle.Render("Adam Krupa (@pcstyle)"))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	b.WriteString("\n\n")

	// WHO section
	b.WriteString(LabelStyle.Render("WHO"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("18 years old • Częstochowa, Poland"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("AI Student @ Politechnika Częstochowska"))
	b.WriteString("\n\n")

	// WHAT section
	b.WriteString(LabelStyle.Render("WHAT"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("Blending AI, design, and creative coding. Focused on neo-brutalist"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("design aesthetics combined with interactive and generative technologies."))
	b.WriteString("\n\n")

	// SKILLS section
	b.WriteString(LabelStyle.Render("SKILLS"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Frontend: Next.js 16, React 19, TypeScript, Tailwind v4, Framer Motion"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Graphics: WebGL shaders, generative art"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• AI/Backend: Python, custom generative pipelines"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Areas: AI, ML, Creative Coding, Interactive Design"))
	b.WriteString("\n\n")

	// PROJECTS section
	b.WriteString(LabelStyle.Render("PROJECTS"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Clock Gallery - Interactive animated art (clock.pcstyle.dev)"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• AimDrift - Precision aim trainer (driftfield.pcstyle.dev)"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• PoliCalc - Grade calculator (kalkulator.pcstyle.dev)"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• PixelForge - AI-powered image editor (pixlab.pcstyle.dev)"))
	b.WriteString("\n\n")

	// EXPLORING section
	b.WriteString(LabelStyle.Render("EXPLORING"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Realtime AI workflow agents for animations"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Neo-brutalist design system tokenization"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Interactive SSH contact UX with WebRTC fallback"))
	b.WriteString("\n\n")

	// CONNECT section
	b.WriteString(LabelStyle.Render("CONNECT"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("GitHub: github.com/pcstyle"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("Twitter: @pcstyle"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("Email: adamkrupa@tuta.io"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("Calendar: cal.com/pcstyle"))
	b.WriteString("\n\n")

	// Footer
	b.WriteString(HelpStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("Built with Go + Charm (Wish, Bubble Tea, Lip Gloss)"))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("Source: github.com/pc-style/pcstyledev-ssh"))
	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("Press Enter or Esc to go back"))

	return BoxStyle.Render(b.String())
}

// GoodbyeView renders the goodbye message
func GoodbyeView() string {
	goodbye := `
  _____ _                 _                       _
 |_   _| |__   __ _ _ __ | | __  _   _  ___  _  _| |
   | | | '_ \ / _' | '_ \| |/ / | | | |/ _ \| || | |
   | | | | | | (_| | | | |   <  | |_| | (_) | || |_|
   |_| |_| |_|\__,_|_| |_|_|\_\  \__, |\___/ \_,_(_)
                                 |___/

Thanks for visiting pcstyle.dev via SSH!
`
	return TitleStyle.Render(goodbye) + "\n"
}
