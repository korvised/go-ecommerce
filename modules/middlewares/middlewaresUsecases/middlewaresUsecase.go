package middlewaresUsecases

import (
	"github.com/korvised/go-ecommerce/modules/middlewares"
	"github.com/korvised/go-ecommerce/modules/middlewares/middlewaresRepositories"
)

type IMiddlewareUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRoles() ([]*middlewares.Role, error)
}

type middlewareUsecase struct {
	middlewareRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewareUsecase(middlewareRepository middlewaresRepositories.IMiddlewaresRepository) IMiddlewareUsecase {
	return &middlewareUsecase{
		middlewareRepository: middlewareRepository,
	}
}

func (u *middlewareUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewareRepository.FindAccessToken(userId, accessToken)
}

func (u *middlewareUsecase) FindRoles() ([]*middlewares.Role, error) {
	return u.middlewareRepository.FindRole()
}
