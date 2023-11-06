package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
)

var sessionStore *session.Store

func index(c *fiber.Ctx) error {
	viewModel := make(fiber.Map)
	viewModel["Title"] = "Hello, World!"
	return c.Render("index", viewModel)
}

func main() {
	engine := html.New("./views", ".html")
	// Reload the templates on each render, good for development
	engine.Reload(true) // Optional. Default: false
	// Debug will print each template that is parsed, good for debugging
	engine.Debug(true) // Optional. Default: false

	sessionStore = session.New()

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./wwwroot")

	app.Get("/", index)

	log.Fatal(app.Listen(":3000"))
}
