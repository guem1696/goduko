package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Notes  key.Binding
	Quit   key.Binding
	Number key.Binding
	Help   key.Binding
	New    key.Binding
	Reset  key.Binding
	Delete key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Notes, k.New, k.Reset}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit, k.Notes, k.New, k.Reset},
		{k.Up, k.Down, k.Left, k.Right},
		{k.Number},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "w"),
		key.WithHelp("↑/w", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "s"),
		key.WithHelp("↓/s", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "a"),
		key.WithHelp("←/a", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "d"),
		key.WithHelp("→/d", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Notes: key.NewBinding(
		key.WithKeys("e", "n"),
		key.WithHelp("e", "toggle note mode"),
	),
	New: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "new sudoku"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "clear board"),
	),
	Number: key.NewBinding(
		key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
		key.WithHelp("1-9", "enter number"),
	),
	Delete: key.NewBinding(
		key.WithKeys("backspace", "c"),
		key.WithHelp("c/backspace", "clear cell"),
	),
}

var selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("#a89e32"))
var notesSelectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("#98c466"))
var lockedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#265aeb"))
var wrongStyle = lipgloss.NewStyle().Background(lipgloss.Color("#d47659"))
var highlighted = lipgloss.NewStyle().Background(lipgloss.Color("#424242"))
var highlightedAndWrong = lipgloss.NewStyle().Background(lipgloss.Color("#806d67"))

type Cell struct {
	number      int64
	filled      bool
	prefilled   bool
	selected    bool
	wrong       bool
	highlighted bool
	notes       map[int]bool
}

func (c Cell) getTextNumber() string {
	if c.filled || c.prefilled {
		switch c.number {
		case 0:
			res := ""
			res += "   \n"
			res += "   \n"
			res += "   \n"
			res += "   \n"
			res += "   "
			return res
		case 1:
			res := ""
			res += "  █\n"
			res += "  █\n"
			res += "  █\n"
			res += "  █\n"
			res += "  █"
			return res
		case 2:
			res := ""
			res += "███\n"
			res += "  █\n"
			res += "███\n"
			res += "█  \n"
			res += "███"
			return res
		case 3:
			res := ""
			res += "███\n"
			res += "  █\n"
			res += "███\n"
			res += "  █\n"
			res += "███"
			return res
		case 4:
			res := ""
			res += "█ █\n"
			res += "█ █\n"
			res += "███\n"
			res += "  █\n"
			res += "  █"
			return res
		case 5:
			res := ""
			res += "███\n"
			res += "█  \n"
			res += "███\n"
			res += "  █\n"
			res += "███"
			return res
		case 6:
			res := ""
			res += "███\n"
			res += "█  \n"
			res += "███\n"
			res += "█ █\n"
			res += "███"
			return res
		case 7:
			res := ""
			res += "███\n"
			res += "  █\n"
			res += "  █\n"
			res += "  █\n"
			res += "  █"
			return res
		case 8:
			res := ""
			res += "███\n"
			res += "█ █\n"
			res += "███\n"
			res += "█ █\n"
			res += "███"
			return res
		case 9:
			res := ""
			res += "███\n"
			res += "█ █\n"
			res += "███\n"
			res += "  █\n"
			res += "  █"
			return res
		}
	}

	res := [][]string{
		{" ", " ", " "},
		{" ", " ", " "},
		{" ", " ", " "},
		{" ", " ", " "},
		{" ", " ", " "},
	}

	for key, val := range c.notes {
		if val {
			x := (key - 1) % 3
			y := ((key - 1) / 3) + 1

			res[y][x] = fmt.Sprint(key)
		}
	}

	return strings.Join(Map(res, func(val []string) string { return strings.Join(val, "") }), "\n")
}

func (p Cell) Render(mode int) string {
	n := p.getTextNumber()

	selStyle := selectedStyle.Render

	if mode == NOTES {
		selStyle = notesSelectedStyle.Render
	}

	if p.selected {
		n = selStyle(n)
	} else if p.wrong && p.highlighted {
		n = highlightedAndWrong.Render(n)
	} else if p.highlighted {
		n = highlighted.Render(n)
	} else if p.wrong {
		n = wrongStyle.Render(n)
	}

	if p.prefilled {
		n = lockedStyle.Render(n)
	}

	return n
}

const (
	PLAYING = iota
	NOTES
)

type gameState struct {
	board [][]Cell
	state int
	x     int
	y     int
	done  bool
	mode  int
	help  help.Model
	keys keyMap
	windowWith int
}

func createCell(selected bool) Cell {
	return Cell{
		number:      0,
		filled:      false,
		prefilled:   false,
		selected:    selected,
		wrong:       false,
		highlighted: false,
		notes: map[int]bool{
			1: false,
			2: false,
			3: false,
			4: false,
			5: false,
			6: false,
			7: false,
			8: false,
			9: false,
		},
	}
}

