package productsPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/korvised/go-ecommerce/modules/products"
	"github.com/korvised/go-ecommerce/pkg/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

type IFindProductsBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*products.Product
	Count() int
	PrintQuery()
}

type findProductsBuilder struct {
	db             *sqlx.DB
	req            *products.ProductFilter
	query          string
	lastStackIndex int
	values         []any
}

func FindProductsBuild(db *sqlx.DB, req *products.ProductFilter) IFindProductsBuilder {
	return &findProductsBuilder{
		db:  db,
		req: req,
	}
}

func (b *findProductsBuilder) openJsonQuery() {
	b.query += `
	 SELECT 
    	array_to_json(array_agg("t")) 
	 FROM (
	`
}

func (b *findProductsBuilder) initQuery() {
	b.query += `
	 SELECT p.id,
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
	 WHERE 1 = 1
	`
}

func (b *findProductsBuilder) countQuery() {
	b.query += `
	 SELECT count(*) AS count
	 FROM products p
	 WHERE 1 = 1
	`
}

func (b *findProductsBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// ID check
	if b.req.ID != "" {
		b.values = append(b.values, b.req.ID)

		query := `
		 AND p.id = ?
		`
		queryWhereStack = append(queryWhereStack, query)
	}

	// Search check
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		queryWhereStack = append(queryWhereStack, `
		AND (LOWER(p.title) LIKE ? OR LOWER(p.description) LIKE ?)`)
	}

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
		} else {
			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
			queryWhere = strings.Replace(queryWhere, "?", "$"+strconv.Itoa(i+2), 1)
		}
	}

	// Last stack record
	b.lastStackIndex = len(b.values)
	// Summary query
	b.query += queryWhere
}

func (b *findProductsBuilder) sort() {
	orderByReq := strings.ToLower(b.req.OrderBy)
	orderByMap := map[string]string{
		"id":    "p.id",
		"title": "p.title",
		"price": "p.price",
	}

	if orderByMap[orderByReq] == "" {
		b.req.OrderBy = orderByMap["title"]
	} else {
		b.req.OrderBy = orderByMap[orderByReq]
	}

	sortReq := strings.ToUpper(b.req.Sort)
	sortMap := map[string]string{
		"ASC":  "ASC",
		"DESC": "DESC",
	}

	if sortMap[sortReq] == "" {
		b.req.Sort = sortMap["title"]
	} else {
		b.req.Sort = sortMap[sortReq]
	}

	b.values = append(b.values, b.req.OrderBy)
	b.query += fmt.Sprintf(`
	 ORDER BY $%d %s`, b.lastStackIndex+1, b.req.Sort)
	b.lastStackIndex = len(b.values)
}

func (b *findProductsBuilder) paginate() {
	// offset (page - 1) * limit (size)
	b.values = append(b.values, (b.req.Page-1)*b.req.Size, b.req.Size)
	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}

func (b *findProductsBuilder) closeJsonQuery() {
	b.query += `
	 ) AS t;
	`
}

func (b *findProductsBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}

func (b *findProductsBuilder) Result() []*products.Product {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	productsBytes := make([]byte, 0)
	productsData := make([]*products.Product, 0)

	b.PrintQuery()

	if err := b.db.Get(&productsBytes, b.query, b.values...); err != nil {
		log.Printf("query products failed: %v\n", err)
		return make([]*products.Product, 0)
	}

	if err := json.Unmarshal(productsBytes, &productsData); err != nil {
		log.Printf("unmarshal products failed: %v\n", err)
		return make([]*products.Product, 0)
	}

	b.resetQuery()
	return productsData
}

func (b *findProductsBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count products failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}

func (b *findProductsBuilder) PrintQuery() {
	utils.Debug(b.values)
	utils.Debug(b.query)
}

type findProductsEngineer struct {
	builder IFindProductsBuilder
}

func FindProductsEngineer(builder IFindProductsBuilder) *findProductsEngineer {
	return &findProductsEngineer{
		builder: builder,
	}
}

func (en *findProductsEngineer) FindProducts() IFindProductsBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findProductsEngineer) CountProduct() IFindProductsBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}
