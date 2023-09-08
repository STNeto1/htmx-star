package main

import (
	"net/http"

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
		"show":       container.ShouldContinue(),
	}, layout)
}

func HandleIndex(c *fiber.Ctx) error {
	if container == nil {
		return c.Redirect("/setup")
	}

	return handleRender(c, "layouts/main")
}

type SetupBody struct {
	Rows int `form:"rows" json:"rows"`
	Cols int `form:"cols" json:"cols"`
}

func HandleSetup(c *fiber.Ctx) error {
	if container != nil {
		return c.Redirect("/")
	}

	if c.Method() == "POST" {
		body := new(SetupBody)

		if err := c.BodyParser(body); err != nil {
			// return error
		}

		container = NewContainer(body.Cols, body.Rows)
		return c.Redirect("/", http.StatusTemporaryRedirect)
		// return somethiong
	}

	return c.Render("setup", fiber.Map{}, "layouts/main")
}

func HandleTick(c *fiber.Ctx) error {
	if container == nil {
		return c.Redirect("/", http.StatusTemporaryRedirect)
	}

	// TODO - this is a hack to get the first tick to work
	if len(container.OpenSet) == 1 && len(container.ClosedSet) == 0 {
		container.Tick()
	}

	container.Tick()

	return handleRender(c, "")
}

func HandleReset(c *fiber.Ctx) error {
	if container == nil {
		return c.Redirect("/setup")
	}

	container = NewContainer(container.Cols, container.Rows)

	return handleRender(c, "")
}

func HandleFinish(c *fiber.Ctx) error {
	if container == nil {
		return c.Redirect("/setup")
	}

	for true {
		if !container.ShouldContinue() {
			break
		}

		container.Tick()
	}

	return handleRender(c, "")
}

func createRange(min, max int) []int {
	result := []int{}

	for i := min; i < max; i++ {
		result = append(result, i)
	}

	return result
}
