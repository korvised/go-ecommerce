package middlewaresHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/middlewares/middlewaresUsecases"
	"github.com/korvised/go-ecommerce/pkg/auth"
	"github.com/korvised/go-ecommerce/pkg/utils"
	"strings"
)

type middlewaresHandlerErrCode string

const (
	UserId                                   = "UserId"
	UserRoleId                               = "UserRoleId"
	routerCheckErr middlewaresHandlerErrCode = "middleware-001"
	jwtAuthErr     middlewaresHandlerErrCode = "middleware-002"
	paramsCheckErr middlewaresHandlerErrCode = "middleware-003"
	authorizeErr   middlewaresHandlerErrCode = "middleware-004"
)

const unauthorizedMsg = "unauthorized, no permission to access this route"

type IMiddlewaresHandler interface {
	Cor() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRoleIds ...int) fiber.Handler
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
		ctx.Locals(UserId, claims.Id)
		ctx.Locals(UserRoleId, claims.RoleId)

		return ctx.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId := ctx.Locals(UserId)

		if ctx.Params("user_id") != userId {
			return entities.NewResponse(ctx).Error(fiber.StatusUnauthorized, string(paramsCheckErr), unauthorizedMsg).Res()
		}

		return ctx.Next()
	}
}

func (h *middlewaresHandler) Authorize(expectRoleIds ...int) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userRoleId, ok := ctx.Locals(UserRoleId).(int)
		if !ok {
			return entities.NewResponse(ctx).Error(fiber.StatusUnauthorized, string(authorizeErr), unauthorizedMsg).Res()
		}

		roles, err := h.middlewaresUsecase.FindRoles()
		if err != nil {
			return entities.NewResponse(ctx).Error(fiber.StatusInternalServerError, string(authorizeErr), err.Error()).Res()
		}

		sum := 0
		for _, v := range expectRoleIds {
			sum += v
		}

		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))

		// user ->     0 1 0
		// expected -> 1 1 0

		for i := range userValueBinary {
			if userValueBinary[i]&expectedValueBinary[i] == 1 {
				return ctx.Next()
			}
		}

		return entities.NewResponse(ctx).Error(fiber.StatusUnauthorized, string(authorizeErr), unauthorizedMsg).Res()
	}
}