func initialState() gameState {
	return gameState{
		board: [][]Cell{
			{createCell(true), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false)},
			{createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false)},
			{createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false)},
			{createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false)},
			{createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false)},
			{createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false)},
			{createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false)},
			{createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false)},
			{createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false), createCell(false)},
		},
		state: PLAYING,
		x:     0,
		y:     0,
		done:  false,
		mode:  PLAYING,
		help:  help.New(),
		keys: keys,
	}
}

func (s *gameState) generate() {
	s.reset()
	s.resetWrong()

	s.generateCell(0, true)

	for i := 0; i < 81; i++ {
		num := rand.Intn(81)
		x := num % 9
		y := num / 9

		oldNum := s.board[y][x].number
		s.board[y][x].number = 0
		s.board[y][x].prefilled = false

		if !s.generateCell(0, false) {
			s.board[y][x].number = oldNum
			s.board[y][x].prefilled = true
		}
	}
}

func (s *gameState) generateCell(i int, keepFilled bool) bool {
	x := i % 9
	y := i / 9

	if x > 8 || y > 8 {
		return true
	}
	if s.board[y][x].prefilled || s.board[y][x].filled {
		return s.generateCell(i+1, keepFilled)
	}

	numbers := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}

	rand.Shuffle(9, func(i, j int) { numbers[i], numbers[j] = numbers[j], numbers[i] })

	for n := 0; n < 9; n++ {
		num := numbers[n]

		s.board[y][x].number = num
		s.board[y][x].prefilled = true

		res := s.checkBoard(false)

		if res != WRONG {
			r := s.generateCell(i+1, keepFilled)
			if r {
				if !keepFilled {
					s.board[y][x].number = 0
					s.board[y][x].prefilled = false
				}
				return true
			}
		}

		s.board[y][x].number = 0
		s.board[y][x].prefilled = false

	}
	return false
}

func (s *gameState) reset() {
	for y := range s.board {
		for x := range s.board[y] {
			s.board[y][x] = createCell(false)
		}
	}
	s.board[0][0].selected = true
	s.x = 0
	s.y = 0
	s.resetHighlight()
	s.resetWrong()
}

func (s *gameState) resetWrong() {
	for y := range s.board {
		for x := range s.board[y] {
			s.board[y][x].wrong = false
		}
	}
}

func (s *gameState) resetHighlight() {
	for y := range s.board {
		for x := range s.board[y] {
			s.board[y][x].highlighted = false
			s.board[y][x].highlighted = false
		}
	}
}

type Status int

const (
	WRONG = iota
	OK
	DONE
)

func (s *gameState) checkBoard(markWrong bool) Status {
	s.resetWrong()

	res := DONE
	for i := 0; i < 9; i++ {
		if col := s.checkColumn(i); col != DONE {
			if col == WRONG {
				res = WRONG
				if markWrong {
					s.setWrongColumn(i)
				} else {
					return WRONG
				}
			} else if col == OK && res != WRONG {
				res = OK
			}
		}
		if row := s.checkRow(i); row != DONE {
			if row == WRONG {
				res = WRONG
				if markWrong {
					s.setWrongRow(i)
				} else {
					return WRONG
				}
			} else if row == OK && res != WRONG {
				res = OK
			}
		}
		if box := s.checkBox(i); box != DONE {
			if box == WRONG {
				res = WRONG
				if markWrong {
					s.setWrongBox(i)
				} else {
					return WRONG
				}
			} else if box == OK && res != WRONG {
				res = OK
			}
		}
	}

	return Status(res)
}

func (s gameState) checkRow(row int) Status {
	count := []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	for _, num := range s.board[row] {
		if num.filled || num.prefilled {
			count[num.number-1]++
		}

	}

	done := true

	for _, num := range count {
		if num > 1 {
			return WRONG
		}
		if num != 1 {
			done = false
		}
	}

	if done {
		return DONE
	}
	return OK
}

func (s gameState) checkColumn(column int) Status {
	count := []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	for _, row := range s.board {
		num := row[column]
		if num.filled || num.prefilled {
			count[num.number-1]++
		}

	}

	done := true

	for _, num := range count {
		if num > 1 {
			return WRONG
		}
		if num != 1 {
			done = false
		}
	}

	if done {
		return DONE
	}
	return OK
}

func (s gameState) checkBox(box int) Status {
	count := []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	row := (box / 3) * 3
	column := (box % 3) * 3

	//fmt.Printf("row: %v column: %v \n", row, column)

	for y := row; y < row+3; y++ {
		for x := column; x < column+3; x++ {
			if s.board[y][x].filled || s.board[y][x].prefilled {
				count[s.board[y][x].number-1]++
			}
		}
	}

	done := true

	for _, num := range count {
		if num > 1 {
			return WRONG
		}
		if num != 1 {
			done = false
		}
	}

	if done {
		return DONE
	}
	return OK
}

func (s *gameState) setWrongRow(row int) {
	for i := range s.board[row] {
		s.board[row][i].wrong = true
	}
}

func (s *gameState) setWrongColumn(column int) {
	for row := range s.board {
		s.board[row][column].wrong = true
	}
}

