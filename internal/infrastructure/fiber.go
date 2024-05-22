package infrastructure

import (
	"fmt"
	"go-multiple-query/internal/docs"
	"go-multiple-query/internal/voucher"
	"go-multiple-query/pkg/xlogger"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Run() {
	logger := xlogger.Logger

	app := fiber.New(fiber.Config{
		ProxyHeader:           cfg.ProxyHeader,
		DisableStartupMessage: true,
		ErrorHandler:          defaultErrorHandler,
	})

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
		Fields: cfg.LogFields,
	}))
	app.Use(recover2.New())
	app.Use(etag.New())
	app.Use(requestid.New())

	// Grouping Routes
	api := app.Group("/api")
	docs.NewHttpHandler(api.Group("/docs"))
	voucher.NewHTTPHandler(api.Group("/vouchers"), voucherService, logger)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info().Msgf("Server is running on address: %s", addr)
	if err := app.Listen(addr); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
