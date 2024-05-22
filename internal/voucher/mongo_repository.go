package voucher

import (
	"context"
	"go-multiple-query/internal/domain"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbRepository struct {
	db *mongo.Database
}

// FindByID implements domain.VoucherRepository.
func (m *mongodbRepository) FindByID(id primitive.ObjectID) (*domain.Voucher, error) {
	coll := m.db.Collection("vouchers")

	var voucher domain.Voucher
	err := coll.FindOne(context.TODO(), primitive.M{"_id": id}).Decode(&voucher)
	if err != nil {
		return nil, err
	}

	return &voucher, nil
}

// Count implements domain.VoucherRepository.
func (m *mongodbRepository) Count(filter domain.VoucherFilter) (int64, error) {
	var count int64

	coll := m.db.Collection("vouchers")

	// loop over filter and build query
	query := bson.M{}
	v := reflect.ValueOf(filter)
	typeOfFilter := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfFilter.Field(i).Name

		// Skip pagination and sorting
		if fieldName == "Page" || fieldName == "Size" || fieldName == "OrderBy" || fieldName == "SortOrder" {
			continue
		}

		fieldValue := v.Field(i).Interface()
		if str, ok := fieldValue.(string); ok && str != "" {
			query[typeOfFilter.Field(i).Tag.Get("query")] = fieldValue
		}
	}

	count, err := coll.CountDocuments(context.TODO(), query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindWithFilter implements domain.VoucherRepository.
func (m *mongodbRepository) FindWithFilter(filter domain.VoucherFilter) ([]*domain.Voucher, int, error) {
	coll := m.db.Collection("vouchers")
	var vouchers []*domain.Voucher

	page, _ := strconv.Atoi(filter.Page)
	size, _ := strconv.Atoi(filter.Size)
	offset := (page - 1) * size

	// loop over filter and build query
	query := bson.M{}
	v := reflect.ValueOf(filter)
	typeOfFilter := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfFilter.Field(i).Name

		// Skip pagination and sorting
		if fieldName == "Page" || fieldName == "Size" || fieldName == "OrderBy" || fieldName == "SortOrder" {
			continue
		}

		fieldValue := v.Field(i).Interface()
		if str, ok := fieldValue.(string); ok && str != "" {
			query[typeOfFilter.Field(i).Tag.Get("query")] = fieldValue
		}
	}

	var sortOrder int
	if filter.SortOrder == "asc" {
		sortOrder = 1
	} else {
		sortOrder = -1
	}

	findOptions := options.Find().
		SetLimit(int64(size)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: filter.OrderBy, Value: sortOrder}})

	cursor, err := coll.Find(context.TODO(), query, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var voucher domain.Voucher
		err := cursor.Decode(&voucher)
		if err != nil {
			return nil, 0, err
		}

		vouchers = append(vouchers, &voucher)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	if len(vouchers) == 0 {
		return nil, 0, mongo.ErrNoDocuments
	}

	var nextCursor int
	if len(vouchers) > 0 {
		nextCursor = page + 1
	}

	return vouchers, nextCursor, nil
}

// Store implements domain.VoucherRepository.
func (m *mongodbRepository) Store(voucher *domain.Voucher) (*domain.Voucher, error) {
	coll := m.db.Collection("vouchers")

	result, err := coll.InsertOne(context.TODO(), voucher)
	if err != nil {
		return &domain.Voucher{}, err
	}

	// get by id
	voucher, err = m.FindByID(result.InsertedID.(primitive.ObjectID))
	if err != nil {
		return &domain.Voucher{}, err
	}

	return voucher, nil
}

func NewMongoRepository(db *mongo.Database) domain.VoucherRepository {
	return &mongodbRepository{db}
}
