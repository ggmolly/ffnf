package routes

import (
	"strconv"

	"github.com/ggmolly/ffnf/orm"
	"github.com/gofiber/fiber/v2"
)

func GetNoticesAfter(c *fiber.Ctx) error {
	var releases []orm.Notice

	lastId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}

	if err := orm.GormDB.Where("id > ?", lastId).Order("id ASC").Find(&releases).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to get releases: " + err.Error())
	}

	return c.JSON(releases)
}
