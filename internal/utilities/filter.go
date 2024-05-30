package utilities

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ExtractEntityesFromFilter[V any](c *fiber.Ctx) V {
	v, ok := c.Locals("entities").(V)
	if !ok {
		return v
	}
	return v
}

func ExtractQueryAndFindOptions(c *fiber.Ctx) (primitive.M, *options.FindOptions, int, int) {
	query, ok := c.Locals("query").(primitive.M)
	if !ok {
		return nil, nil, 0, 0
	}

	findOptions, ok := c.Locals("findOptions").(*options.FindOptions)
	if !ok {
		return nil, nil, 0, 0
	}

	page, ok := c.Locals("page").(int)
	if !ok {
		return nil, nil, 0, 0
	}

	limit, ok := c.Locals("limit").(int)
	if !ok {
		return nil, nil, 0, 0
	}

	return query, findOptions, page, limit
}
