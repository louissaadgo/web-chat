package routes

import (
	"github.com/gofiber/fiber/v2"
)

//Home is the home page handler
func Home(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{})
}
