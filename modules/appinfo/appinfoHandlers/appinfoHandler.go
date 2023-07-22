package appinfoHandlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/appinfo"
	"github.com/korvised/go-ecommerce/modules/appinfo/appinfoUsecases"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/pkg/auth"
	"strconv"
)

type appinfoHandlersErrCode string

const (
	generateApiKeyErr   appinfoHandlersErrCode = "appinfo-001"
	findCategoriesErr   appinfoHandlersErrCode = "appinfo-002"
	addCategoriesErr    appinfoHandlersErrCode = "appinfo-003"
	deleteCategoriesErr appinfoHandlersErrCode = "appinfo-004"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategories(c *fiber.Ctx) error
	AddCategories(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
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

func (h *appinfoHandler) FindCategories(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(findCategoriesErr), err.Error()).Res()
	}

	categories, err := h.appinfoUsecase.FindCategories(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(findCategoriesErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, categories).Res()
}

func (h *appinfoHandler) AddCategories(c *fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		fmt.Println("err BodyParser", err)
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(addCategoriesErr), err.Error()).Res()
	}

	fmt.Println("req body", req)

	if err := h.appinfoUsecase.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(addCategoriesErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, req).Res()
}

func (h *appinfoHandler) DeleteCategory(c *fiber.Ctx) error {
	categoryId, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(deleteCategoriesErr),
			"category id must be a number",
		).Res()
	}

	if err := h.appinfoUsecase.DeleteCategory(categoryId); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(deleteCategoriesErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
