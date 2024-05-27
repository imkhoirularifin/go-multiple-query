package utilities

import "github.com/gofiber/fiber/v2"

func ExtractEntityesFromFilter[V any](c *fiber.Ctx) V {
	v, ok := c.Locals("entities").(V)
	if !ok {
		return v
	}
	return v
}
