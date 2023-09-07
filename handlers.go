package main

import (
	"github.com/gofiber/fiber/v2"
)

var container *Container

func handleRender(c *fiber.Ctx, layout string) error {
	return c.Render("index", fiber.Map{
		"rows":       createRange(0, container.Rows),
		"cols":       createRange(0, container.Cols),
		"labels":     container.GetLabelGrid(),
		"finished":   container.Finished,
		"noSolution": container.NoSolution,
		"show":       !container.NoSolution && !container.Finished,
	}, layout)
}

func HandleIndex(c *fiber.Ctx) error {
	if container == nil {
		container = NewContainer(15, 15)
	}

	return handleRender(c, "layouts/main")
}

func HandleTick(c *fiber.Ctx) error {
	// TODO - this is a hack to get the first tick to work
	if len(container.OpenSet) == 1 && len(container.ClosedSet) == 0 {
		container.Tick()
	}

	container.Tick()

	return handleRender(c, "")
}

func HandleReset(c *fiber.Ctx) error {
	container = NewContainer(10, 10)

	return handleRender(c, "")
}

func createRange(min, max int) []int {
	result := []int{}

	for i := min; i < max; i++ {
		result = append(result, i)
	}

	return result
}
