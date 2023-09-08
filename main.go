package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", HandleIndex)
	app.Post("/tick", HandleTick)
	app.Post("/reset", HandleReset)
	app.Post("/finish", HandleFinish)

	app.Get("/setup", HandleSetup)
	app.Post("/setup", HandleSetup)

	app.Listen(":3000")
}
