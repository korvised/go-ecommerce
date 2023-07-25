package productsRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/files/filesUsecases"
	"github.com/korvised/go-ecommerce/modules/products"
	"github.com/korvised/go-ecommerce/modules/products/productsPatterns"
	"time"
)

type IProductsRepository interface {
	FindOneProduct(productID string) (*products.Product, error)
	FindManyProducts(req *products.ProductFilter) ([]*products.Product, int)
	InsertProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productID string) error
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

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

	productBytes := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.GetContext(ctx, &productBytes, query, productID); err != nil {
		return nil, fmt.Errorf("get product failed: %v", err)
	}
	if err := json.Unmarshal(productBytes, &product); err != nil {
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

func (r *productsRepository) InsertProduct(req *products.Product) (*products.Product, error) {
	builder := productsPatterns.InsertProductBuilder(r.db, req)
	productID, err := productsPatterns.InsertProductEngineer(builder).InsertProduct()
	if err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productsRepository) UpdateProduct(req *products.Product) (*products.Product, error) {
	builder := productsPatterns.UpdateProductBuilder(r.db, req, r.filesUsecase)
	engineer := productsPatterns.UpdateProductEngineer(builder)

	if err := engineer.UpdateProduct(); err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(req.ID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productsRepository) DeleteProduct(productID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `DELETE FROM products WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, query, productID); err != nil {
		return fmt.Errorf("delete product failed: %v", err)
	}

	return nil
}