func (s *gameState) setWrongBox(box int) {
	row := (box / 3) * 3
	column := (box % 3) * 3

	for y := row; y < row+3; y++ {
		for x := column; x < column+3; x++ {
			s.board[y][x].wrong = true
		}
	}
}

func (s *gameState) setHighlight() {
	s.setHighlightedColumn()
	s.setHighlightedRow()
}

func (s *gameState) setHighlightedRow() {
	for i := range s.board[s.y] {
		s.board[s.y][i].highlighted = true
	}
}

func (s *gameState) setHighlightedColumn() {
	for row := range s.board {
		s.board[row][s.x].highlighted = true
	}
}

func (s gameState) Init() tea.Cmd {
	return nil
}

func (s gameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.windowWith = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return s, tea.Quit
		case key.Matches(msg, keys.Up):
			if s.y != 0 {
				newY := s.y - 1
				s.board[s.y][s.x].selected = false
				s.y = newY
				s.board[s.y][s.x].selected = true
			}
		case key.Matches(msg, keys.Down):
			if s.y != 8 {
				newY := s.y + 1
				s.board[s.y][s.x].selected = false
				s.y = newY
				s.board[s.y][s.x].selected = true
			}
		case key.Matches(msg, keys.Left):
			if s.x != 0 {
				newX := s.x - 1
				s.board[s.y][s.x].selected = false
				s.x = newX
				s.board[s.y][s.x].selected = true
			}
		case key.Matches(msg, keys.Right):
			if s.x != 8 {
				newX := s.x + 1
				s.board[s.y][s.x].selected = false
				s.x = newX
				s.board[s.y][s.x].selected = true
			}
		case key.Matches(msg, keys.Notes):
			if s.mode == NOTES {
				s.mode = PLAYING
			} else {
				s.mode = NOTES
			}
		case key.Matches(msg, keys.Number):
			if !s.board[s.y][s.x].prefilled {
				n, _ := strconv.ParseInt(msg.String(), 10, 64)
				if s.mode == NOTES {
					s.board[s.y][s.x].notes[int(n)] = !s.board[s.y][s.x].notes[int(n)]
				} else {
					s.board[s.y][s.x].number = n
					s.board[s.y][s.x].filled = true
				}
			}
		case key.Matches(msg, keys.Delete):
			if !s.board[s.y][s.x].prefilled {
				s.board[s.y][s.x].number = 0
				s.board[s.y][s.x].filled = false
			}
		case key.Matches(msg, keys.Reset):
			s.reset()
		case key.Matches(msg, keys.New):
			s.generate()
		case key.Matches(msg, keys.Help):
			s.help.ShowAll = !s.help.ShowAll
		}
	}
	s.resetHighlight()
	s.setHighlight()

	result := s.checkBoard(true)

	if result == DONE {
		s.done = true
	} else {
		s.done = false
	}

	return s, nil
}

func renderHorSeperator(height int, thick bool) string {
	res := ""

	for i := 0; i < height; i++ {
		if thick {
			res += "║"
		} else {
			res += " "
		}
		if i != height-1 {
			res += "\n"
		}
	}

	return res
}

func renderVerSeperator(cellWidth int, thick bool) string {
	char := " "
	if thick {
		char = "═"
	}

	thinVer := " "
	thickVer := "║"

	if thick {
		thinVer = char
		thickVer = "╬"
	}
	res := Join(Expand(strings.Join(Expand(char, cellWidth), ""), 9), thinVer)
	res[5] = thickVer
	res[11] = thickVer

	return strings.Join(res, "")
}

func renderRow(row []Cell, mode int) string {
	stringNumbers := Map(row, func(v Cell) string {
		return v.Render(mode)
	})
	stringNumbers = Join(stringNumbers, renderHorSeperator(lipgloss.Height(stringNumbers[0]), false))
	stringNumbers[5] = renderHorSeperator(lipgloss.Height(stringNumbers[0]), true)
	stringNumbers[11] = renderHorSeperator(lipgloss.Height(stringNumbers[0]), true)

	return lipgloss.JoinHorizontal(lipgloss.Top, stringNumbers...)
}

func (s gameState) View() string {
	rows := Map(s.board, func(val []Cell) string {
		return renderRow(val, s.mode)
	})

	rows = Join(rows, renderVerSeperator(3, false))
	rows[5] = renderVerSeperator(3, true)
	rows[11] = renderVerSeperator(3, true)

	board := lipgloss.JoinVertical(lipgloss.Left, rows...)

	s.help.Width = s.windowWith - lipgloss.Width(board)
	help := s.help.View(s.keys)
	return lipgloss.JoinHorizontal(lipgloss.Top, board, help)
}

func main() {
	rand.Seed(time.Now().Unix())
	godotenv.Load(".env")
	if os.Getenv("HELP_DEBUG") == "true" {
		if f, err := tea.LogToFile("debug.log", "help"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}
	}
	p := tea.NewProgram(initialState())
	if err := p.Start(); err != nil {
		fmt.Printf("Complete failure!")
		os.Exit(1)
	}
}
