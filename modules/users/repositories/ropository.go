package usersRepositories

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/korvised/go-ecommerce/modules/users"
	usersPatterns "github.com/korvised/go-ecommerce/modules/users/patterns"
	"time"
)

type IUsersRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserPassport) error
}

type usersRepository struct {
	db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := usersPatterns.InsertUser(r.db, req, isAdmin)

	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	// Get result from inserting
	user, err := result.Result()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *usersRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
	 SELECT id, email, password, username, role_id
	 FROM users
	 WHERE email = $1;
	`

	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return user, nil
}

func (r *usersRepository) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `
	 INSERT INTO oauth (user_id, access_token, refresh_token)
	 VALUES ($1, $2, $3)
	 RETURNING "id";
	`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.Id,
		req.Token.AccessToken,
		req.Token.RefreshToken,
	).Scan(&req.Token.Id); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}

	return nil
}
