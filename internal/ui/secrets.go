package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// SecretsModel is just vibes, trochę dziennik, trochę spam
type SecretsModel struct {
	index       int
	entries     []secretEntry
	statusFlash string
	lastTick    time.Time
	width       int
	height      int
}

type secretEntry struct {
	title string
	body  []string
}

// NewSecretsModel spawns the list
func NewSecretsModel() SecretsModel {
	return SecretsModel{
		entries: buildSecretEntries(),
	}
}

// Enter resets the view state
func (m *SecretsModel) Enter() tea.Cmd {
	m.index = time.Now().Nanosecond() % len(m.entries)
	m.statusFlash = "otwieram dziennik... chwila"
	return tea.Tick(280*time.Millisecond, func(time.Time) tea.Msg {
		return secretsBlinkMsg(time.Now())
	})
}

// Update handles events
func (m SecretsModel) Update(msg tea.Msg) (SecretsModel, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
	case tea.KeyMsg:
		switch typed.String() {
		case "left", "h", "k":
			m.index = wrapSecretsIndex(m.index-1, len(m.entries))
			m.statusFlash = "cofnąłem wpis • chill"
		case "right", "l", "j", " ":
			m.index = wrapSecretsIndex(m.index+1, len(m.entries))
			m.statusFlash = "next log... don't judge"
		case "enter":
			m.index = wrapSecretsIndex(m.index+1, len(m.entries))
			m.statusFlash = "skipping bo tak"
		case "esc", "q":
			return m, func() tea.Msg { return BackMsg{} }
		}
	case secretsBlinkMsg:
		m.statusFlash = ""
	}
	return m, nil
}

// View prints the current entry
func (m SecretsModel) View() string {
	entry := m.entries[m.index]
	var b strings.Builder

	header := fmt.Sprintf("log %02d :: %s", m.index+1, entry.title)
	b.WriteString(TitleStyle.Render(header))
	b.WriteString("\n\n")

	for _, line := range entry.body {
		b.WriteString(NavItemStyle.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render(" ←/→ (hjkl też) zmienia wpis • enter skip • esc wraca "))
	if m.statusFlash != "" {
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render(m.statusFlash))
	}

	return BoxStyle.Render(b.String())
}

func wrapSecretsIndex(idx, total int) int {
	if total == 0 {
		return 0
	}
	for idx < 0 {
		idx += total
	}
	if idx >= total {
		idx = idx % total
	}
	return idx
}

type secretsBlinkMsg time.Time

func buildSecretEntries() []secretEntry {
	return []secretEntry{
		{
			title: "ssh onboarding chaos",
			body: []string{
				"• lesson #1: ludzie kochają wejścia ascii, so keep them.",
				"• lesson #2: zostaw małe bugi, wyglądają na autentyczne.",
				"• note: tak, snake został napisany za późno w nocy.",
			},
		},
		{
			title: "todo? maybe?",
			body: []string{
				"- build bubble tea ui for fridge? czemu nie.",
				"- zakończyć shader o 3 w nocy. again.",
				"- find keyboard that nie hałasuje aż tak.",
			},
		},
		{
			title: "fave commands of the week",
			body: []string{
				"`curl wttr.in` // bo pogoda ma vibe",
				"`rg \"ugh\"` // sprawdzam gdzie narzekałem w kodzie",
				"`ssh` // obvious",
			},
		},
		{
			title: "konserwy audio",
			body: []string{
				"-> synthwave w tle, bo inaczej snake zasypia",
				"-> czasem white noise, seriously",
				"-> 3AM playlist: alt-J, nosowska, jak zwykle mix totalny",
			},
		},
		{
			title: "pcstyle lore dump",
			body: []string{
				"1. pierwsze portfolio było css zrobione w notatniku.",
				"2. potem generative art, bo czemu nie.",
				"3. teraz można mnie pingować przez ssh, wild.",
			},
		},
		{
			title: "easter egg roadmap",
			body: []string{
				"[ ] ascii art generator (jakieś glitchy logo).",
				"[x] snake ale taki w 2 kolorach.",
				"[ ] hidden chat bot? maybe.",
			},
		},
	}
}
