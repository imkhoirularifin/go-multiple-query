package voucher

import (
	"go-multiple-query/internal/domain"
)

type voucherService struct {
	voucherRepo domain.VoucherRepository
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
