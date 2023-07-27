package ordersRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/korvised/go-ecommerce/modules/orders"
	"github.com/korvised/go-ecommerce/modules/orders/ordersPatterns"
	"strings"
	"time"
)

type IOrdersRepository interface {
	FindOneOrder(orderID string) (*orders.Order, error)
	FindManyOrders(req *orders.OrderFilter) ([]*orders.Order, int)
	InsertOrder(req *orders.Order) (string, error)
	UpdateOrder(req *orders.UpdateOrderReq) error
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
		Products: make([]*orders.ProductsOrder, 0),
	}

	if err := r.db.GetContext(ctx, &orderBytes, query, orderID); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(orderBytes, &order); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}

	return order, nil
}

func (r *ordersRepository) FindManyOrders(req *orders.OrderFilter) ([]*orders.Order, int) {
	builder := ordersPatterns.FindOrderBuilder(r.db, req)
	engineer := ordersPatterns.FindOrderEngineer(builder)

	return engineer.FindOrder(), engineer.CountOrder()
}

func (r *ordersRepository) InsertOrder(req *orders.Order) (string, error) {
	builder := ordersPatterns.InsertOrderBuilder(r.db, req)
	engineer := ordersPatterns.InsertOrderEngineer(builder)

	return engineer.InsertOrder()
}

func (r *ordersRepository) UpdateOrder(req *orders.UpdateOrderReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	UPDATE orders SET `

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	if req.Status != "" {
		values = append(values, req.Status)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		status = $%d ?`, lastIndex))

		lastIndex++
	}

	if req.TransferSlip != nil {
		values = append(values, req.TransferSlip)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		transfer_slip = $%d ?`, lastIndex))

		lastIndex++
	}

	values = append(values, req.ID)

	queryClose := fmt.Sprintf(`
	WHERE id = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)

		}
	}

	query += queryClose

	if _, err := r.db.ExecContext(ctx, query, values...); err != nil {
		return fmt.Errorf("update order failed: %v", err)
	}

	return nil
}
