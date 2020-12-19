package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := tea.NewProgram(initialModel()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	fmt.Print("Ba-bam\n")
}

const size = 15

type cell struct {
	x, y uint8
}

type direction uint8

const (
	up = direction(iota)
	down
	left
	right
)

type model struct {
	snake      []cell
	dir        direction
	dirChanged bool // dir was changed, but no affected snake position yet
	food       cell
}

func initialModel() model {
	m := model{
		snake: []cell{{x: size / 2, y: size / 2}},
		dir:   right,
	}
	m.food = rndFood(m)
	return m
}

type move struct{}

func (m model) Init() tea.Cmd {
	return tea.Tick(500*time.Microsecond, func(t time.Time) tea.Msg {
		return move{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case move:
		snake := make([]cell, 1, len(m.snake))
		switch m.dir {
		case up:
			y := m.snake[0].y
			if y == 0 {
				y = size - 1
			} else {
				y--
			}
			snake[0] = cell{x: m.snake[0].x, y: y}
		case down:
			y := m.snake[0].y + 1
			if y >= size {
				y = 0
			}
			snake[0] = cell{x: m.snake[0].x, y: y}
		case left:
			x := m.snake[0].x
			if x == 0 {
				x = size - 1
			} else {
				x--
			}
			snake[0] = cell{x: x, y: m.snake[0].y}
		case right:
			x := m.snake[0].x + 1
			if x >= size {
				x = 0
			}
			snake[0] = cell{x: x, y: m.snake[0].y}
		}
		if cellIn(snake[0], m.snake) {
			return m, tea.Quit
		}
		if snake[0] == m.food {
			m.snake = append(snake, m.snake[0:len(m.snake)]...)
			m.food = rndFood(m)
		} else {
			m.snake = append(snake, m.snake[0:len(m.snake)-1]...)
		}
		m.dirChanged = false
		return m, tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
			return move{}
		})
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		// Cycle between inputs
		case "up", "k":
			if m.dir != down && !m.dirChanged {
				m.dirChanged = true
				m.dir = up
			}
			return m, nil
		// Cycle between inputs
		case "down", "j":
			if m.dir != up && !m.dirChanged {
				m.dirChanged = true
				m.dir = down
			}
			return m, nil
		// Cycle between inputs
		case "left", "h":
			if m.dir != right && !m.dirChanged {
				m.dirChanged = true
				m.dir = left
			}
			return m, nil
		// Cycle between inputs
		case "right", "l":
			if m.dir != left && !m.dirChanged {
				m.dirChanged = true
				m.dir = right
			}
			return m, nil
		}
	}

	return m, cmd
}

func (m model) View() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("╔═══════════════╗\n")
	for i := uint8(0); i < size; i++ {
		buf.WriteRune('║')
		for j := uint8(0); j < size; j++ {
			c := cell{x: j, y: i}
			if c == m.food {
				buf.WriteRune('F')
			} else if cellIn(c, m.snake) {
				buf.WriteRune('X')
			} else {
				buf.WriteRune(' ')
			}
		}
		buf.WriteString("║\n")
	}
	buf.WriteString("╚═══════════════╝\n")
	return buf.String()
}

// return true if c in arr
func cellIn(c cell, arr []cell) bool {
	for _, c2 := range arr {
		if c == c2 {
			return true
		}
	}
	return false
}

func rndFood(m model) cell {
	freeCells := make([]cell, 0, size*size)
	for i := uint8(0); i < size; i++ {
		for j := uint8(0); j < size; j++ {
			c := cell{x: j, y: i}
			if cellIn(c, m.snake) {
				continue
			}
			freeCells = append(freeCells, c)
		}
	}
	return freeCells[rand.Intn(len(freeCells))]
}
