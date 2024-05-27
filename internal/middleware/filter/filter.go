package filter

import (
	"context"
	"fmt"
	"go-multiple-query/internal/domain"
	"go-multiple-query/internal/utilities"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FilterMiddleware[entity any](db *mongo.Database, collectionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		queries := c.Queries()
		var entities []entity
		var cursor *mongo.Cursor
		page := queries["page"]
		limit := queries["limit"]
		orderBy := queries["orderBy"]
		sortOrder := queries["sortOrder"]
		var queryRequest domain.QueryRequest
		var findOptions *options.FindOptions
		var totalItem int64
		var nextCursor int64
		var maxPage int

		queryFilters, err := parseQueriesToFilters(queries)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.Response{
				Code:    fiber.StatusBadRequest,
				Errors:  []string{err.Error()},
				Message: "invalid query key",
			})
		}

		queryRequest.Filters = queryFilters

		if page == "" {
			queryRequest.Page = 1
		} else {
			p, err := strconv.ParseInt(page, 10, 64)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(domain.Response{
					Code:    fiber.StatusBadRequest,
					Status:  "error",
					Message: "Invalid page query",
				})
			}
			queryRequest.Page = int(p)
		}

		nextCursor = int64(queryRequest.Page + 1)

		if limit == "" {
			queryRequest.Limit = 10
		} else {
			l, err := strconv.ParseInt(limit, 10, 64)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(domain.Response{
					Code:    fiber.StatusBadRequest,
					Status:  "error",
					Message: "Invalid limit query",
				})
			}
			queryRequest.Limit = int(l)
		}

		if orderBy != "" {
			convertedOrderBy := utilities.ToSnakeCase(orderBy)
			queryRequest.OrderBy = convertedOrderBy

			if sortOrder != "" {
				// check if sortOrder is valid
				if sortOrder != "asc" && sortOrder != "desc" {
					return c.Status(fiber.StatusBadRequest).JSON(domain.Response{
						Code:    fiber.StatusBadRequest,
						Status:  "error",
						Message: "Invalid sort order",
					})
				}

				queryRequest.SortOrder = sortOrder
			} else {
				queryRequest.SortOrder = "asc"
			}
		}

		coll := db.Collection(collectionName)
		query := bson.M{}
		offset := (queryRequest.Page - 1) * queryRequest.Limit

		if queryRequest.OrderBy != "" {
			var sortOrder int
			if queryRequest.SortOrder == "" {
				sortOrder = 1
			} else if queryRequest.SortOrder == "asc" {
				sortOrder = 1
			} else {
				sortOrder = -1
			}

			findOptions = options.Find().
				SetLimit(int64(queryRequest.Limit)).
				SetSkip(int64(offset)).
				SetSort(bson.D{{Key: queryRequest.OrderBy, Value: sortOrder}})
		}

		for _, q := range queryRequest.Filters {
			// q.Key format is camelCase, convert to snake_case
			q.Key = utilities.ToSnakeCase(q.Key)

			// loop over filter criteria string map
			for k, v := range domain.FilterCriteriaString {
				if q.Filter == k {
					switch v {
					case "equal":
						if q.Key == "id" {
							id, err := primitive.ObjectIDFromHex(q.Value[0])
							if err != nil {
								c.Status(fiber.StatusBadRequest).JSON(domain.Response{
									Code:    fiber.StatusBadRequest,
									Status:  "error",
									Message: "Invalid ID",
								})
							} else {
								query["_id"] = id
							}
						} else {
							value, err := strconv.ParseInt(q.Value[0], 10, 64)
							if err == nil {
								query[q.Key] = value
							} else {
								query[q.Key] = q.Value[0]
							}
						}
					case "notEqual":
						if q.Key == "id" {
							id, err := primitive.ObjectIDFromHex(q.Value[0])
							if err != nil {
								c.Status(fiber.StatusBadRequest).JSON(domain.Response{
									Code:    fiber.StatusBadRequest,
									Status:  "error",
									Message: "Invalid ID",
								})
							} else {
								query["_id"] = bson.M{"$ne": id}
							}
						} else {
							value, err := strconv.ParseInt(q.Value[0], 10, 64)
							if err == nil {
								query[q.Key] = bson.M{"$ne": value}
							} else {
								query[q.Key] = bson.M{"$ne": q.Value[0]}
							}
						}
					case "greaterThanOrEqual":
						value, err := strconv.ParseInt(q.Value[0], 10, 64)
						if err == nil {
							query[q.Key] = bson.M{"$gte": value}
						} else {
							query[q.Key] = bson.M{"$gte": q.Value[0]}
						}
					case "lessThanOrEqual":
						value, err := strconv.ParseInt(q.Value[0], 10, 64)
						if err == nil {
							query[q.Key] = bson.M{"$lte": value}
						} else {
							query[q.Key] = bson.M{"$lte": q.Value[0]}
						}
					case "greaterThan":
						value, err := strconv.ParseInt(q.Value[0], 10, 64)
						if err == nil {
							query[q.Key] = bson.M{"$gt": value}
						} else {
							query[q.Key] = bson.M{"$gt": q.Value[0]}
						}
					case "lessThan":
						value, err := strconv.ParseInt(q.Value[0], 10, 64)
						if err == nil {
							query[q.Key] = bson.M{"$lt": value}
						} else {
							query[q.Key] = bson.M{"$lt": q.Value[0]}
						}
					case "in":
						if q.Key == "id" {
							values := make([]primitive.ObjectID, len(q.Value))
							for i, str := range q.Value {
								id, err := primitive.ObjectIDFromHex(str)
								if err != nil {
									// handle error, for example log it and continue to the next iteration
									fmt.Println("error parsing id", err.Error())
									continue
								}
								values[i] = id
							}
							query["_id"] = bson.M{"$in": values}
						} else {
							// assuming q.Value is a slice of strings
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
								query[q.Key] = bson.M{"$in": valuesInt}
							} else {
								query[q.Key] = bson.M{"$in": valuesStr}
							}
						}
					}
				}
			}
		}

		// check if findOptions is nil
		if findOptions != nil {
			cursor, err = coll.Find(context.TODO(), query, findOptions)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(domain.Response{
					Code:    fiber.StatusInternalServerError,
					Status:  "error",
					Message: err.Error(),
				})
			}
		} else {
			cursor, err = coll.Find(context.TODO(), query)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(domain.Response{
					Code:    fiber.StatusInternalServerError,
					Status:  "error",
					Message: err.Error(),
				})
			}
		}

		defer cursor.Close(context.TODO())

		for cursor.Next(context.TODO()) {
			// data is single data of the collection
			var data entity

			err := cursor.Decode(&data)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(domain.Response{
					Code:    fiber.StatusInternalServerError,
					Status:  "error",
					Message: err.Error(),
				})
			}

			entities = append(entities, data)
		}

		// Check for any errors during cursor iteration
		if err := cursor.Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.Response{
				Code:    fiber.StatusInternalServerError,
				Status:  "error",
				Message: err.Error(),
			})
		}

		if len(entities) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(domain.Response{
				Code:    fiber.StatusNotFound,
				Status:  "error",
				Message: "No data found",
			})
		}

		// get total item
		totalItem, err = coll.CountDocuments(context.TODO(), query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.Response{
				Code:    fiber.StatusInternalServerError,
				Status:  "error",
				Message: err.Error(),
			})
		}

		maxPage = int(math.Ceil(float64(totalItem) / float64(queryRequest.Limit)))

		c.Locals("entities", entities)
		c.Locals("totalItem", totalItem)
		c.Locals("nextCursor", nextCursor)
		c.Locals("maxPage", maxPage)
		return c.Next()
	}
}

func parseQueriesToFilters(queries map[string]string) ([]domain.QueryFilter, error) {
	var filters []domain.QueryFilter

	for k, v := range queries {
		// skip page, limit, orderBy, and sortOrder
		if k == "page" || k == "limit" || k == "orderBy" || k == "sortOrder" {
			continue
		}

		parts := strings.Split(k, ".")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid query key: %s", k)
		}

		// field is key, criteriaStr is filter criteria and v is value
		field := parts[0]
		criteriaStr := parts[1]

		var criteria domain.FilterCriteria
		for k, v := range domain.FilterCriteriaString {
			if v == criteriaStr {
				criteria = k
				break
			}
		}

		// v is array of string, loop over it and append to values
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
