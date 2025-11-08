package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type arcadeState int

const (
	arcadeStateMenu arcadeState = iota
	arcadeStateSnake
	arcadeStateScreensaver
)

// ArcadeModel ogarnia hidden arcade, lowkey chaos
type ArcadeModel struct {
	state        arcadeState
	cursor       int
	menu         []arcadeEntry
	snake        *snakeGame
	statusLine   string
	lastBootPing time.Time
	width        int
	height       int
}

type arcadeEntry struct {
	title       string
	description string
	state       arcadeState
}

// newArcadeMenu zwraca entries, bo czemu nie
func newArcadeMenu() []arcadeEntry {
	return []arcadeEntry{
		{
			title:       "SNAKE.exe",
			description: "classic borderline laggy snake",
			state:       arcadeStateSnake,
		},
		{
			title:       "CRT DREAM",
			description: "just viby screensaver, bez kontroli",
			state:       arcadeStateScreensaver,
		},
	}
}

// NewArcadeModel odpala arcade view, jak stary emulator
func NewArcadeModel() ArcadeModel {
	return ArcadeModel{
		state:      arcadeStateMenu,
		menu:       newArcadeMenu(),
		statusLine: "booting tiny arcade... chwila",
	}
}

// Enter odpala mini boot sequence
func (m *ArcadeModel) Enter() tea.Cmd {
	m.state = arcadeStateMenu
	m.statusLine = "booting tiny arcade... chwila"
	m.lastBootPing = time.Now()
	return tea.Tick(350*time.Millisecond, func(time.Time) tea.Msg {
		return arcadeBootMsg(time.Now())
	})
}

// Update łapie eventy i wysyła dalej jak trzeba
func (m ArcadeModel) Update(msg tea.Msg) (ArcadeModel, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case arcadeStateMenu:
			return m.handleMenuKey(typed)
		case arcadeStateSnake:
			return m.forwardToSnake(typed)
		case arcadeStateScreensaver:
			if typed.String() == "esc" || typed.String() == "enter" {
				m.statusLine = "ok screensaver off, wracamy"
				m.state = arcadeStateMenu
			} else if typed.String() == "q" {
				return m, func() tea.Msg { return BackMsg{} }
			}
		}

	case arcadeBootMsg:
		if time.Since(m.lastBootPing) > 200*time.Millisecond {
			m.statusLine = "ok arcade ready, let's go"
		}

	case snakeTickMsg:
		if m.state == arcadeStateSnake && m.snake != nil {
			var cmd tea.Cmd
			m.snake, cmd = m.snake.updateTick()
			if !m.snake.alive {
				m.statusLine = "rip snake, press r żeby zrespawnić"
			}
			return m, cmd
		}
	}

	return m, nil
}

// View showtime, zależnie od state
func (m ArcadeModel) View() string {
	switch m.state {
	case arcadeStateSnake:
		return m.renderSnake()
	case arcadeStateScreensaver:
		return m.renderScreensaver()
	default:
		return m.renderMenu()
	}
}

func (m ArcadeModel) renderMenu() string {
	var b strings.Builder

	b.WriteString(drawArcadeBanner())
	b.WriteString("\n")
	b.WriteString(NavItemStyle.Render("mini arcade hub"))
	b.WriteString("\n\n")

	for i, entry := range m.menu {
		cursor := "  "
		if i == m.cursor {
			cursor = NavArrowStyle.Render("→ ")
		}

		itemStyle := NavItemStyle
		if i == m.cursor {
			itemStyle = NavItemSelectedStyle
		}

		line := fmt.Sprintf("%s%s - %s\n", cursor, itemStyle.Render(entry.title), HelpStyle.Render(entry.description))
		b.WriteString(line)
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("Enter to launch • esc/q żeby wrócić • r resets snake gdy padniesz"))
	if m.statusLine != "" {
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render(m.statusLine))
	}

	return BoxStyle.Render(b.String())
}

