package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// MatrixModel represents the matrix rain effect
type MatrixModel struct {
	columns    []MatrixColumn
	width      int
	height     int
	startTime  time.Time
	duration   time.Duration
}

// MatrixColumn represents a single column in the matrix
type MatrixColumn struct {
	chars     []rune
	positions []int
	speeds    []int
	lengths   []int
}

// NewMatrixModel creates a new matrix rain effect
func NewMatrixModel() MatrixModel {
	return MatrixModel{
		startTime: time.Now(),
		duration:  10 * time.Second, // show for 10 seconds
	}
}

// Init initializes the matrix effect
func (m MatrixModel) Init() tea.Cmd {
	return m.tick()
}

// tickCmd sends a tick for animation
func (m MatrixModel) tick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(time.Time) tea.Msg {
		return MatrixTickMsg{}
	})
}

// MatrixTickMsg is sent for matrix animation
type MatrixTickMsg struct{}

// Update handles matrix updates
func (m MatrixModel) Update(msg tea.Msg) (MatrixModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter":
			return m, func() tea.Msg {
				return BackMsg{}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// initialize columns
		if m.width > 0 && m.height > 0 {
			m.columns = make([]MatrixColumn, m.width)
			for i := range m.columns {
				m.columns[i] = m.initColumn(m.height)
			}
		}
		return m, nil
	
	case MatrixTickMsg:
		// check if duration exceeded
		if time.Since(m.startTime) > m.duration {
			return m, func() tea.Msg {
				return BackMsg{}
			}
		}
		
		// update columns
		if len(m.columns) > 0 {
			for i := range m.columns {
				m.columns[i] = m.updateColumn(m.columns[i], m.height)
			}
		}
		
		return m, m.tick()
	}
	
	return m, nil
}

// initColumn creates a new matrix column
func (m MatrixModel) initColumn(height int) MatrixColumn {
	col := MatrixColumn{
		chars:     make([]rune, 0),
		positions: make([]int, 0),
		speeds:    make([]int, 0),
		lengths:   make([]int, 0),
	}
	
	// create 1-3 streams per column
	numStreams := rand.Intn(3) + 1
	for i := 0; i < numStreams; i++ {
		col.chars = append(col.chars, m.randomChar())
		col.positions = append(col.positions, rand.Intn(height))
		col.speeds = append(col.speeds, rand.Intn(2)+1)
		col.lengths = append(col.lengths, rand.Intn(8)+5)
	}
	
	return col
}

// updateColumn updates a single column
func (m MatrixModel) updateColumn(col MatrixColumn, height int) MatrixColumn {
	for i := range col.positions {
		col.positions[i] += col.speeds[i]
		
		// reset if off screen
		if col.positions[i] > height+col.lengths[i] {
			col.positions[i] = -col.lengths[i]
			col.chars[i] = m.randomChar()
		}
	}
	
	return col
}

// randomChar returns a random character (katakana, numbers, symbols)
func (m MatrixModel) randomChar() rune {
	chars := []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'ア', 'イ', 'ウ', 'エ', 'オ', 'カ', 'キ', 'ク', 'ケ', 'コ',
		'サ', 'シ', 'ス', 'セ', 'ソ', 'タ', 'チ', 'ツ', 'テ', 'ト',
		'ナ', 'ニ', 'ヌ', 'ネ', 'ノ', 'ハ', 'ヒ', 'フ', 'ヘ', 'ホ',
		'マ', 'ミ', 'ム', 'メ', 'モ', 'ヤ', 'ユ', 'ヨ',
		'ラ', 'リ', 'ル', 'レ', 'ロ', 'ワ', 'ヲ', 'ン',
	}
	return chars[rand.Intn(len(chars))]
}

// View renders the matrix effect
func (m MatrixModel) View() string {
	var b strings.Builder
	
	b.WriteString(TitleStyle.Render("MATRIX RAIN"))
	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("(press Esc/Q/Enter to exit)"))
	b.WriteString("\n\n")
	
	if len(m.columns) == 0 {
		b.WriteString(HelpStyle.Render("Resizing..."))
		return BoxStyle.Render(b.String())
	}
	
	// create screen buffer
	screen := make([][]rune, m.height)
	for i := range screen {
		screen[i] = make([]rune, m.width)
		for j := range screen[i] {
			screen[i][j] = ' '
		}
	}
	
	// draw columns
	for x, col := range m.columns {
		if x >= m.width {
			break
		}
		for i := range col.positions {
			startY := col.positions[i]
			length := col.lengths[i]
			
			for j := 0; j < length; j++ {
				y := startY + j
				if y >= 0 && y < m.height {
					// brighter at head, dimmer at tail
					char := col.chars[i]
					if j == 0 {
						screen[y][x] = char
					} else if j < length/2 {
						screen[y][x] = char
					} else {
						screen[y][x] = '·'
					}
				}
			}
		}
	}
	
	// render screen
	for _, row := range screen {
		for _, char := range row {
			b.WriteRune(char)
		}
		b.WriteString("\n")
	}
	
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render(fmt.Sprintf("Time remaining: %.1fs", m.duration.Seconds()-time.Since(m.startTime).Seconds())))
	
	return b.String()
}
