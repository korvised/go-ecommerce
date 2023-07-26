package ordersPatterns

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/korvised/go-ecommerce/modules/orders"
	"time"
)

type IInsertOrderBuilder interface {
	initTransaction() error
	insertOrder() error
	insertProductsOrder() error
	getOrderId() string
	commit() error
}

type insertOrderBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *orders.Order
}

type insertOrderEngineer struct {
	builder IInsertOrderBuilder
}

func (b *insertOrderBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertOrderBuilder) insertOrder() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `
	INSERT INTO orders (user_id, address, contact, transfer_slip, status)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;`

	if err := b.tx.QueryRowContext(
		ctx,
		query,
		b.req.UserID,
		b.req.Address,
		b.req.Contact,
		b.req.TransferSlip,
		b.req.Status,
	).Scan(&b.req.ID); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert order failed: %v", err)
	}

	return nil
}

func (b *insertOrderBuilder) insertProductsOrder() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `
	INSERT INTO products_orders (order_id, qty, product)
	VALUES`

	values := make([]any, 0)
	lastIndex := 0
	for i, pro := range b.req.Products {
		values = append(values, b.req.ID, pro.Qty, pro.Product)

		if i != len(b.req.Products)-1 {
			query += fmt.Sprintf(`
			( $%d, $%d, $%d ),`, lastIndex+1, lastIndex+2, lastIndex+3)
		} else {
			query += fmt.Sprintf(`
			( $%d, $%d, $%d );`, lastIndex+1, lastIndex+2, lastIndex+3)
		}

		lastIndex += 3
	}

	if _, err := b.tx.ExecContext(ctx, query, values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert product order failed: %v", err)
	}

	return nil
}

func (b *insertOrderBuilder) getOrderId() string { return b.req.ID }

func (b *insertOrderBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func InsertOrderBuilder(db *sqlx.DB, req *orders.Order) IInsertOrderBuilder {
	return &insertOrderBuilder{
		db:  db,
		req: req,
	}
}

func InsertOrderEngineer(b IInsertOrderBuilder) *insertOrderEngineer {
	return &insertOrderEngineer{builder: b}
}

func (en *insertOrderEngineer) InsertOrder() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}

	if err := en.builder.insertOrder(); err != nil {
		return "", err
	}

	if err := en.builder.insertProductsOrder(); err != nil {
		return "", err
	}

	if err := en.builder.commit(); err != nil {
		return "", err
	}

	return en.builder.getOrderId(), nil
}
