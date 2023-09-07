package main

import "github.com/gofiber/fiber/v2"

var container *Container

func HandleIndex(c *fiber.Ctx) error {
	if container == nil {
		container = NewContainer(10, 10)
	}

	return c.Render("index", fiber.Map{
		"rows":       createRange(0, container.Rows),
		"cols":       createRange(0, container.Cols),
		"labels":     container.GetLabelGrid(),
		"finished":   container.Finished,
		"noSolution": container.NoSolution,
	}, "layouts/main")

}

func createRange(min, max int) []int {
	result := []int{}

	for i := min; i < max; i++ {
		result = append(result, i)
	}

	return result
}
