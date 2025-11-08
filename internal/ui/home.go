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
	menuItems     []MenuItem
	cursor        int
	width         int
	height        int
	secretBuffer  string // buffer for secret commands
	konamiBuffer  string // buffer for konami code
	hiddenVisible bool   // show hidden menu items
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
		cursor:        0,
		secretBuffer:  "",
		konamiBuffer:  "",
		hiddenVisible: false,
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
		key := msg.String()
		
		// handle konami code detection
		konamiCode := "↑↑↓↓←→←→ba"
		m.konamiBuffer += key
		if len(m.konamiBuffer) > len(konamiCode) {
			m.konamiBuffer = m.konamiBuffer[len(m.konamiBuffer)-len(konamiCode):]
		}
		if strings.Contains(m.konamiBuffer, konamiCode) {
			m.hiddenVisible = !m.hiddenVisible
			m.konamiBuffer = ""
			return m, nil
		}
		
		// handle secret command buffer (typing commands)
		if len(key) == 1 && ((key[0] >= 'a' && key[0] <= 'z') || (key[0] >= 'A' && key[0] <= 'Z')) {
			m.secretBuffer += strings.ToLower(key)
			if len(m.secretBuffer) > 20 {
				m.secretBuffer = m.secretBuffer[len(m.secretBuffer)-20:]
			}
			
			// check for secret commands
			if strings.Contains(m.secretBuffer, "snake") {
				m.secretBuffer = ""
				return m, func() tea.Msg {
					return NavigateToGameMsg{GameType: "snake"}
				}
			}
			if strings.Contains(m.secretBuffer, "matrix") {
				m.secretBuffer = ""
				return m, func() tea.Msg {
					return NavigateToGameMsg{GameType: "matrix"}
				}
			}
			if strings.Contains(m.secretBuffer, "konami") {
				m.secretBuffer = ""
				m.hiddenVisible = !m.hiddenVisible
				return m, nil
			}
			if strings.Contains(m.secretBuffer, "help") {
				m.secretBuffer = ""
				return m, func() tea.Msg {
					return NavigateToGameMsg{GameType: "help"}
				}
			}
		} else {
			// reset buffer on non-letter keys
			m.secretBuffer = ""
		}
		
		switch key {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			maxItems := len(m.menuItems)
			if m.hiddenVisible {
				maxItems += 2 // snake + matrix
			}
			if m.cursor < maxItems-1 {
				m.cursor++
			}
		case "enter", " ":
			// check if selecting hidden item
			if m.hiddenVisible && m.cursor >= len(m.menuItems) {
				hiddenIndex := m.cursor - len(m.menuItems)
				if hiddenIndex == 0 {
					return m, func() tea.Msg {
						return NavigateToGameMsg{GameType: "snake"}
					}
				}
				if hiddenIndex == 1 {
					return m, func() tea.Msg {
						return NavigateToGameMsg{GameType: "matrix"}
					}
				}
			}
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
	
	// hidden menu items (easter eggs)
	if m.hiddenVisible {
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("🎮 HIDDEN GAMES 🎮"))
		b.WriteString("\n")
		
		hiddenItems := []struct {
			title string
			desc  string
		}{
			{"🐍 Snake", "Classic snake game"},
			{"💊 Matrix", "Matrix rain effect"},
		}
		
		for i, item := range hiddenItems {
			idx := len(m.menuItems) + i
			cursor := "  "
			if idx == m.cursor {
				cursor = NavArrowStyle.Render("→ ")
			}
			
			itemStyle := NavItemStyle
			if idx == m.cursor {
				itemStyle = NavItemSelectedStyle
			}
			
			title := itemStyle.Render(item.title)
			desc := HelpStyle.Render(item.desc)
			b.WriteString(fmt.Sprintf("%s%s - %s\n", cursor, title, desc))
		}
	}

	// Help text
	b.WriteString("\n")
	helpText := "Use ↑/↓ or j/k to navigate • Press Enter to select • Press q to quit"
	b.WriteString(HelpStyle.Render(helpText))
	
	// easter egg hint (very subtle)
	if !m.hiddenVisible {
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("💡 Tip: Try typing 'snake' or 'matrix'..."))
	}

	return BaseStyle.Render(b.String())
}
