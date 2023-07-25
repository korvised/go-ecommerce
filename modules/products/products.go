package products

import (
	"github.com/korvised/go-ecommerce/modules/appinfo"
	"github.com/korvised/go-ecommerce/modules/entities"
)

type Product struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Category    *appinfo.Category `json:"category"`
	Images      []*entities.Image `json:"images"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

type ProductFilter struct {
	ID     string `query:"id"`
	Search string `query:"search"`
	*entities.PaginationReq
	*entities.SortReq
}
