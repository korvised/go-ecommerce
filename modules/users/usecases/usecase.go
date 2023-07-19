package usersUsecases

import (
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/users"
	userRepositories "github.com/korvised/go-ecommerce/modules/users/repositories"
)

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type userUsecase struct {
	cfg            config.IConfig
	userRepository userRepositories.IUserRepository
}

func UserUsecase(cfg config.IConfig, userRepository userRepositories.IUserRepository) IUserUsecase {
	return &userUsecase{
		cfg:            cfg,
		userRepository: userRepository,
	}
}

func (u userUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.userRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}
