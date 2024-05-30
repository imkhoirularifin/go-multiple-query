package voucher

import (
	"go-multiple-query/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type voucherService struct {
	voucherRepo domain.VoucherRepository
}

// FindByFilter implements domain.VoucherService.
func (v *voucherService) FindByFilter(query primitive.M, findOptions *options.FindOptions, page int, limit int) ([]*domain.Voucher, int64, int64, int, error) {
	vouchers, totalItem, nextCursor, maxPage, err := v.voucherRepo.FindByFilter(query, findOptions, page, limit)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return vouchers, totalItem, nextCursor, maxPage, nil
}

// Store implements domain.VoucherUsecase.
func (v *voucherService) Store(voucher *domain.Voucher) (*domain.Voucher, error) {
	voucher, err := v.voucherRepo.Store(voucher)
	if err != nil {
		return &domain.Voucher{}, err
	}

	return voucher, err
}

// NewVoucherService creates a new instance of VoucherService.
func NewVoucherService(voucherRepo domain.VoucherRepository) domain.VoucherService {
	return &voucherService{
		voucherRepo: voucherRepo,
	}
}
