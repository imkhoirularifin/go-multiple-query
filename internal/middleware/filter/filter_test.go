package filter

import (
	"encoding/json"
	"errors"
	"go-multiple-query/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestParseQueriesToFilters(t *testing.T) {
	t.Run("ParseQueriesToFilters Success", func(t *testing.T) {
		queries := map[string]string{
			"field1.equal":       "value1",
			"field2.notEqual":    "value2",
			"field3.greaterThan": "10",
			"field4.lessThan":    "20",
			"field5.in":          "value3,value4",
			"page":               "1",
		}

		_, err := parseQueriesToFilters(queries)
		assert.NoError(t, err)
	})

	t.Run("ParseQueriesToFilters Failed", func(t *testing.T) {
		queries := map[string]string{
			"field1.": "value1",
			"field2":  "value2",
		}

		_, err := parseQueriesToFilters(queries)
		assert.Error(t, err)
	})

	t.Run("ParseQueriesToFilters Failed2", func(t *testing.T) {
		queries := map[string]string{
			"field1.someoneyoulove": "value1",
		}

		_, err := parseQueriesToFilters(queries)
		assert.Error(t, err)
	})

}

func TestBuildQuery(t *testing.T) {
	t.Run("BuildQuery with valid filters", func(t *testing.T) {
		request := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "field1",
					Filter: domain.Equal,
					Value:  []string{"value1"},
				},
				{
					Key:    "field2",
					Filter: domain.NotEqual,
					Value:  []string{"value2"},
				},
				{
					Key:    "field3",
					Filter: domain.GreaterThanOrEqual,
					Value:  []string{"10"},
				},
				{
					Key:    "field4",
					Filter: domain.LessThanOrEqual,
					Value:  []string{"20"},
				},
				{
					Key:    "field5",
					Filter: domain.GreaterThan,
					Value:  []string{"30"},
				},
				{
					Key:    "field6",
					Filter: domain.LessThan,
					Value:  []string{"40"},
				},
				{
					Key:    "field7",
					Filter: domain.In,
					Value:  []string{"value3", "value4"},
				},
				{
					Key:    "field8",
					Filter: domain.In,
					Value:  []string{"50", "60"},
				},
				{
					Key:    "field9",
					Filter: domain.Equal,
					Value:  []string{"70"},
				},
				{
					Key:    "field10",
					Filter: domain.NotEqual,
					Value:  []string{"80"},
				},
			},
		}

		expectedQuery := primitive.M{
			"field1":  primitive.M{"$eq": "value1"},
			"field2":  primitive.M{"$ne": "value2"},
			"field3":  primitive.M{"$gte": int64(10)},
			"field4":  primitive.M{"$lte": int64(20)},
			"field5":  primitive.M{"$gt": int64(30)},
			"field6":  primitive.M{"$lt": int64(40)},
			"field7":  primitive.M{"$in": []string{"value3", "value4"}},
			"field8":  primitive.M{"$in": []int64{50, 60}},
			"field9":  primitive.M{"$eq": int64(70)},
			"field10": primitive.M{"$ne": int64(80)},
		}

		query, err := buildQuery(request)
		assert.NoError(t, err)
		assert.Equal(t, expectedQuery, query)
	})

	t.Run("BuildQuery with invalid filters", func(t *testing.T) {
		request1 := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "field1",
					Filter: domain.GreaterThanOrEqual,
					Value:  []string{"value1"},
				},
			},
		}

		_, err := buildQuery(request1)
		assert.Error(t, err)

		request2 := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "field2",
					Filter: domain.LessThanOrEqual,
					Value:  []string{"value2"},
				},
			},
		}

		_, err = buildQuery(request2)
		assert.Error(t, err)

		request3 := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "field3",
					Filter: domain.GreaterThan,
					Value:  []string{"value3"},
				},
			},
		}

		_, err = buildQuery(request3)
		assert.Error(t, err)

		request4 := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "field4",
					Filter: domain.LessThan,
					Value:  []string{"value4"},
				},
			},
		}

		_, err = buildQuery(request4)
		assert.Error(t, err)
	})

	t.Run("BuildQuery with invalid ID", func(t *testing.T) {
		requestEqual := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "id",
					Filter: domain.Equal,
					Value:  []string{"invalidID"},
				},
			},
		}

		_, err := buildQuery(requestEqual)
		assert.Error(t, err)

		requestNotEqual := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "id",
					Filter: domain.NotEqual,
					Value:  []string{"invalidID"},
				},
			},
		}

		_, err = buildQuery(requestNotEqual)
		assert.Error(t, err)

		requestIn := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "id",
					Filter: domain.In,
					Value:  []string{"invalidID"},
				},
			},
		}

		_, err = buildQuery(requestIn)
		assert.Error(t, err)
	})

	t.Run("BuildQuery with valid ID", func(t *testing.T) {
		requestEqual := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "id",
					Filter: domain.Equal,
					Value:  []string{"60f1b0b3d1f3f3b3b3b3b3b3"},
				},
			},
		}

		_, err := buildQuery(requestEqual)
		assert.NoError(t, err)

		requestNotEqual := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "id",
					Filter: domain.NotEqual,
					Value:  []string{"60f1b0b3d1f3f3b3b3b3b3b3"},
				},
			},
		}

		_, err = buildQuery(requestNotEqual)
		assert.NoError(t, err)

		requestIn := domain.QueryRequest{
			Filters: []domain.QueryFilter{
				{
					Key:    "id",
					Filter: domain.In,
					Value:  []string{"60f1b0b3d1f3f3b3b3b3b3b3"},
				},
			},
		}

		_, err = buildQuery(requestIn)
		assert.NoError(t, err)
	})
}

