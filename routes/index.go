package routes

import (
	"github.com/gofiber/fiber/v2"
)

// GET /
func Index(c *fiber.Ctx) error {
	return c.Render("pages/index", fiber.Map{
		"title": "FFNF - Index",
		"page":  "index",
	}, "layouts/main")
}
