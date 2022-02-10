package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Home(ctx *fiber.Ctx) error {
	return ctx.Render("main", fiber.Map{})
}
