package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// NavigateToGameMsg is sent when navigating to a game
type NavigateToGameMsg struct {
	GameType string
}

// HelpEasterEggModel shows help for easter eggs
type HelpEasterEggModel struct {
	startTime time.Time
}

// NewHelpEasterEggModel creates a new help easter egg
func NewHelpEasterEggModel() HelpEasterEggModel {
	return HelpEasterEggModel{
		startTime: time.Now(),
	}
}

// Init initializes the help model
func (m HelpEasterEggModel) Init() tea.Cmd {
	return nil
}

// Update handles help updates
func (m HelpEasterEggModel) Update(msg tea.Msg) (HelpEasterEggModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter":
			return m, func() tea.Msg {
				return BackMsg{}
			}
		}
	}
	return m, nil
}

// View renders the help easter egg
func (m HelpEasterEggModel) View() string {
	var b strings.Builder
	
	b.WriteString(TitleStyle.Render("EASTER EGGS & SECRETS"))
	b.WriteString("\n\n")
	
	b.WriteString(LabelStyle.Render("HIDDEN COMMANDS:"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Type 'snake' - Play Snake game"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Type 'matrix' - Matrix rain effect"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Type 'konami' - Toggle hidden menu"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Type 'help' - Show this help"))
	b.WriteString("\n\n")
	
	b.WriteString(LabelStyle.Render("KONAMI CODE:"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("↑ ↑ ↓ ↓ ← → ← → B A"))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("(or type 'konami' as shortcut)"))
	b.WriteString("\n\n")
	
	b.WriteString(LabelStyle.Render("GAMES:"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("Snake - Classic arcade game"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("  Controls: WASD/Arrow keys/HJKL"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("  Space: Pause, R: Restart"))
	b.WriteString("\n\n")
	b.WriteString(NavItemStyle.Render("Matrix - Digital rain effect"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("  Just watch and enjoy!"))
	b.WriteString("\n\n")
	
	b.WriteString(LabelStyle.Render("FUN FACTS:"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• This SSH interface is built with Go + Bubble Tea"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• The matrix effect uses katakana characters"))
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("• Snake game speed increases as you score"))
	b.WriteString("\n\n")
	
	b.WriteString(HelpStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("Press Esc/Q/Enter to go back"))
	
	return BoxStyle.Render(b.String())
}

// KonamiSuccessModel shows konami code success message
type KonamiSuccessModel struct {
	startTime time.Time
}

// NewKonamiSuccessModel creates konami success message
func NewKonamiSuccessModel() KonamiSuccessModel {
	return KonamiSuccessModel{
		startTime: time.Now(),
	}
}

// Init initializes konami model
func (m KonamiSuccessModel) Init() tea.Cmd {
	return nil
}

// Update handles konami updates
func (m KonamiSuccessModel) Update(msg tea.Msg) (KonamiSuccessModel, tea.Cmd) {
	// auto-close after 3 seconds
	if time.Since(m.startTime) > 3*time.Second {
		return m, func() tea.Msg {
			return BackMsg{}
		}
	}
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter":
			return m, func() tea.Msg {
				return BackMsg{}
			}
		}
	}
	return m, nil
}

// View renders konami success
func (m KonamiSuccessModel) View() string {
	var b strings.Builder
	
	art := `
╔═══════════════════════════════════════╗
║                                       ║
║     KONAMI CODE ACTIVATED!             ║
║                                       ║
║     Hidden features unlocked!         ║
║                                       ║
╚═══════════════════════════════════════╝
`
	
	b.WriteString(TitleStyle.Render(art))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("Hidden menu items are now visible!"))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("(Auto-closing in a moment...)"))
	
	return BoxStyle.Render(b.String())
}
