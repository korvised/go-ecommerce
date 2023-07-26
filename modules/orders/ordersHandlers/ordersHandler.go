package ordersHandlers

import (
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/orders/ordersUsecases"
)

type IOrdersHandler interface {
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
