package middlewareHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/entities"
	middlewareUsecases "github.com/korvised/go-ecommerce/modules/middlewares/usecase"
)

type middlewareHandlerErrCode string

const (
	routerCheckErr middlewareHandlerErrCode = "middleware-001"
)

type IMiddlewareHandler interface {
	Cor() fiber.Handler
	RouterCheck() fiber.Handler
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

func (h *middlewareHandler) RouterCheck() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return entities.NewResponse(ctx).Error(
			fiber.StatusNotFound,
			string(routerCheckErr),
			"route not found",
		).Res()
	}
}
