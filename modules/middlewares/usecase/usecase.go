package middlewareUsecases

import (
	middlewareRepositories "github.com/korvised/go-ecommerce/modules/middlewares/repositories"
)

type IMiddlewareUsecase interface {
}

type middlewareUsecase struct {
	middlewareRepository middlewareRepositories.IMiddlewareRepository
}

func MiddlewareUsecase(middlewareRepository middlewareRepositories.IMiddlewareRepository) IMiddlewareUsecase {
	return &middlewareUsecase{
		middlewareRepository: middlewareRepository,
	}
}
