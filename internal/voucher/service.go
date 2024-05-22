package voucher

import (
	"go-multiple-query/internal/domain"
)

type voucherService struct {
	voucherRepo domain.VoucherRepository
}

// Count implements domain.VoucherUsecase.
func (v *voucherService) Count(filter domain.VoucherFilter) (int64, error) {
	count, err := v.voucherRepo.Count(filter)
	if err != nil {
		return 0, err
	}

	return count, err
}

// FindWithFilter implements domain.VoucherUsecase.
func (v *voucherService) FindWithFilter(filter domain.VoucherFilter) ([]*domain.Voucher, int, error) {
	vouchers, nextCursor, err := v.voucherRepo.FindWithFilter(filter)
	if err != nil {
		return []*domain.Voucher{}, 0, err
	}

	return vouchers, nextCursor, err
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
