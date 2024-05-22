package infrastructure

import (
	"go-multiple-query/internal/config"
	"go-multiple-query/internal/domain"
	"go-multiple-query/internal/voucher"
	"go-multiple-query/pkg/xlogger"

	"github.com/caarlos0/env/v10"

	_ "github.com/joho/godotenv/autoload"
)

var (
	cfg config.Config

	voucherRepo domain.VoucherRepository

	voucherService domain.VoucherService
)

func init() {
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	xlogger.Setup(cfg)

	db := mongodbSetup()

	voucherRepo = voucher.NewMongoRepository(db)

	voucherService = voucher.NewVoucherService(voucherRepo)
}
