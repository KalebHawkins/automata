package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Status int

const (
	Initializing Status = iota
	Ready
)

const (
	AliveSymbol rune = '.'
	DeadSymbol  rune = '@'
)

type Model struct {
	Status
	*Grid
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleGridSizing(msg)
	case tea.KeyMsg:
		return m.handleKeyEvents(msg)
	case tea.MouseMsg:
		return m.handleMouseEvents(msg)
	}

	return m, nil
}

func (m Model) View() string {
	var s strings.Builder

	// s.WriteString(m.Grid.Draw())

	for y := 0; y < m.Grid.Height(); y++ {
		for x := 0; x < m.Grid.Width(); x++ {
			cellState := m.Grid.Cell(x, y)

			switch cellState {
			case Alive:
				s.WriteRune(AliveSymbol)
			default:
				s.WriteRune(DeadSymbol)
			}
		}
		s.WriteRune('\n')
	}

	s.WriteString("\n\nThis is working... Kinda?")
	return s.String()
}

func NewModel() Model {
	g := NewGrid(0, 0)

	return Model{
		Grid: g,
	}
}

func (m Model) handleKeyEvents(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) handleMouseEvents(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonLeft:
		// If the Cell is dead and clicked set it to an alive state.
		m.Grid.SetCell(msg.X, msg.Y, Alive)
	case tea.MouseButtonRight:
		// if the Cell is alive and right clicked... unalive it.
		m.Grid.SetCell(msg.X, msg.Y, Dead)
	}

	return m, nil
}

func (m Model) handleGridSizing(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.SetWidth(msg.Width - msg.Width/4)
	m.SetHeight(msg.Height - msg.Height/4)
	return m, nil
}
