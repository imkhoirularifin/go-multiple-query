package voucher

import (
	"context"
	"go-multiple-query/internal/domain"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
