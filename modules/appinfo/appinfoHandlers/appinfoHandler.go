package appinfoHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/appinfo/appinfoUsecases"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/pkg/auth"
)

type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfo-001"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg            config.IConfig
	appinfoUsecase appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:            cfg,
		appinfoUsecase: appinfoUsecase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := auth.NewAuth(auth.ApiKey, h.cfg.Jwt(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(generateApiKeyErr), err.Error()).Res()
	}

	data := &struct {
		Key string `json:"key"`
	}{
		Key: apiKey.SignToken(),
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, data).Res()
}
