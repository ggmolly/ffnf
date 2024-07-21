package routes

import (
	"strconv"

	"github.com/ggmolly/ffnf/orm"
	"github.com/gofiber/fiber/v2"
)

// GET /api/v1/releases/:n
func GetLastReleases(c *fiber.Ctx) error {
	n, err := strconv.Atoi(c.Params("n"))
	if err != nil || n < 1 || n > 20 {
		return c.Status(fiber.StatusBadRequest).SendString("invalid range, must be between 1 and 20")
	}

	var releases []orm.Release
	if err := orm.GormDB.Find(&releases).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to get releases: " + err.Error())
	}

	return c.JSON(releases)
}
