package voucher

import (
	"go-multiple-query/internal/domain"
	"go-multiple-query/internal/middleware/filter"
	"go-multiple-query/internal/middleware/validation"
	"go-multiple-query/internal/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type httpHandler struct {
	voucherService domain.VoucherService
	db             *mongo.Database
}

// NewHTTPHandler creates a new instance of HTTPHandler.
func NewHTTPHandler(r fiber.Router, voucherService domain.VoucherService, db *mongo.Database, logger *zerolog.Logger) {
	handler := &httpHandler{
		voucherService: voucherService,
		db:             db,
	}

	r.Post("/", validation.New[domain.StoreVoucherRequest](), handler.Store)
	r.Get("/filter", filter.FilterMiddleware(), handler.FindWithFilter)
}

// Store handles the store voucher request.
func (h *httpHandler) Store(c *fiber.Ctx) error {
	storeVoucherReq := utilities.ExtractStructFromValidator[domain.StoreVoucherRequest](c)

	voucher := domain.Voucher{
		BrandCode:        storeVoucherReq.BrandCode,
		Sku:              storeVoucherReq.Sku,
		SkuName:          storeVoucherReq.SkuName,
		Nominal:          storeVoucherReq.Nominal,
		DistributorPrice: storeVoucherReq.DistributorPrice,
		ProductStatus:    storeVoucherReq.ProductStatus,
		OrderDestination: storeVoucherReq.OrderDestination,
		Stock:            storeVoucherReq.Stock,
		Vendor:           storeVoucherReq.Vendor,
	}

	result, err := h.voucherService.Store(&voucher)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Response{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Voucher has been stored successfully",
		Data:    result,
	})
}

func (h *httpHandler) FindWithFilter(c *fiber.Ctx) error {
	query, findOptions, page, limit := utilities.ExtractQueryAndFindOptions(c)

	vouchers, totalItem, nextCursor, maxPage, err := h.voucherService.FindByFilter(query, findOptions, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Response{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	if int(nextCursor) > 0 && int(nextCursor) <= maxPage {
		c.Set("X-Cursor", strconv.Itoa(int(nextCursor)))
	}
	c.Set("X-Total-Count", strconv.Itoa(int(totalItem)))
	c.Set("X-Max-Page", strconv.Itoa(maxPage))

	return c.Status(fiber.StatusOK).JSON(domain.Response{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Vouchers have been fetched successfully",
		Data:    vouchers,
	})
}
