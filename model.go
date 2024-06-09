package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/KalebHawkins/automata/grid"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// State represents the state of a cell in the grid.
type State int

const (
	// Dead represents a dead cell.
	Dead State = iota
	// Alice represenets an alive cell.
	Alive
)

const (
	// AliveSymbol is the rune used to display cells in an Alive State
	AliveSymbol rune = '@'
	// DeadSymbol is the rune used to display cells in an Dead State
	DeadSymbol rune = '.'
	// FPS is how many frames per second
	FPS = 60
)

// keyMap defines a set of keybindings for the help menu.
// To work for help it must satisfy the key.Map interface.
type keyMap struct {
	Start key.Binding
	Help  key.Binding
	Next  key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Start, k.Next},
		{k.Help, k.Quit},
	}

}

var keys = keyMap{
	Start: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("[Space]", "start/stop animation"),
	),
	Next: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next generation"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help menu"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type updateMsg struct{}

// Model provides the structure for the simulation world.
type Model struct {
	*grid.Grid
	isPaused   bool
	generation int
	mouseLoc   [2]int
	keys       keyMap
	help       help.Model
}

// Init performs and model initialization.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles any update events for the model. Key events, mouse clicks, windows resizing etc.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleGridSizing(msg)
	case tea.KeyMsg:
		return m.handleKeyEvents(msg)
	case tea.MouseMsg:
		return m.handleMouseEvents(msg)
	case updateMsg:
		if !m.isPaused {
			m.updateCells()
			m.generation++
			return m, m.Tick()
		}
	}

	return m, nil
}

// View declares how to output the model.
func (m Model) View() string {
	var s strings.Builder

	// s.WriteString(m.Grid.Draw())

	for y := 0; y < m.Grid.Height(); y++ {
		for x := 0; x < m.Grid.Width(); x++ {

			cellState := m.Grid.Cell(x, y)

			switch cellState {
			case int(Alive): // Alive
				s.WriteRune(AliveSymbol)
			default: // Dead
				s.WriteRune(DeadSymbol)
			}
		}
		s.WriteRune('\n')
	}

	s.WriteString(fmt.Sprintf("\nMouse: (%d, %d)\n", m.mouseLoc[0], m.mouseLoc[1]))
	s.WriteString(fmt.Sprintf("\n\nGeneration: %d\n", m.generation))
	s.WriteString("This is working... Kinda?")
	return s.String()
}

// NewModel returns a new model.
func NewModel() Model {
	g := grid.NewGrid(0, 0)

	return Model{
		Grid:     g,
		isPaused: true,
		keys:     keys,
		help:     help.New(),
	}
}

// handleKeyEvents is a helper function for the models Update method.
func (m Model) handleKeyEvents(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	}

	switch msg.String() {
	case " ":
		m.isPaused = !m.isPaused
		return m, m.Tick()
	case "n":
		m.updateCells()
		m.generation++
	}
	return m, nil
}

// handleMouseEvents is a helper function for the models Update method.
func (m Model) handleMouseEvents(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	m.mouseLoc[0] = msg.X
	m.mouseLoc[1] = msg.Y

	switch msg.Button {
	case tea.MouseButtonLeft:
		// If the Cell is dead and clicked set it to an alive state.
		m.Grid.SetCell(msg.X, msg.Y, int(Alive))
	case tea.MouseButtonRight:
		// if the Cell is alive and right clicked... unalive it.
		m.Grid.SetCell(msg.X, msg.Y, int(Dead))
	}

	return m, nil
}

// handleGridSizing is a helper function for the models Update method.
// This handles resizing of the grid when the window size is changed.
func (m Model) handleGridSizing(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {

	// handle grid resizing
	newHeightWidth := min(msg.Width-msg.Width/4, msg.Height-msg.Height/4)
	m.Grid.Resize(newHeightWidth, newHeightWidth)

	// handle help resizing
	m.help.Width = msg.Width

	return m, nil
}

// countNeighbors counts the number of neighbors of a corresponding cell.
func (m Model) countNeighbors(x, y int) int {
	columns := m.Grid.Width()
	rows := m.Grid.Height()
	neighbors := 0

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			// Handle wrapping out of bounds index.
			neighborRow := (x + j + rows) % rows
			neighborCol := (y + i + columns) % columns
			neighbors += m.Grid.Cell(neighborRow, neighborCol)
		}
	}
	return neighbors
}

// updateCells gets each cells neighbor count and updates them accordingly
func (m Model) updateCells() {
	/*
		Live Cell

		- [x] Loneliness: If a live cell has fewer than 2 live neighbors, it dies in the next generation.
		- [x] Overcrowding: If a live cell has more than 3 live neighbors, it dies in the next generation.
		- [x] Stasis: If a live cell has 2 or 3 live neighbors, it remains alive in the next generation.

		Dead Cell

		- [x] Birth: If a dead cell has exactly 3 live neighbors, it becomes alive in the next generation.
	*/

	nextGen := make([]int, m.Grid.Width()*m.Grid.Height())

	for y := 0; y < m.Grid.Height(); y++ {
		for x := 0; x < m.Grid.Width(); x++ {
			nextState := m.getNextState(m.Grid.Cell(y, x), m.countNeighbors(y, x))
			nextGen[y*m.Grid.Width()+x] = nextState
		}
	}

	for y := 0; y < m.Grid.Height(); y++ {
		for x := 0; x < m.Grid.Width(); x++ {
			m.Grid.SetCell(y, x, nextGen[y*m.Grid.Width()+x])
		}
	}
}

// getNextState returns the next state of a cell in the grid.
func (m Model) getNextState(currentState, neighborCount int) int {
	switch {
	case neighborCount < 2 || neighborCount > 3 && currentState == int(Alive):
		return int(Dead)
	case neighborCount == 2 || neighborCount == 3 && currentState == int(Alive):
		return currentState
	case neighborCount == 3 && currentState == int(Dead):
		return int(Alive)
	default:
		return currentState
	}
}

// Tick plays the animation at the constant framerate of 60 FPS.
func (m Model) Tick() tea.Cmd {
	return tea.Tick(time.Second/FPS, func(t time.Time) tea.Msg {
		return updateMsg{}
	})
}
