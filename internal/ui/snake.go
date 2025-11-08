package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// SnakeGameModel represents the snake game state
type SnakeGameModel struct {
	snake      []Position
	food       Position
	direction  Direction
	score      int
	gameOver   bool
	paused     bool
	width      int
	height     int
	gameWidth  int
	gameHeight int
	speed      time.Duration
}

// Position represents a coordinate
type Position struct {
	X int
	Y int
}

// Direction represents snake direction
type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

// NewSnakeGameModel creates a new snake game
func NewSnakeGameModel() SnakeGameModel {
	rand.Seed(time.Now().UnixNano())
	
	// initial snake position - start in middle
	startX := 20
	startY := 10
	
	return SnakeGameModel{
		snake: []Position{
			{X: startX, Y: startY},
			{X: startX - 1, Y: startY},
			{X: startX - 2, Y: startY},
		},
		direction: DirRight,
		score:     0,
		gameOver:  false,
		paused:    false,
		gameWidth: 40,
		gameHeight: 20,
		speed:     150 * time.Millisecond, // start speed
	}
}

// Init initializes the snake game
func (m SnakeGameModel) Init() tea.Cmd {
	return tea.Batch(
		m.spawnFood(),
		m.tick(), // start game loop
	)
}

// tickCmd sends a tick message for game loop
func (m SnakeGameModel) tick() tea.Cmd {
	return tea.Tick(m.speed, func(time.Time) tea.Msg {
		return TickMsg{}
	})
}

// TickMsg is sent periodically for game updates
type TickMsg struct{}

// spawnFoodCmd spawns food at random position
func (m SnakeGameModel) spawnFood() tea.Cmd {
	return func() tea.Msg {
		// find empty spot
		for {
			food := Position{
				X: rand.Intn(m.gameWidth),
				Y: rand.Intn(m.gameHeight),
			}
			
			// check if food is on snake
			onSnake := false
			for _, seg := range m.snake {
				if seg.X == food.X && seg.Y == food.Y {
					onSnake = true
					break
				}
			}
			
			if !onSnake {
				return FoodSpawnedMsg{food}
			}
		}
	}
}

// FoodSpawnedMsg is sent when food spawns
type FoodSpawnedMsg struct {
	Food Position
}

// Update handles snake game updates
func (m SnakeGameModel) Update(msg tea.Msg) (SnakeGameModel, tea.Cmd) {
	if m.gameOver {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "r", "R":
				// restart game
				return NewSnakeGameModel(), NewSnakeGameModel().Init()
			case "esc", "q":
				return m, func() tea.Msg {
					return BackMsg{}
				}
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "w", "k":
			if m.direction != DirDown {
				m.direction = DirUp
			}
		case "down", "s", "j":
			if m.direction != DirUp {
				m.direction = DirDown
			}
		case "left", "a", "h":
			if m.direction != DirRight {
				m.direction = DirLeft
			}
		case "right", "d", "l":
			if m.direction != DirLeft {
				m.direction = DirRight
			}
		case " ":
			m.paused = !m.paused
		case "esc", "q":
			return m, func() tea.Msg {
				return BackMsg{}
			}
		}

	case FoodSpawnedMsg:
		m.food = msg.Food
		return m, nil

	case TickMsg:
		// game loop tick - move snake
		if !m.paused && !m.gameOver {
			// calculate new head position
			head := m.snake[0]
			newHead := Position{X: head.X, Y: head.Y}
			
			switch m.direction {
			case DirUp:
				newHead.Y--
			case DirDown:
				newHead.Y++
			case DirLeft:
				newHead.X--
			case DirRight:
				newHead.X++
			}
			
			// check wall collision
			if newHead.X < 0 || newHead.X >= m.gameWidth ||
			   newHead.Y < 0 || newHead.Y >= m.gameHeight {
				m.gameOver = true
				return m, nil
			}
			
			// check self collision
			for _, seg := range m.snake {
				if seg.X == newHead.X && seg.Y == newHead.Y {
					m.gameOver = true
					return m, nil
				}
			}
			
			// add new head
			m.snake = append([]Position{newHead}, m.snake...)
			
			// check food collision
			if newHead.X == m.food.X && newHead.Y == m.food.Y {
				m.score++
				// increase speed slightly
				if m.speed > 50*time.Millisecond {
					m.speed -= 2 * time.Millisecond
				}
				// spawn new food
				var cmd tea.Cmd
				m, cmd = m, m.spawnFood()
				return m, tea.Batch(cmd, m.tick())
			} else {
				// remove tail if no food eaten
				m.snake = m.snake[:len(m.snake)-1]
			}
			
			return m, m.tick()
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// View renders the snake game
func (m SnakeGameModel) View() string {
	var b strings.Builder
	
	// title
	b.WriteString(TitleStyle.Render("🐍 SNAKE GAME 🐍"))
	b.WriteString("\n\n")
	
	if m.gameOver {
		b.WriteString(ErrorStyle.Render("GAME OVER!"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("Final Score: %d\n", m.score))
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("Press R to restart • Esc/Q to quit"))
		return BoxStyle.Render(b.String())
	}
	
	if m.paused {
		b.WriteString(HelpStyle.Render("⏸ PAUSED - Press Space to resume"))
		b.WriteString("\n\n")
	}
	
	// score
	b.WriteString(fmt.Sprintf("Score: %d\n", m.score))
	b.WriteString("\n")
	
	// game board
	board := make([][]rune, m.gameHeight)
	for i := range board {
		board[i] = make([]rune, m.gameWidth)
		for j := range board[i] {
			board[i][j] = ' '
		}
	}
	
	// draw snake
	for i, seg := range m.snake {
		if i == 0 {
			board[seg.Y][seg.X] = '●' // head
		} else {
			board[seg.Y][seg.X] = '○' // body
		}
	}
	
	// draw food
	board[m.food.Y][m.food.X] = '*' // food
	
	// render board with borders
	b.WriteString("┌")
	for i := 0; i < m.gameWidth; i++ {
		b.WriteString("─")
	}
	b.WriteString("┐\n")
	
	for _, row := range board {
		b.WriteString("│")
		for _, cell := range row {
			b.WriteRune(cell)
		}
		b.WriteString("│\n")
	}
	
	b.WriteString("└")
	for i := 0; i < m.gameWidth; i++ {
		b.WriteString("─")
	}
	b.WriteString("┘\n")
	
	// controls
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("WASD/↑↓←→/HJKL to move • Space to pause • Esc/Q to quit"))
	
	return BoxStyle.Render(b.String())
}
