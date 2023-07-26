package ordersHandlers

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/orders/ordersUsecases"
	"strings"
)

type ordersHandlersErrCode string

const (
	fineOneOrderErr ordersHandlersErrCode = "orders-001"
)

type IOrdersHandler interface {
	FindOneOrder(c *fiber.Ctx) error
}

type ordersHandler struct {
	cfg           config.IConfig
	ordersUsecase ordersUsecases.IOrdersUsecase
}

func OrdersHandler(cfg config.IConfig, ordersUsecase ordersUsecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandler{
		cfg:           cfg,
		ordersUsecase: ordersUsecase,
	}
}

func (h *ordersHandler) FindOneOrder(c *fiber.Ctx) error {
	orderID := strings.Trim(c.Params("order_id"), " ")

	order, err := h.ordersUsecase.FindOneOrder(orderID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(fineOneOrderErr),
				"order not found",
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(fineOneOrderErr),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
