package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/jet"
	"github.com/louissaadgo/web-chat/routes"
)

var port = ":3000"

func main() {
	engine := jet.New("./html", ".jet")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/home", routes.Home)

	app.Listen(port)
}
