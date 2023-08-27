package main

import (
	"errors"
	"fmt"
)

type GridElement struct {
	X int
	Y int
	F float32
	G float32
	H float32
}

func (elem *GridElement) Compare(toCompare *GridElement) bool {
	return elem.X == toCompare.X && elem.Y == toCompare.Y
}

func (elem *GridElement) Label(grid *Container) string {
	for _, inner := range grid.OpenSet {
		if elem.Compare(inner) {
			return "O"
		}
	}

	for _, inner := range grid.ClosedSet {
		if elem.Compare(inner) {
			return "C"
		}
	}

	return "G"
}

func NewElement(x, y int) *GridElement {
	return &GridElement{
		X: x,
		Y: y,
		F: 0,
		G: 0,
		H: 0,
	}
}

type Container struct {
	Grid      [][]*GridElement // move to a pointer later?
	OpenSet   []*GridElement
	ClosedSet []*GridElement
	Start     *GridElement
	End       *GridElement
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

	openSet := []*GridElement{start}
	closedSet := make([]*GridElement, 0)

	return &Container{
		Grid:      grid,
		OpenSet:   openSet,
		ClosedSet: closedSet,
		Start:     start,
		End:       end,
	}
}

func (c *Container) Tick() error {
	if len(c.ClosedSet) == 0 {
		return errors.New("no solution")
	}

	return nil
}

func (c *Container) Print() {
	for _, col := range c.Grid {
		for _, row := range col {
			fmt.Printf("%s ", row.Label(c))
		}
		fmt.Println()
	}
}
