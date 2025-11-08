package ui

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
)

// MenuItem opisuje menu entry, niby obvious ale trzeba
type MenuItem struct {
	Title       string
	Description string
	Target      View
	isSecret    bool
}

// NavigateMsg leci gdy wybierzesz coś z listy, simple
type NavigateMsg struct {
	Target View
}

// HomeModel pilnuje strony startowej, zero magii
type HomeModel struct {
	menuItems      []MenuItem
	secretItems    []MenuItem
	cursor         int
	width          int
	height         int
	secretUnlocked bool
	secretBuffer   string
	secretMessage  string
	lastUnlockPing time.Time
}

// NewHomeModel składa menu bazowe, plus secret stash
func NewHomeModel() HomeModel {
	base := []MenuItem{
		{
			Title:       "Contact",
			Description: "Send me a message",
			Target:      ViewContact,
		},
		{
			Title:       "About",
			Description: "Learn more about this project",
			Target:      ViewAbout,
		},
		{
			Title:       "Exit",
			Description: "Disconnect from SSH",
			Target:      ViewExit,
		},
	}

	secrets := []MenuItem{
		{
			Title:       "Arcade",
			Description: "play snake + dziwne rzeczy",
			Target:      ViewArcade,
			isSecret:    true,
		},
		{
			Title:       "???",
			Description: "weird logbook, nie oceniaj",
			Target:      ViewSecrets,
			isSecret:    true,
		},
	}

	return HomeModel{
		menuItems:   base,
		secretItems: secrets,
		cursor:      0,
	}
}

// Init niby nic nie robi, ale bubbletea chce
func (m HomeModel) Init() tea.Cmd {
	return nil
}

// Update ogarnia klawisze, trochę też odpala sekrety
func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.KeyMsg:
		if cmd := m.observeForSecrets(typed); cmd != nil {
			return m, cmd
		}

		switch typed.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.menuItems)-1 {
				m.cursor++
			}
		case "enter", " ":
			// send nav message, bo bubbletea tak lubi
			item := m.menuItems[m.cursor]
			return m, func() tea.Msg {
				return NavigateMsg{Target: item.Target}
			}
		}
	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
	}

	return m, nil
}

func (m *HomeModel) observeForSecrets(key tea.KeyMsg) tea.Cmd {
	str := strings.ToLower(key.String())

	// bail jeśli to control key, bo tam nie ma liter
	if len(str) > 1 {
		return nil
	}

	r := []rune(str)
	if len(r) != 1 || !unicode.IsLetter(r[0]) {
		return nil
	}

	m.secretBuffer += str
	if len(m.secretBuffer) > 16 {
		m.secretBuffer = m.secretBuffer[len(m.secretBuffer)-16:]
	}

	if strings.Contains(m.secretBuffer, "snake") || strings.Contains(m.secretBuffer, "games") {
		if !m.secretUnlocked {
			m.secretUnlocked = true
			m.secretMessage = "ok... arcade booted, powodzenia"
			m.lastUnlockPing = time.Now()
			m.menuItems = append(m.menuItems, m.secretItems...)
			return tea.Tick(420*time.Millisecond, func(time.Time) tea.Msg {
				return NavigateMsg{Target: ViewArcade}
			})
		}
	}

	return nil
}

// View rysuje ekran główny
func (m HomeModel) View() string {
	var b strings.Builder

	// ascii banner bo inaczej nudno
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

	// welcome, bo tak wypada
	welcome := "Welcome to pcstyle.dev SSH interface"
	b.WriteString(TitleStyle.Width(m.width).Render(welcome))
	b.WriteString("\n\n")

	// navigation menu aka główne decyzje
	for i, item := range m.menuItems {
		cursor := "  "
		if i == m.cursor {
			cursor = NavArrowStyle.Render("→ ")
		}

		// styling mix, nie pytaj
		itemStyle := NavItemStyle
		if i == m.cursor {
			itemStyle = NavItemSelectedStyle
		}

		title := itemStyle.Render(markSecretTitle(item))
		desc := HelpStyle.Render(item.Description)

		b.WriteString(fmt.Sprintf("%s%s - %s\n", cursor, title, desc))
	}

	// help + chaos
	b.WriteString("\n")
	helpText := "Use ↑/↓ or j/k to navigate • Enter to select • Press q to quit"
	b.WriteString(HelpStyle.Render(helpText))

	if m.secretUnlocked {
		b.WriteString("\n")
		secretLine := fmt.Sprintf("bonus: type 'snake' albo 'games' kiedyś tam. %s", m.secretMessage)
		b.WriteString(HelpStyle.Render(secretLine))
	}

	return BaseStyle.Render(b.String())
}

func markSecretTitle(item MenuItem) string {
	if item.isSecret {
		return item.Title + " *"
	}
	return item.Title
}
