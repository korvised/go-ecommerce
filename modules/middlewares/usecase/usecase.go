package middlewareUsecases

import (
	"github.com/korvised/go-ecommerce/modules/middlewares"
	middlewareRepositories "github.com/korvised/go-ecommerce/modules/middlewares/repositories"
)

type IMiddlewareUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRoles() ([]*middlewares.Role, error)
}

type middlewareUsecase struct {
	middlewareRepository middlewareRepositories.IMiddlewaresRepository
}

func MiddlewareUsecase(middlewareRepository middlewareRepositories.IMiddlewaresRepository) IMiddlewareUsecase {
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
