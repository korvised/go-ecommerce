package productsHandlers

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/appinfo"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/files"
	"github.com/korvised/go-ecommerce/modules/files/filesUsecases"
	"github.com/korvised/go-ecommerce/modules/products"
	"github.com/korvised/go-ecommerce/modules/products/productsUsecases"
	"log"
	"strings"
)

type productsHandlersErrCode string

const (
	findOneProductErr  productsHandlersErrCode = "products-001"
	findManyProductErr productsHandlersErrCode = "products-002"
	addProductErr      productsHandlersErrCode = "products-003"
	updateProductErr   productsHandlersErrCode = "products-004"
	deleteProductErr   productsHandlersErrCode = "products-005"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindManyProducts(c *fiber.Ctx) error
	AddProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type productsHandler struct {
	cfg             config.IConfig
	productsUsecase productsUsecases.IProductsUsecase
	filesUsecase    filesUsecases.IFilesUsecase
}

func ProductsHandler(
	cfg config.IConfig,
	productsUsecase productsUsecases.IProductsUsecase,
	filesUsecase filesUsecases.IFilesUsecase,
) IProductsHandler {
	return &productsHandler{
		cfg:             cfg,
		productsUsecase: productsUsecase,
		filesUsecase:    filesUsecase,
	}
}

func (h *productsHandler) FindOneProduct(c *fiber.Ctx) error {
	productID := strings.Trim(c.Params("product_id"), "")

	product, err := h.productsUsecase.FindOneProduct(productID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(findOneProductErr),
				"product not found",
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(findOneProductErr),
				err.Error(),
			).Res()
		}

	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productsHandler) FindManyProducts(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(findManyProductErr), err.Error()).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Size < 5 {
		req.Size = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}

	if req.Sort == "" {
		req.Sort = "ASC"
	}

	data := h.productsUsecase.FindManyProducts(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, data).Res()
}

func (h *productsHandler) AddProduct(c *fiber.Ctx) error {
	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(addProductErr), err.Error()).Res()
	}

	product, err := h.productsUsecase.AddProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(addProductErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, product).Res()
}

func (h *productsHandler) UpdateProduct(c *fiber.Ctx) error {
	productID := strings.Trim(c.Params("product_id"), " ")

	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(updateProductErr), err.Error()).Res()
	}
	req.ID = productID

	product, err := h.productsUsecase.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(updateProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, product).Res()
}

func (h *productsHandler) DeleteProduct(c *fiber.Ctx) error {
	productID := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecase.FindOneProduct(productID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(findOneProductErr),
				"product not found",
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(findOneProductErr),
				err.Error(),
			).Res()
		}

	}

	// Delete images
	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, img := range product.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("products/%s", img.FileName),
		})
	}

	if err = h.filesUsecase.DeleteFileOnStorage(deleteFileReq); err != nil {
		log.Printf("delete image failed: %v", err)
	}

	// Delete product
	if err := h.productsUsecase.DeleteProduct(productID); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
