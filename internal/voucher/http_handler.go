package voucher

import (
	"go-multiple-query/internal/domain"
	"go-multiple-query/internal/middleware/validation"
	"go-multiple-query/internal/utilities"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type httpHandler struct {
	voucherService domain.VoucherService
}

// NewHTTPHandler creates a new instance of HTTPHandler.
func NewHTTPHandler(r fiber.Router, voucherService domain.VoucherService, logger *zerolog.Logger) {
	handler := &httpHandler{
		voucherService: voucherService,
	}

	r.Post("/", validation.New[domain.StoreVoucherRequest](), handler.Store)
	r.Get("/filter", handler.FindWithFilter)
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

// FindWithFilter handles the find with filter request.
func (h *httpHandler) FindWithFilter(c *fiber.Ctx) error {
	filter := new(domain.VoucherFilter)

	if err := c.QueryParser(filter); err != nil {
		return err
	}

	defaults := map[string]*string{
		"1":        &filter.Page,
		"10":       &filter.Size,
		"sku_name": &filter.OrderBy,
		"asc":      &filter.SortOrder,
	}

	for defaultValue, field := range defaults {
		if *field == "" {
			*field = defaultValue
		}
	}

	vouchers, nextPage, err := h.voucherService.FindWithFilter(*filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(domain.Response{
				Code:    fiber.StatusNotFound,
				Status:  "error",
				Message: "Vouchers not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Response{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	totalItem, err := h.voucherService.Count(*filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Response{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	size, _ := strconv.Atoi(filter.Size)
	maxPage := int(math.Ceil(float64(totalItem) / float64(size)))

	if nextPage > 0 && nextPage <= maxPage {
		c.Set("X-Cursor", strconv.Itoa(nextPage))
	}
	c.Set("X-Total-Count", strconv.Itoa(int(totalItem)))
	c.Set("X-Max-Page", strconv.Itoa(maxPage))

	return c.JSON(domain.Response{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Vouchers have been fetched successfully",
		Data:    vouchers,
	})
}
