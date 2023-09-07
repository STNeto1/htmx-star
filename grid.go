package main

import (
	"math"
	"math/rand"
)

type GridElement struct {
	X         int
	Y         int
	F         float32
	G         float32
	H         float32
	Neighbors []*GridElement
	Previous  *GridElement
	Wall      bool
}

func (elem *GridElement) Compare(toCompare *GridElement) bool {
	return elem.X == toCompare.X && elem.Y == toCompare.Y
}

func (elem *GridElement) Label(grid *Container) string {
	if elem.Wall {
		return "wall"
	}

	if elem.Compare(grid.Start) {
		return "start"
	}

	if elem.Compare(grid.End) {
		return "end"
	}

	if grid.InPath(elem) {
		return "path"
	}

	for _, inner := range grid.OpenSet {
		if elem.Compare(inner) {
			return "open"
		}
	}

	for _, inner := range grid.ClosedSet {
		if elem.Compare(inner) {
			return "closed"
		}
	}

	return "void"
}

func (elem *GridElement) CalculareNeighbors(grid *Container) {

	if elem.X < grid.Cols-1 {
		elem.Neighbors = append(elem.Neighbors, grid.Grid[elem.X+1][elem.Y])
	}

	if elem.X > 0 {
		elem.Neighbors = append(elem.Neighbors, grid.Grid[elem.X-1][elem.Y])
	}

	if elem.Y < grid.Rows-1 {
		elem.Neighbors = append(elem.Neighbors, grid.Grid[elem.X][elem.Y+1])
	}

	if elem.Y > 0 {
		elem.Neighbors = append(elem.Neighbors, grid.Grid[elem.X][elem.Y-1])
	}

	// Diagonal
	// if elem.X > 0 && elem.Y > 0 {
	// 	elem.Neighbors = append(elem.Neighbors, grid.Grid[elem.X-1][elem.Y-1])
	// }
	//
	// if elem.X < grid.Cols-1 && elem.Y > 0 {
	// 	elem.Neighbors = append(elem.Neighbors, grid.Grid[elem.X+1][elem.Y-1])
	// }
	//
	// if elem.X > 0 && elem.Y < grid.Cols-1 {
	// 	elem.Neighbors = append(elem.Neighbors, grid.Grid[elem.X-1][elem.Y+1])
	// }
	//
	// if elem.X < grid.Rows-1 && elem.Y < grid.Cols-1 {
	// 	elem.Neighbors = append(elem.Neighbors, grid.Grid[elem.X+1][elem.Y+1])
	// }
}

func (elem *GridElement) UpdateHeuristics(to *GridElement) {
	elem.H = float32(math.Abs(float64(elem.X-to.X)) + math.Abs(float64(elem.Y-to.Y)))
}

func NewElement(x, y int) *GridElement {
	return &GridElement{
		X:         x,
		Y:         y,
		F:         0,
		G:         0,
		H:         0,
		Neighbors: []*GridElement{},
		Previous:  nil,
		Wall:      rand.Float64() < 0.3,
	}
}

type Container struct {
	Grid       [][]*GridElement // move to a pointer later?
	OpenSet    []*GridElement
	ClosedSet  []*GridElement
	Start      *GridElement
	End        *GridElement
	Path       []*GridElement
	Cols       int
	Rows       int
	NoSolution bool
	Finished   bool
}

func NewContainer(cols, rows int) *Container {
	grid := [][]*GridElement{}
	for i := 0; i < cols; i++ {
		row := []*GridElement{}

		for j := 0; j < rows; j++ {
			row = append(row, NewElement(i, j))
		}

		grid = append(grid, row)
	}

	start := grid[0][0]
	end := grid[cols-1][rows-1]

	start.Wall = false
	end.Wall = false

	openSet := []*GridElement{start}
	closedSet := []*GridElement{}

	container := &Container{
		Grid:      grid,
		OpenSet:   openSet,
		ClosedSet: closedSet,
		Start:     start,
		End:       end,
		Path:      []*GridElement{},
		Cols:      cols,
		Rows:      rows,
		Finished:  false,
	}

	for _, col := range container.Grid {
		for _, elem := range col {
			elem.CalculareNeighbors(container)
		}
	}

	return container
}

func (c *Container) Tick() {
	if c.Finished || c.NoSolution {
		return
	}

	if len(c.OpenSet) == 0 {
		c.NoSolution = true
		return
	}

	lowest := c.LowestInOpenSet()
	if lowest.Compare(c.End) {
		c.Backtrack()
		c.Finished = true
	}

	c.ClosedSet = append(c.ClosedSet, lowest)
	c.RemoveFromOpenSet(lowest)

	for _, neighbor := range lowest.Neighbors {
		if c.InClosedSet(neighbor) || neighbor.Wall {
			continue
		}

		tmpG := lowest.G + 1

		if c.InOpenSet(neighbor) {
			if tmpG < neighbor.G {
				// tmpG = neighbor.G
				neighbor.G = tmpG
			}
		} else {
			neighbor.G = tmpG
			c.OpenSet = append(c.OpenSet, neighbor)
		}

		if tmpG == neighbor.G {
			neighbor.UpdateHeuristics(c.End)
			neighbor.F = neighbor.G + neighbor.H
			neighbor.Previous = lowest
		}
	}

}

func (c *Container) RemoveFromOpenSet(toRemove *GridElement) {
	result := []*GridElement{}

	for _, elem := range c.OpenSet {
		if !elem.Compare(toRemove) {
			result = append(result, elem)
		}
	}

	c.OpenSet = result
}

func (c *Container) InClosedSet(toCheck *GridElement) bool {
	for _, elem := range c.ClosedSet {
		if elem.Compare(toCheck) {
			return true
		}

	}

	return false
}

func (c *Container) InOpenSet(toCheck *GridElement) bool {
	for _, elem := range c.OpenSet {
		if elem.Compare(toCheck) {
			return true
		}

	}

	return false
}

func (c *Container) InPath(toCheck *GridElement) bool {
	for _, elem := range c.Path {
		if elem.Compare(toCheck) {
			return true
		}

	}

	return false
}

func (c *Container) LowestInOpenSet() *GridElement {
	root := c.OpenSet[0]

	for _, elem := range c.OpenSet {
		if !root.Compare(elem) && elem.F < root.F {
			root = elem
		}
	}

	return root
}

func (c *Container) GetLabelGrid() [][]string {
	result := [][]string{}

	for _, col := range c.Grid {
		row := []string{}

		for _, elem := range col {
			row = append(row, elem.Label(c))
		}

		result = append(result, row)
	}

	return result
}

func (c *Container) Backtrack() {
	start := c.End
	for true {
		if start.Previous == nil {
			break
		}

		c.Path = append(c.Path, start.Previous)
		start = start.Previous
	}
}
