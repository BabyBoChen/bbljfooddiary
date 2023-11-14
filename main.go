package main

import (
	"fmt"
	"log"

	"github.com/BabyBoChen/bbljfooddiary/models"
	"github.com/BabyBoChen/bbljfooddiary/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
)

var sessionStore *session.Store

func index(c *fiber.Ctx) error {
	viewModel := models.NewIndexViewModel()
	return c.Render("index", viewModel)
}

func newCuisine(c *fiber.Ctx) error {
	return c.Render("newCuisine", nil)
}

func main() {
	envVars := services.ReadEnvironmentVariables()
	engine := html.New("./views", ".html")
	// Reload the templates on each render, good for development
	engine.Reload(envVars.Config == "debug") // Optional. Default: false
	// Debug will print each template that is parsed, good for debugging
	engine.Debug(envVars.Config == "debug") // Optional. Default: false

	sessionStore = session.New()

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./wwwroot")

	app.Get("/", index)
	app.Get("/newCuisine", newCuisine)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", envVars.Port)))
}
