package voucher

import (
	"context"
	"errors"
	"go-multiple-query/internal/domain"
	"math"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbRepository struct {
	db *mongo.Database
}

// FindByFilter implements domain.VoucherRepository.
func (m *mongodbRepository) FindByFilter(query primitive.M, findOptions *options.FindOptions, page int, limit int) ([]*domain.Voucher, int64, int64, int, error) {
	coll := m.db.Collection("vouchers")
	var cursor *mongo.Cursor
	var err error
	vouchers := []*domain.Voucher{}
	nextCursor := int64(page + 1)
	maxPage := 0

	totalItem, err := coll.CountDocuments(context.TODO(), query)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if findOptions != nil {
		cursor, err = coll.Find(context.TODO(), query, findOptions)
	} else {
		cursor, err = coll.Find(context.TODO(), query)
	}
	if err != nil {
		return nil, 0, 0, 0, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var voucher domain.Voucher
		if err := cursor.Decode(&voucher); err != nil {
			continue
		}
		vouchers = append(vouchers, &voucher)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, 0, 0, err
	}

	if len(vouchers) == 0 {
		return nil, 0, 0, 0, errors.New("no vouchers found")
	}

	maxPage = int(math.Ceil(float64(totalItem) / float64(limit)))

	return vouchers, totalItem, nextCursor, maxPage, nil
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
