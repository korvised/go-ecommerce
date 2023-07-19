package middlewareHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/korvised/go-ecommerce/config"
	middlewareUsecases "github.com/korvised/go-ecommerce/modules/middlewares/usecase"
)

type IMiddlewareHandler interface {
	Cor() fiber.Handler
}

type middlewareHandler struct {
	cfg               config.IConfig
	middlewareUsecase middlewareUsecases.IMiddlewareUsecase
}

func MiddlewareHandler(cfg config.IConfig, middlewareUsecase middlewareUsecases.IMiddlewareUsecase) IMiddlewareHandler {
	return &middlewareHandler{
		cfg:               cfg,
		middlewareUsecase: middlewareUsecase,
	}
}

func (h *middlewareHandler) Cor() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}
