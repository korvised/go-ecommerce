package ordersUsecases

import (
	"fmt"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/orders"
	"github.com/korvised/go-ecommerce/modules/orders/ordersRepositories"
	"github.com/korvised/go-ecommerce/modules/products/productsRepositories"
	"math"
)

type IOrdersUsecase interface {
	FindOneOrder(orderID string) (*orders.Order, error)
	FindManyOrders(req *orders.OrderFilter) *entities.PaginateRes
	InsertOrder(req *orders.Order) (*orders.Order, error)
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

func (u *ordersUsecase) InsertOrder(req *orders.Order) (*orders.Order, error) {
	// Check if product is exits
	for i, pro := range req.Products {
		if pro.Product == nil {
			return nil, fmt.Errorf("product %d is empty", i+1)
		}

		product, err := u.productsRepository.FindOneProduct(pro.Product.ID)
		if err != nil {
			return nil, err
		}

		// Summary price
		req.TotalPaid += pro.Product.Price * float64(pro.Qty)
		req.Products[i].Product = product
	}

	orderID, err := u.ordersRepository.InsertOrder(req)
	if err != nil {
		return nil, err
	}

	order, err := u.ordersRepository.FindOneOrder(orderID)
	if err != nil {
		return nil, err
	}

	return order, nil
}
