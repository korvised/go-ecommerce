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
	ApiKey                                   = "X-Api-Key"
	routerCheckErr middlewaresHandlerErrCode = "middleware-001"
	jwtAuthErr     middlewaresHandlerErrCode = "middleware-002"
	paramsCheckErr middlewaresHandlerErrCode = "middleware-003"
	authorizeErr   middlewaresHandlerErrCode = "middleware-004"
	apiKeyErr      middlewaresHandlerErrCode = "middleware-005"
)

const unauthorizedMsg = "unauthorized, no permission to access this route"
const requiredApiKeyMsg = "unauthorized, api key is required"
const invalidApiKeyMsg = "unauthorized, invalid api key"

type IMiddlewaresHandler interface {
	Cor() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRoleIds ...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
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
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
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
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" {
			return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(jwtAuthErr), unauthorizedMsg).Res()

		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(jwtAuthErr), err.Error()).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(jwtAuthErr), unauthorizedMsg).Res()
		}

		// Set userId
		c.Locals(UserId, claims.Id)
		c.Locals(UserRoleId, claims.RoleId)

		return c.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals(UserId)

		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(paramsCheckErr), unauthorizedMsg).Res()
		}

		return c.Next()
	}
}

func (h *middlewaresHandler) Authorize(expectRoleIds ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals(UserRoleId).(int)
		if !ok {
			return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(authorizeErr), unauthorizedMsg).Res()
		}

		roles, err := h.middlewaresUsecase.FindRoles()
		if err != nil {
			return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(authorizeErr), err.Error()).Res()
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
				return c.Next()
			}
		}

		return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(authorizeErr), unauthorizedMsg).Res()
	}
}

func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get(ApiKey)
		if apiKey == "" {
			return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(apiKeyErr), requiredApiKeyMsg).Res()
		}

		if _, err := auth.ParseApiKey(h.cfg.Jwt(), apiKey); err != nil {
			return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(apiKeyErr), invalidApiKeyMsg).Res()
		}

		return c.Next()
	}
}
