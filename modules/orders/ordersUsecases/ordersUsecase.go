package ordersUsecases

import (
	"github.com/korvised/go-ecommerce/modules/orders/ordersRepositories"
	"github.com/korvised/go-ecommerce/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
}

type ordersUsecase struct {
	ordersRepository   ordersRepositories.IOrdersRepository
	productsRepository productsRepositories.IProductsRepository
}

func OrdersUsecase(
	ordersRepository ordersRepositories.IOrdersRepository,
	productsRepository productsRepositories.IProductsRepository,
) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepository:   ordersRepository,
		productsRepository: productsRepository,
	}
}