func (m ArcadeModel) renderSnake() string {
	if m.snake == nil {
		return BoxStyle.Render("snake not booted??? jak to w ogóle...")
	}

	board := m.snake.draw()
	lines := []string{
		TitleStyle.Render("SNAKE.exe // karma dla nostalgii"),
		board,
		HelpStyle.Render(fmt.Sprintf("score: %d  • steruj strzałkami / wasd • esc/q żeby wyjść • r respawn", m.snake.score)),
	}

	if !m.snake.alive {
		lines = append(lines, ErrorStyle.Render("padłeś. r = retry, esc = wyjdź"))
	}

	if m.statusLine != "" {
		lines = append(lines, HelpStyle.Render(m.statusLine))
	}

	return BoxStyle.Render(strings.Join(lines, "\n\n"))
}

func (m ArcadeModel) renderScreensaver() string {
	frame := drawScreensaver(time.Now())
	lines := []string{
		TitleStyle.Render("CRT DREAM // tak, kompletnie bezużyteczne"),
		frame,
		HelpStyle.Render("press esc lub enter, bo inaczej będziesz patrzeć na ten glitch forever"),
	}
	if m.statusLine != "" {
		lines = append(lines, HelpStyle.Render(m.statusLine))
	}
	return BoxStyle.Render(strings.Join(lines, "\n\n"))
}

func (m ArcadeModel) handleMenuKey(key tea.KeyMsg) (ArcadeModel, tea.Cmd) {
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		} else {
			m.cursor = len(m.menu) - 1
		}
	case "down", "j":
		if m.cursor < len(m.menu)-1 {
			m.cursor++
		} else {
			m.cursor = 0
		}
	case "enter", " ":
		entry := m.menu[m.cursor]
		switch entry.state {
		case arcadeStateSnake:
			m.snake = newSnakeGame()
			m.state = arcadeStateSnake
			m.statusLine = "snake loaded, nie crashuj w siebie pls"
			return m, m.snake.init()
		case arcadeStateScreensaver:
			m.state = arcadeStateScreensaver
			m.statusLine = "enjoy the glitch, chyba"
		}
	case "esc", "q":
		return m, func() tea.Msg { return BackMsg{} }
	}
	return m, nil
}

func (m ArcadeModel) forwardToSnake(key tea.KeyMsg) (ArcadeModel, tea.Cmd) {
	if m.snake == nil {
		m.snake = newSnakeGame()
		cmd := m.snake.init()
		return m, cmd
	}

	switch key.String() {
	case "esc", "q":
		m.state = arcadeStateMenu
		m.statusLine = "snake paused, pewnie głodny"
		return m, nil
	case "r":
		m.snake.reset()
		m.statusLine = "respawned, powodzenia elo"
		return m, m.snake.init()
	}

	var cmd tea.Cmd
	m.snake, cmd = m.snake.handleKey(key)
	return m, cmd
}

// drawArcadeBanner robi ascii bo czemu nie
func drawArcadeBanner() string {
	return `
╔═ ARC4D3 ═══════════════╗
║  █░█ █▀█ █▄░█ ▄▀█ ▄▀█  ║
║  ▀▄▀ █▀▀ █░▀█ █▀█ █▀█  ║
╚════════════════════════╝`
}

type arcadeBootMsg time.Time

// snake internals -----------------------------------------------------------

type snakeTickMsg struct{}

type snakeGame struct {
	width    int
	height   int
	snake    []snakePoint
	dir      snakePoint
	nextDir  snakePoint
	apple    snakePoint
	score    int
	alive    bool
	speed    time.Duration
	random   *rand.Rand
	lastMove time.Time
}

type snakePoint struct {
	x int
	y int
}

func newSnakeGame() *snakeGame {
	seed := time.Now().UnixNano()
	g := &snakeGame{
		width:  22,
		height: 12,
		snake: []snakePoint{
			{1, 1},
			{0, 1},
		},
		dir:     snakePoint{1, 0},
		nextDir: snakePoint{1, 0},
		alive:   true,
		speed:   180 * time.Millisecond,
		random:  rand.New(rand.NewSource(seed)),
	}
	g.spawnApple()
	return g
}

