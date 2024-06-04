package utilities

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestExtractQueryAndFindOptions(t *testing.T) {
	c := fiber.New()
	ctx := c.AcquireCtx(&fasthttp.RequestCtx{})

	t.Run("failure query", func(t *testing.T) {
		query, _, _, _ := ExtractQueryAndFindOptions(ctx)

		assert.Equal(t, primitive.M{}, query)
	})

	t.Run("failure findOptions", func(t *testing.T) {
		ctx.Locals("query", primitive.M{"field1": "value1"})

		_, findOptions, _, _ := ExtractQueryAndFindOptions(ctx)

		assert.Equal(t, &options.FindOptions{}, findOptions)
	})

	t.Run("failure page", func(t *testing.T) {
		ctx.Locals("query", primitive.M{"field1": "value1"})
		ctx.Locals("findOptions", options.Find().SetSort(primitive.M{"field2": 1}))

		_, _, page, limit := ExtractQueryAndFindOptions(ctx)

		assert.Equal(t, 0, page)
		assert.Equal(t, 0, limit)
	})

	t.Run("failure limit", func(t *testing.T) {
		ctx.Locals("query", primitive.M{"field1": "value1"})
		ctx.Locals("findOptions", options.Find().SetSort(primitive.M{"field2": 1}))
		ctx.Locals("page", 1)

		_, _, page, limit := ExtractQueryAndFindOptions(ctx)

		assert.Equal(t, 0, page)
		assert.Equal(t, 0, limit)
	})

	t.Run("success", func(t *testing.T) {
		ctx.Locals("query", primitive.M{"field1": "value1"})
		ctx.Locals("findOptions", options.Find().SetSort(primitive.M{"field2": 1}))
		ctx.Locals("page", 1)
		ctx.Locals("limit", 10)

		query, findOptions, page, limit := ExtractQueryAndFindOptions(ctx)

		assert.Equal(t, primitive.M{"field1": "value1"}, query)
		assert.Equal(t, options.Find().SetSort(primitive.M{"field2": 1}), findOptions)
		assert.Equal(t, 1, page)
		assert.Equal(t, 10, limit)

	})
}
