package ordersRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/korvised/go-ecommerce/modules/orders"
	"time"
)

type IOrdersRepository interface {
	FindOneOrder(orderID string) (*orders.Order, error)
}

type ordersRepository struct {
	db *sqlx.DB
}

func OrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{db: db}
}

func (r *ordersRepository) FindOneOrder(orderID string) (*orders.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	SELECT to_jsonb(t)
	FROM (SELECT o.id,
				 o.user_id,
				 o.transfer_slip,
				 (SELECT array_to_json(array_agg(pt))
				  FROM (SELECT spo.id, spo.qty, spo.product
						FROM products_orders spo
						WHERE spo.order_id = o.id) AS pt) AS products,
				 o.address,
				 o.contact,
				 o.status,
				 (SELECT SUM(COALESCE((po.product ->> 'price')::FLOAT * (po.qty)::FLOAT, 0))
				  FROM products_orders po
				  WHERE po.order_id = o.id)               AS total_paid,
				 o.created_at,
				 o.updated_at
		  FROM orders o
		  WHERE o.id = $1) AS t;
	`

	orderBytes := make([]byte, 0)
	order := &orders.Order{
		TransferSlip: &orders.TransferSlip{},
		Products:     make([]*orders.ProductsOrder, 0),
	}

	if err := r.db.GetContext(ctx, &orderBytes, query, orderID); err != nil {
		return nil, err
	}

	fmt.Printf(string(orderBytes))

	if err := json.Unmarshal(orderBytes, &order); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}

	return order, nil
}