func (g *snakeGame) reset() {
	g.snake = []snakePoint{{1, 1}, {0, 1}}
	g.dir = snakePoint{1, 0}
	g.nextDir = g.dir
	g.alive = true
	g.score = 0
	g.spawnApple()
}

func (g *snakeGame) init() tea.Cmd {
	return tea.Tick(g.speed, func(time.Time) tea.Msg {
		return snakeTickMsg{}
	})
}

func (g *snakeGame) handleKey(key tea.KeyMsg) (*snakeGame, tea.Cmd) {
	switch key.String() {
	case "up", "k", "w":
		g.queueDir(0, -1)
	case "down", "j", "s":
		g.queueDir(0, 1)
	case "left", "h", "a":
		g.queueDir(-1, 0)
	case "right", "l", "d":
		g.queueDir(1, 0)
	}
	return g, nil
}

func (g *snakeGame) updateTick() (*snakeGame, tea.Cmd) {
	if !g.alive {
		return g, nil
	}

	g.dir = g.nextDir
	head := g.nextHead()

	if g.hitWall(head) || g.hitSelf(head) {
		g.alive = false
		return g, nil
	}

	g.snake = append([]snakePoint{head}, g.snake...)

	if head == g.apple {
		g.score += 10
		g.spawnApple()
		// przyspiesz trochę, bo czemu nie
		if g.speed > 80*time.Millisecond {
			g.speed -= 8 * time.Millisecond
		}
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}

	return g, g.init()
}

func (g *snakeGame) draw() string {
	var b strings.Builder
	borderTop := "╔" + strings.Repeat("═", g.width) + "╗"
	borderBottom := "╚" + strings.Repeat("═", g.width) + "╝"
	b.WriteString(borderTop + "\n")

	body := make(map[snakePoint]bool)
	for _, part := range g.snake {
		body[part] = true
	}

	for y := 0; y < g.height; y++ {
		b.WriteString("║")
		for x := 0; x < g.width; x++ {
			p := snakePoint{x: x, y: y}
			switch {
			case p == g.snake[0]:
				b.WriteString("■")
			case p == g.apple:
				b.WriteString("●")
			case body[p]:
				b.WriteString("░")
			default:
				if (x+y)%7 == 3 && g.random.Intn(40) == 0 {
					b.WriteString("·")
				} else {
					b.WriteString(" ")
				}
			}
		}
		b.WriteString("║\n")
	}

	b.WriteString(borderBottom)
	return b.String()
}

func (g *snakeGame) queueDir(dx, dy int) {
	if -dx == g.dir.x && -dy == g.dir.y {
		return // nie zawracamy w miejscu
	}
	g.nextDir = snakePoint{dx, dy}
}

func (g *snakeGame) nextHead() snakePoint {
	head := g.snake[0]
	return snakePoint{x: head.x + g.dir.x, y: head.y + g.dir.y}
}

func (g *snakeGame) hitWall(p snakePoint) bool {
	return p.x < 0 || p.x >= g.width || p.y < 0 || p.y >= g.height
}

func (g *snakeGame) hitSelf(p snakePoint) bool {
	for _, part := range g.snake {
		if part == p {
			return true
		}
	}
	return false
}

func (g *snakeGame) spawnApple() {
	for {
		p := snakePoint{
			x: g.random.Intn(g.width),
			y: g.random.Intn(g.height),
		}
		if !g.hitSelf(p) && p != g.snake[0] {
			g.apple = p
			return
		}
	}
}

// drawScreensaver generuje random ascii tv
func drawScreensaver(now time.Time) string {
	width := 30
	height := 10
	seed := now.UnixNano() / int64(1*time.Millisecond)
	r := rand.New(rand.NewSource(seed))

	var b strings.Builder
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			switch r.Intn(6) {
			case 0:
				b.WriteString("▒")
			case 1:
				b.WriteString("▓")
			case 2:
				b.WriteString("░")
			case 3:
				b.WriteString("┼")
			default:
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
