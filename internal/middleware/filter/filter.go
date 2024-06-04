package filter

import (
	"errors"
	"fmt"
	"go-multiple-query/internal/domain"
	"go-multiple-query/internal/utilities"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// set query, findOptions, page, and limit to fiber context locals
func FilterMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		queries := c.Queries()
		var queryRequest domain.QueryRequest
		var findOptions *options.FindOptions

		queryFilters, err := parseQueriesToFilters(queries)
		if err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "Invalid query key", err)
		}

		queryRequest.Filters = queryFilters

		queryRequest.Page, err = parseQueryInt(queries, "page", 1)
		if err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "Invalid page query", err)
		}

		queryRequest.Limit, err = parseQueryInt(queries, "limit", 10)
		if err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "Invalid limit query", err)
		}

		queryRequest.OrderBy, queryRequest.SortOrder, err = parseQueryOrder(queries)
		if err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "Invalid sortOrder query", err)
		}

		query, err := buildQuery(queryRequest)
		if err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "Invalid filter query", err)
		}

		offset := (queryRequest.Page - 1) * queryRequest.Limit
		if queryRequest.OrderBy != "" {
			sortOrder := 1
			if queryRequest.SortOrder == "desc" {
				sortOrder = -1
			}

			findOptions = options.Find().
				SetLimit(int64(queryRequest.Limit)).
				SetSkip(int64(offset)).
				SetSort(bson.D{{Key: queryRequest.OrderBy, Value: sortOrder}})
		}

		c.Locals("query", query)
		c.Locals("findOptions", findOptions)
		c.Locals("page", queryRequest.Page)
		c.Locals("limit", queryRequest.Limit)
		return c.Next()
	}
}

func errorResponse(c *fiber.Ctx, code int, message string, err error) error {
	return c.Status(code).JSON(domain.Response{
		Code:    code,
		Status:  "error",
		Message: message,
		Errors:  []string{err.Error()},
	})
}

func parseQueryInt(queries map[string]string, key string, defaultValue int) (int, error) {
	str, ok := queries[key]
	if !ok {
		return defaultValue, nil
	}
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(value), nil
}

func parseQueryOrder(queries map[string]string) (string, string, error) {
	orderBy, ok := queries["orderBy"]
	if !ok {
		return "", "", nil
	}
	sortOrder, ok := queries["sortOrder"]
	if !ok {
		sortOrder = "asc"
	} else if sortOrder != "asc" && sortOrder != "desc" {
		return "", "", errors.New("invalid sort order")
	}
	return utilities.ToSnakeCase(orderBy), sortOrder, nil
}

func buildQuery(request domain.QueryRequest) (bson.M, error) {
	query := bson.M{}

	for _, q := range request.Filters {
		q.Key = utilities.ToSnakeCase(q.Key)
		for k, v := range domain.FilterCriteriaString {
			if q.Filter == k {
				if _, ok := query[q.Key]; !ok {
					query[q.Key] = bson.M{}
				}

				switch v {
				case "equal":
					if q.Key == "id" {
						id, err := primitive.ObjectIDFromHex(q.Value[0])
						if err != nil {
							return bson.M{}, errors.New("invalid ID")
						} else {
							delete(query, "id")
							query["_id"] = id
						}
					} else {
						value, err := strconv.ParseInt(q.Value[0], 10, 64)
						if err == nil {
							query[q.Key].(bson.M)["$eq"] = value
						} else {
							query[q.Key].(bson.M)["$eq"] = q.Value[0]
						}
					}
				case "notEqual":
					if q.Key == "id" {
						id, err := primitive.ObjectIDFromHex(q.Value[0])
						if err != nil {
							return bson.M{}, errors.New("invalid ID")
						} else {
							delete(query, "id")
							query["_id"] = bson.M{"$ne": id}
						}
					} else {
						value, err := strconv.ParseInt(q.Value[0], 10, 64)
						if err == nil {
							query[q.Key].(bson.M)["$ne"] = value
						} else {
							query[q.Key].(bson.M)["$ne"] = q.Value[0]
						}
					}
				case "greaterThanOrEqual":
					value, err := strconv.ParseInt(q.Value[0], 10, 64)
					if err != nil {
						return bson.M{}, errors.New("invalid value in greaterThanOrEqual query, value must be integer")
					} else {
						query[q.Key].(bson.M)["$gte"] = value
					}
				case "lessThanOrEqual":
					value, err := strconv.ParseInt(q.Value[0], 10, 64)
					if err != nil {
						return bson.M{}, errors.New("invalid value in lessThanOrEqual query, value must be integer")
					} else {
						query[q.Key].(bson.M)["$lte"] = value
					}
				case "greaterThan":
					value, err := strconv.ParseInt(q.Value[0], 10, 64)
					if err != nil {
						return bson.M{}, errors.New("invalid value in greaterThan query, value must be integer")
					} else {
						query[q.Key].(bson.M)["$gt"] = value
					}
				case "lessThan":
					value, err := strconv.ParseInt(q.Value[0], 10, 64)
					if err != nil {
						return bson.M{}, errors.New("invalid value in lessThan query, value must be integer")
					} else {
						query[q.Key].(bson.M)["$lt"] = value
					}
				case "in":
					if q.Key == "id" {
						values := make([]primitive.ObjectID, len(q.Value))
						for i, str := range q.Value {
							id, err := primitive.ObjectIDFromHex(str)
							if err != nil {
								return bson.M{}, errors.New("invalid ID")
							}
							values[i] = id
						}
						delete(query, "id")
						query["_id"] = bson.M{"$in": values}
					} else {
						valuesInt := make([]int64, len(q.Value))
						valuesStr := make([]string, len(q.Value))
						hasIntValues := false
						for i, str := range q.Value {
							val, err := strconv.ParseInt(str, 10, 64)
							if err != nil {
								valuesStr[i] = str
							} else {
								valuesInt[i] = val
								hasIntValues = true
							}
						}

						if hasIntValues {
							query[q.Key].(bson.M)["$in"] = valuesInt
						} else {
							query[q.Key].(bson.M)["$in"] = valuesStr
						}
					}
				}
			}
		}
	}

	return query, nil
}

func parseQueriesToFilters(queries map[string]string) ([]domain.QueryFilter, error) {
	var filters []domain.QueryFilter

	for k, v := range queries {
		if k == "page" || k == "limit" || k == "orderBy" || k == "sortOrder" {
			continue
		}

		parts := strings.Split(k, ".")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid query key: %s", k)
		}

		field := parts[0]
		criteriaStr := parts[1]

		var criteria domain.FilterCriteria
		for k, v := range domain.FilterCriteriaString {
			if v == criteriaStr {
				criteria = k
				break
			}
		}

		if criteria == domain.NoMatch {
			return nil, fmt.Errorf("invalid filter criteria: %s", criteriaStr)
		}

		var values []string
		if strings.Split(v, ",") != nil && len(strings.Split(v, ",")) > 1 {
			values = strings.Split(v, ",")
		} else {
			values = append(values, v)
		}

		filters = append(filters, domain.QueryFilter{
			Key:    field,
			Filter: criteria,
			Value:  values,
		})
	}

	return filters, nil
}
