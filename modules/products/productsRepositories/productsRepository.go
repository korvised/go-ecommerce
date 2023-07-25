package productsRepositories

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/files/filesUsecases"
	"github.com/korvised/go-ecommerce/modules/products"
	"github.com/korvised/go-ecommerce/modules/products/productsPatterns"
)

type IProductsRepository interface {
	FindOneProduct(productID string) (*products.Product, error)
	FindManyProducts(req *products.ProductFilter) ([]*products.Product, int)
}

type productsRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func ProductsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IProductsRepository {
	return &productsRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *productsRepository) FindOneProduct(productID string) (*products.Product, error) {
	query := `
	 SELECT to_jsonb(t)
	 FROM (SELECT p.id,
             p.title,
             p.description,
             p.price,
             (SELECT to_jsonb(ct)
              FROM (SELECT c.id,
                           c.title
                    FROM categories c
                             LEFT JOIN products_categories pc ON pc.category_id = c.id
                    WHERE pc.product_id = p.id) AS ct) AS category,
             (SELECT COALESCE(array_to_json(array_agg(it)), '[]'::json)
              FROM (SELECT i.id,
                           i.filename,
                           i.url
                    FROM images i
                    WHERE i.product_id = p.id) AS it)  AS images,
             p.created_at,
             p.updated_at

      FROM products p
      WHERE p.id = $1
      LIMIT 1) AS t;
	`

	productByte := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&productByte, query, productID); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(productByte, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product failed: %v", err)
	}

	return product, nil
}

func (r *productsRepository) FindManyProducts(req *products.ProductFilter) ([]*products.Product, int) {
	builder := productsPatterns.FindProductBuilder(r.db, req)
	engineer := productsPatterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()

	return result, count
}