func TestParseQueryOrder(t *testing.T) {
	t.Run("ParseQueryOrder Success", func(t *testing.T) {
		queries := map[string]string{
			"orderBy":   "field1",
			"sortOrder": "asc",
		}

		orderBy, sortOrder, err := parseQueryOrder(queries)
		assert.NoError(t, err)
		assert.Equal(t, "field1", orderBy)
		assert.Equal(t, "asc", sortOrder)

		queries2 := map[string]string{
			"orderBy": "field1",
		}

		orderBy, sortOrder, err = parseQueryOrder(queries2)
		assert.NoError(t, err)
		assert.Equal(t, "field1", orderBy)
		assert.Equal(t, "asc", sortOrder)
	})

	t.Run("ParseQueryOrder Failed", func(t *testing.T) {
		queries1 := map[string]string{
			"orderBy":   "field1",
			"sortOrder": "invalid",
		}

		_, _, err := parseQueryOrder(queries1)
		assert.Error(t, err)

		queries2 := map[string]string{}

		_, _, err = parseQueryOrder(queries2)
		assert.NoError(t, err)
	})
}

func TestParseQueryInt(t *testing.T) {
	t.Run("ParseQueryInt with valid query", func(t *testing.T) {
		queries := map[string]string{
			"page":  "2",
			"limit": "20",
		}

		page, err := parseQueryInt(queries, "page", 1)
		assert.NoError(t, err)
		assert.Equal(t, 2, page)

		limit, err := parseQueryInt(queries, "limit", 10)
		assert.NoError(t, err)
		assert.Equal(t, 20, limit)
	})

	t.Run("ParseQueryInt with missing query", func(t *testing.T) {
		queries := map[string]string{
			"page": "2",
		}

		limit, err := parseQueryInt(queries, "limit", 10)
		assert.NoError(t, err)
		assert.Equal(t, 10, limit)
	})

	t.Run("ParseQueryInt with invalid query", func(t *testing.T) {
		queries := map[string]string{
			"page": "invalid",
		}

		_, err := parseQueryInt(queries, "page", 1)
		assert.Error(t, err)
	})
}

func TestErrorResponse(t *testing.T) {
	c := fiber.New()

	ctx := c.AcquireCtx(&fasthttp.RequestCtx{})

	code := fiber.StatusBadRequest
	message := "Test error message"
	err := errors.New("Test error")

	response := errorResponse(ctx, code, message, err)

	assert.Empty(t, response, "Response should be empty")
}

func TestFilterMiddleware(t *testing.T) {
	c := fiber.New()
	c.Use(FilterMiddleware())
	c.Get("/vouchers/filter", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	t.Run("Valid Query1", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/vouchers/filter?name.equal=ahmad&page=1&limit=10&orderBy=field1&sortOrder=asc", nil)
		resp, _ := c.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Valid Query2", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/vouchers/filter?name.equal=ahmad&page=1&limit=10&orderBy=field1&sortOrder=desc", nil)
		resp, _ := c.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Invalid Query Key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/vouchers/filter?id", nil)
		resp, _ := c.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		// decode response body into domain.Response
		decodedResponse := domain.Response{}
		err := json.NewDecoder(resp.Body).Decode(&decodedResponse)
		assert.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, "Invalid query key", decodedResponse.Message)
	})

	t.Run("Invalid Page Query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/vouchers/filter?page=z", nil)
		resp, _ := c.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		// decode response body into domain.Response
		decodedResponse := domain.Response{}
		err := json.NewDecoder(resp.Body).Decode(&decodedResponse)
		assert.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, "Invalid page query", decodedResponse.Message)
	})

	t.Run("Invalid Limit Query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/vouchers/filter?page=1&limit=z", nil)
		resp, _ := c.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		// decode response body into domain.Response
		decodedResponse := domain.Response{}
		err := json.NewDecoder(resp.Body).Decode(&decodedResponse)
		assert.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, "Invalid limit query", decodedResponse.Message)
	})

	t.Run("Invalid sortOrder Query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/vouchers/filter?orderBy=name&sortOrder=invalid", nil)
		resp, _ := c.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		// decode response body into domain.Response
		decodedResponse := domain.Response{}
		err := json.NewDecoder(resp.Body).Decode(&decodedResponse)
		assert.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, "Invalid sortOrder query", decodedResponse.Message)
	})

	t.Run("Invalid Filter Query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/vouchers/filter?nominal.greaterThanOrEqual=invalid", nil)
		resp, _ := c.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		// decode response body into domain.Response
		decodedResponse := domain.Response{}
		err := json.NewDecoder(resp.Body).Decode(&decodedResponse)
		assert.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, "Invalid filter query", decodedResponse.Message)
	})
}
