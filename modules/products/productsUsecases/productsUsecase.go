package productsUsecases

import (
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/products"
	"github.com/korvised/go-ecommerce/modules/products/productsRepositories"
	"math"
)

type IProductsUsecase interface {
	FindOneProduct(productID string) (*products.Product, error)
	FindManyProducts(req *products.ProductFilter) *entities.PaginateRes
	AddProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productID string) error
}

type productsUsecase struct {
	productsRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productsRepository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productsRepository: productsRepository,
	}
}

func (u *productsUsecase) FindOneProduct(productID string) (*products.Product, error) {
	return u.productsRepository.FindOneProduct(productID)
}

func (u *productsUsecase) FindManyProducts(req *products.ProductFilter) *entities.PaginateRes {
	data, count := u.productsRepository.FindManyProducts(req)

	return &entities.PaginateRes{
		Data:      data,
		Page:      req.Page,
		Size:      req.Size,
		TotalPage: int(math.Ceil(float64(count) / (float64(req.Size)))),
		TotalItem: count,
	}
}

func (u *productsUsecase) AddProduct(req *products.Product) (*products.Product, error) {
	return u.productsRepository.InsertProduct(req)
}

func (u *productsUsecase) UpdateProduct(req *products.Product) (*products.Product, error) {
	return u.productsRepository.UpdateProduct(req)
}

func (u *productsUsecase) DeleteProduct(productID string) error {
	return u.productsRepository.DeleteProduct(productID)
}
