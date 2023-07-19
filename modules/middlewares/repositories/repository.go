package middlewareRepositories

import "github.com/jmoiron/sqlx"

type IMiddlewareRepository interface {
}

type middlewareRepository struct {
	db *sqlx.DB
}

func MiddlewareRepository(db *sqlx.DB) IMiddlewareRepository {
	return &middlewareRepository{
		db: db,
	}
}
