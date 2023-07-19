package usersRepositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/korvised/go-ecommerce/modules/users"
	userPatterns "github.com/korvised/go-ecommerce/modules/users/patterns"
)

type IUserRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
}

type userRepository struct {
	db *sqlx.DB
}

func UserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := userPatterns.InsertUser(r.db, req, isAdmin)

	var err error
	if isAdmin {
		result, err = result.Admin()
	} else {
		result, err = result.Customer()
	}

	// Get result from inserting
	user, err := result.Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}
