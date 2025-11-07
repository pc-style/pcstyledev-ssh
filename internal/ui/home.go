package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// MenuItem represents a navigation item
type MenuItem struct {
	Title       string
	Description string
}

// NavigateMsg is sent when user selects a menu item
type NavigateMsg int

// HomeModel represents the home page with navigation
type HomeModel struct {
	menuItems []MenuItem
	cursor    int
	width     int
	height    int
}

// NewHomeModel creates a new home page model
func NewHomeModel() HomeModel {
	return HomeModel{
		menuItems: []MenuItem{
			{
				Title:       "Contact",
				Description: "Send me a message",
			},
			{
				Title:       "About",
				Description: "Learn more about this project",
			},
			{
				Title:       "Exit",
				Description: "Disconnect from SSH",
			},
		},
		cursor: 0,
	}
}

// Init initializes the home model
func (m HomeModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the home model
func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.menuItems)-1 {
				m.cursor++
			}
		case "enter", " ":
			// Send navigation message
			return m, func() tea.Msg {
				return NavigateMsg(m.cursor)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the home page
func (m HomeModel) View() string {
	var b strings.Builder

	// ASCII Banner
	banner := `
╔═══════════════════════════════════════╗
║                                       ║
║         P C S T Y L E . D E V         ║
║                                       ║
║         SSH Terminal Interface        ║
║                                       ║
╚═══════════════════════════════════════╝
`

	b.WriteString(TitleStyle.Render(banner))
	b.WriteString("\n")

	// Welcome message
	welcome := "Welcome to pcstyle.dev SSH interface"
	b.WriteString(TitleStyle.Width(m.width).Render(welcome))
	b.WriteString("\n\n")

	// Navigation menu
	for i, item := range m.menuItems {
		cursor := "  "
		if i == m.cursor {
			cursor = NavArrowStyle.Render("→ ")
		}

		// Menu item style
		itemStyle := NavItemStyle
		if i == m.cursor {
			itemStyle = NavItemSelectedStyle
		}

		title := itemStyle.Render(item.Title)
		desc := HelpStyle.Render(item.Description)

		b.WriteString(fmt.Sprintf("%s%s - %s\n", cursor, title, desc))
	}

	// Help text
	b.WriteString("\n")
	helpText := "Use ↑/↓ or j/k to navigate • Press Enter to select • Press q to quit"
	b.WriteString(HelpStyle.Render(helpText))

	return BaseStyle.Render(b.String())
}
