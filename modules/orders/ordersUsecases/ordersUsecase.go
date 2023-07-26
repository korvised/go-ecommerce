package ordersUsecases

import (
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/orders"
	"github.com/korvised/go-ecommerce/modules/orders/ordersRepositories"
	"github.com/korvised/go-ecommerce/modules/products/productsRepositories"
	"math"
)

type IOrdersUsecase interface {
	FindOneOrder(orderID string) (*orders.Order, error)
	FindManyOrders(req *orders.OrderFilter) *entities.PaginateRes
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

func (u *ordersUsecase) FindOneOrder(orderID string) (*orders.Order, error) {
	return u.ordersRepository.FindOneOrder(orderID)
}

func (u *ordersUsecase) FindManyOrders(req *orders.OrderFilter) *entities.PaginateRes {
	data, count := u.ordersRepository.FindManyOrders(req)

	return &entities.PaginateRes{
		Page:      req.Page,
		Size:      req.Size,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Size))),
		TotalItem: count,
		Data:      data,
	}
}
