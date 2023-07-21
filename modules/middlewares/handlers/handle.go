package middlewareHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/entities"
	middlewaresUsecases "github.com/korvised/go-ecommerce/modules/middlewares/usecase"
	"github.com/korvised/go-ecommerce/pkg/auth"
	"strings"
)

type middlewaresHandlerErrCode string

const (
	routerCheckErr middlewaresHandlerErrCode = "middleware-001"
	jwtAuthErr     middlewaresHandlerErrCode = "middleware-003"
)

const unauthorizedMsg = "unauthorized, no permission to access this route"

type IMiddlewaresHandler interface {
	Cor() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
}

type middlewaresHandler struct {
	cfg                config.IConfig
	middlewaresUsecase middlewaresUsecases.IMiddlewareUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUsecase middlewaresUsecases.IMiddlewareUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:                cfg,
		middlewaresUsecase: middlewaresUsecase,
	}
}

func (h *middlewaresHandler) Cor() fiber.Handler {
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

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return entities.NewResponse(ctx).Error(
			fiber.StatusNotFound,
			string(routerCheckErr),
			"route not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006 15:04:05",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h middlewaresHandler) JwtAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorization := ctx.Get("Authorization")
		if authorization == "" {
			return entities.NewResponse(ctx).Error(fiber.StatusUnauthorized, string(jwtAuthErr), unauthorizedMsg).Res()

		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(ctx).Error(fiber.StatusUnauthorized, string(jwtAuthErr), err.Error()).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(ctx).Error(fiber.StatusUnauthorized, string(jwtAuthErr), unauthorizedMsg).Res()
		}

		// Set userId
		ctx.Locals("userId", claims.Id)
		ctx.Locals("userRoleId", claims.RoleId)

		return ctx.Next()
	}
}
