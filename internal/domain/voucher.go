package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Voucher struct {
	Id               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BrandCode        string             `json:"brand_code" bson:"brand_code" query:"brand_code"`
	Sku              string             `json:"sku" bson:"sku" query:"sku"`
	SkuName          string             `json:"sku_name" bson:"sku_name" query:"sku_name"`
	Nominal          int                `json:"nominal" bson:"nominal" query:"nominal"`
	DistributorPrice int                `json:"distributor_price" bson:"distributor_price" query:"distributor_price"`
	ProductStatus    string             `json:"product_status" bson:"product_status" query:"product_status"`
	OrderDestination string             `json:"order_destination" bson:"order_destination" query:"order_destination"`
	Stock            int                `json:"stock" bson:"stock" query:"stock"`
	Vendor           string             `json:"vendor" bson:"vendor" query:"vendor"`
}

type VoucherRepository interface {
	FindByID(id primitive.ObjectID) (*Voucher, error)
	FindByFilter(query primitive.M, findOptions *options.FindOptions, page int, limit int) ([]*Voucher, int64, int64, int, error)
	Store(voucher *Voucher) (*Voucher, error)
}

type VoucherService interface {
	FindByFilter(query primitive.M, findOptions *options.FindOptions, page int, limit int) ([]*Voucher, int64, int64, int, error)
	Store(voucher *Voucher) (*Voucher, error)
}

type StoreVoucherRequest struct {
	BrandCode        string `json:"brand_code" validate:"required"`
	Sku              string `json:"sku" validate:"required"`
	SkuName          string `json:"sku_name" validate:"required"`
	Nominal          int    `json:"nominal" validate:"required"`
	DistributorPrice int    `json:"distributor_price" validate:"required"`
	ProductStatus    string `json:"product_status" validate:"required"`
	OrderDestination string `json:"order_destination" validate:"required"`
	Stock            int    `json:"stock" validate:"required"`
	Vendor           string `json:"vendor" validate:"required"`
}

type VoucherFilter struct {
	BrandCode        string `query:"brand_code"`
	Sku              string `query:"sku"`
	SkuName          string `query:"sku_name"`
	Nominal          string `query:"nominal"`
	DistributorPrice string `query:"distributor_price"`
	ProductStatus    string `query:"product_status"`
	OrderDestination string `query:"order_destination"`
	Stock            string `query:"stock"`
	Vendor           string `query:"vendor"`
	OrderBy          string `query:"order_by"`
	SortOrder        string `query:"sort_order"`
	Page             string `query:"page"`
	Size             string `query:"size"`
}
