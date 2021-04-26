package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
	"github.com/louissaadgo/web-chat/routes"
)

var port = ":3000"

func main() {
	engine := html.New("./html", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	go routes.ListenToWsChannel()
	app.Get("/home", routes.Home)
	app.Get("/ws", websocket.New(routes.Ws))

	app.Listen(port)
}
