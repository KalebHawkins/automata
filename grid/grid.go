package grid

// Grid represent a 2 deminsonal grid structure.
type Grid struct {
	width  int
	height int
	cells  []int
}

// Width returns the width of the grid.
func (g *Grid) Width() int { return g.width }

// Height returns the height of the grid.
func (g *Grid) Height() int { return g.height }

// SetWidth will set the width of the grid. Panics if the width is set to < 0.
func (g *Grid) SetWidth(width int) {
	g.Resize(width, g.height)
}

// SetHeight will set the height of the grid. Panics if height is set to < 0
func (g *Grid) SetHeight(height int) {
	g.Resize(g.width, height)
}

// Resize handles the logic for resizing the grid.
func (g *Grid) Resize(width, height int) {
	if width < 0 || height < 0 {
		panic("width and height cannot be negative")
	}

	newCells := make([]int, width*height)
	minWidth := min(g.width, width)
	minHeight := min(g.height, height)

	// Copy existing cells
	for y := 0; y < minHeight; y++ {
		copy(newCells[y*width:], g.cells[y*g.width:y*g.width+minWidth])
	}

	// Fill new rows/columns with zeros
	for y := minHeight; y < height; y++ {
		for x := 0; x < width; x++ {
			newCells[y*width+x] = 0
		}
	}

	g.cells = newCells
	g.width = width
	g.height = height
}

// Cell returns a specified cell at the (x, y) location.
func (g *Grid) Cell(x, y int) int {
	if x < 0 || x >= g.Width() || y < 0 || y >= g.Height() {
		return 0
	}
	return g.cells[y*g.width+x]
}

// SetCell sets the value (v) of the x, y location.
func (g *Grid) SetCell(x, y int, v int) {
	if x < 0 || x >= g.Width() || y < 0 || y >= g.Height() {
		return
	}
	g.cells[y*g.width+x] = v
}

// Cells returns the entire grid of cells.
func (g *Grid) Cells() []int {
	return g.cells
}

// SetCells takes a new grid and replaces the current grid's cells.
func (g *Grid) SetCells(newCells []int) {
	g.cells = newCells
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

// NewGrid generates a fresh empty grid using the width (w) and height (h) specified.
func NewGrid(w, h int) *Grid {
	cells := make([]int, w*h)

	return &Grid{
		width:  w,
		height: h,
		cells:  cells,
	}
}
