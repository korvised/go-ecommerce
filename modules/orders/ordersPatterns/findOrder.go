package ordersPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/korvised/go-ecommerce/modules/orders"
	"log"
	"strings"
	"time"
)

type IFindOrderBuilder interface {
	initQuery()
	initCountQuery()
	buildWhereSearch()
	buildWhereStatus()
	buildWhereDate()
	buildSort()
	buildPaginate()
	closeQuery()
	getQuery() string
	setQuery(query string)
	getValues() []any
	setValues(data []any)
	setLastIndex(n int)
	getDb() *sqlx.DB
	reset()
}

type findOrderBuilder struct {
	db        *sqlx.DB
	req       *orders.OrderFilter
	query     string
	values    []any
	lastIndex int
}

func FindOrderBuilder(db *sqlx.DB, req *orders.OrderFilter) IFindOrderBuilder {
	return &findOrderBuilder{
		db:     db,
		req:    req,
		values: make([]any, 0),
	}
}

type findOrderEngineer struct {
	builder IFindOrderBuilder
}

func FindOrderEngineer(b IFindOrderBuilder) *findOrderEngineer {
	return &findOrderEngineer{builder: b}
}

func (b *findOrderBuilder) initQuery() {
	b.query += `
	SELECT array_to_json(array_agg(at))
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
		  WHERE 1 = 1`
}

func (b *findOrderBuilder) initCountQuery() {
	b.query += `
	SELECT COUNT(*) AS count
	FROM orders o
	WHERE 1 = 1
	`
}

func (b *findOrderBuilder) buildWhereSearch() {
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		query := fmt.Sprintf(`
		AND (
				LOWER(o.user_id) LIKE $%d OR
				LOWER(o.address) LIKE $%d OR
				LOWER(o.contact) LIKE $%d
			)`,
			b.lastIndex+1,
			b.lastIndex+2,
			b.lastIndex+3,
		)

		temp := b.getQuery()
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findOrderBuilder) buildWhereStatus() {
	if b.req.Status != "" {
		b.values = append(b.values, strings.ToLower(b.req.Status))

		query := fmt.Sprintf(`
		AND o.status = $%d`, b.lastIndex+1)

		temp := b.getQuery()
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findOrderBuilder) buildWhereDate() {
	if b.req.StartDate != "" && b.req.EndDate != "" {
		b.values = append(b.values, b.req.StartDate, b.req.EndDate)

		query := fmt.Sprintf(`
		 AND (o.created_at BETWEEN DATE($%d) AND DATE($%d)::DATE + 1)`, b.lastIndex+1, b.lastIndex+2)

		temp := b.getQuery()
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findOrderBuilder) buildSort() {
	b.values = append(b.values, b.req.OrderBy)

	b.query += fmt.Sprintf(`
	ORDER BY $%d %s`, b.lastIndex+1, b.req.Sort)

	b.lastIndex = len(b.values)
}

func (b *findOrderBuilder) buildPaginate() {
	b.values = append(b.values, (b.req.Page-1)*b.req.Size, b.req.Size)

	b.query += fmt.Sprintf(`
	OFFSET $%d LIMIT $%d `, b.lastIndex+1, b.lastIndex+2)

	b.lastIndex = len(b.values)
}

func (b *findOrderBuilder) closeQuery() {
	b.query += `
	) AS at;`
}

func (b *findOrderBuilder) getQuery() string { return b.query }

func (b *findOrderBuilder) setQuery(query string) { b.query = query }

func (b *findOrderBuilder) getValues() []any { return b.values }

func (b *findOrderBuilder) setValues(data []any) { b.values = data }

func (b *findOrderBuilder) setLastIndex(n int) { b.lastIndex = n }

func (b *findOrderBuilder) getDb() *sqlx.DB { return b.db }

func (b *findOrderBuilder) reset() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastIndex = 0
}

func (en *findOrderEngineer) FindOrder() []*orders.Order {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	en.builder.initQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereStatus()
	en.builder.buildWhereDate()
	en.builder.buildSort()
	en.builder.buildPaginate()
	en.builder.closeQuery()

	ordersData := make([]*orders.Order, 0)
	ordersBytes := make([]byte, 0)

	fmt.Println(en.builder.getValues())
	fmt.Println(en.builder.getQuery())

	err := en.builder.getDb().GetContext(ctx, &ordersBytes, en.builder.getQuery(), en.builder.getValues()...)
	if err != nil {
		log.Printf("query orders failed: %v", err)
		return ordersData
	}

	en.builder.reset()

	if err := json.Unmarshal(ordersBytes, &ordersData); err != nil {
		log.Printf("unmarshal order failed: %v", err)
	}

	return ordersData
}

func (en *findOrderEngineer) CountOrder() int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	en.builder.initCountQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereStatus()
	en.builder.buildWhereDate()

	var count int
	err := en.builder.getDb().GetContext(ctx, &count, en.builder.getQuery(), en.builder.getValues()...)
	if err != nil {
		log.Printf("count order failed: %v", err)
		return 0
	}

	en.builder.reset()

	return count
}
