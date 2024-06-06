package main

import (
	"fmt"
)

// State represents the state of a cell in the grid.
type State int

const (
	// Dead represents a dead cell.
	Dead State = iota
	// Alice represenets an alive cell.
	Alive
)

// Grid represent a 2 deminsonal grid structure.
type Grid struct {
	width  int
	height int
	cells  [][]State
}

func (g *Grid) Width() int  { return g.width }
func (g *Grid) Height() int { return g.height }

func (g *Grid) SetWidth(width int) error {
	if g.width == width {
		return nil
	}

	// Handle increasing the grid width.
	if width > g.width {
		for y := range g.cells {
			g.cells[y] = append(g.cells[y], make([]State, width-g.width)...)
		}
	} else { // Handle decreasing the grid width.
		if width < 0 {
			return fmt.Errorf("grid width cannot be negative")
		}
		for i := range g.cells {
			g.cells[i] = g.cells[i][:width]
		}
	}

	g.width = width

	return nil
}

func (g *Grid) SetHeight(height int) error {
	if g.height == height {
		return nil
	}

	// Handle increasing the grid height.
	if height > g.height {
		for i := g.height; i < height; i++ {
			g.cells = append(g.cells, make([]State, g.width))
		}
	} else { // Handle decreasing the grid height.
		if height < 0 {
			return fmt.Errorf("grid height cannot be negative")
		}
		g.cells = g.cells[:height]
	}

	g.height = height

	return nil
}

func (g *Grid) Cell(x, y int) State {
	if x < 0 || x >= g.Width() || y < 0 || y >= g.Height() {
		return -1
	}
	return g.cells[y][x]
}
func (g *Grid) SetCell(x, y int, v State) {
	if x < 0 || x >= g.Width() || y < 0 || y >= g.Height() {
		return
	}
	g.cells[y][x] = v
}

// Draw for debugging purposes.
// func (g *Grid) Draw() string {
// 	var s strings.Builder
// 	for y := range g.cells {
// 		for x := range g.cells[y] {
// 			s.WriteString(fmt.Sprintf("%d", g.Cell(x, y)))
// 		}
// 		s.WriteRune('\n')
// 	}

// 	return s.String()
// }

func NewGrid(w, h int) *Grid {
	cells := make([][]State, h)
	for i := range cells {
		cells[i] = make([]State, w)
	}

	return &Grid{
		width:  w,
		height: h,
		cells:  cells,
	}
}
