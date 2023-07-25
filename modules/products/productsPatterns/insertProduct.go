package productsPatterns

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/korvised/go-ecommerce/modules/products"
	"time"
)

type IInsertProductBuilder interface {
	initTransaction() error
	insertProduct() error
	insertCategory() error
	insertAttachment() error
	commit() error
	getProductId() string
}

type insertProductBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *products.Product
}

func InsertProductBuilder(db *sqlx.DB, req *products.Product) IInsertProductBuilder {
	return &insertProductBuilder{
		db:  db,
		req: req,
	}
}

type insertProductEngineer struct {
	builder IInsertProductBuilder
}

func (b *insertProductBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertProductBuilder) insertProduct() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "products" (
		"title",
		"description",
		"price"
	)
	VALUES ($1, $2, $3)
		RETURNING "id";`

	if err := b.tx.QueryRowContext(
		ctx,
		query,
		b.req.Title,
		b.req.Description,
		b.req.Price,
	).Scan(&b.req.ID); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert product failed: %v", err)
	}

	return nil
}

func (b *insertProductBuilder) insertCategory() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "products_categories" (
		"product_id",
		"category_id"
	)
	VALUES ($1, $2);`

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		b.req.ID,
		b.req.Category.ID,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert products_categories failed: %v", err)
	}

	return nil
}

func (b *insertProductBuilder) insertAttachment() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "images" (
		"filename",
		"url",
		"product_id"
	)
	VALUES`

	valueStack := make([]any, 0)
	var index int
	for i, image := range b.req.Images {
		valueStack = append(valueStack, image.FileName, image.Url, b.req.ID)

		if i != len(b.req.Images)-1 {
			query += fmt.Sprintf(`
			( $%d, $%d, $%d )`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			( $%d, $%d, $%d );`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}

	return nil
}

func (b *insertProductBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (b *insertProductBuilder) getProductId() string {
	return b.req.ID
}

func InsertProductEngineer(b IInsertProductBuilder) *insertProductEngineer {
	return &insertProductEngineer{builder: b}
}

func (en *insertProductEngineer) InsertProduct() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}

	if err := en.builder.insertProduct(); err != nil {
		return "", err
	}

	if err := en.builder.insertCategory(); err != nil {
		return "", err
	}

	if err := en.builder.insertAttachment(); err != nil {
		return "", err
	}

	if err := en.builder.commit(); err != nil {
		return "", err
	}

	return en.builder.getProductId(), nil
}
