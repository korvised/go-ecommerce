package orders

import (
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/products"
)

type OrderFilter struct {
	Search    string `query:"search"` // user_id, address, contact
	Status    string `query:"status"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	*entities.PaginationReq
	*entities.SortReq
}

type Order struct {
	ID           string           `db:"id" json:"id"`
	UserID       string           `db:"user_id" json:"user_id"`
	TransferSlip *TransferSlip    `db:"transfer_slip" json:"transfer_slip"`
	Products     []*ProductsOrder `json:"products"`
	Address      string           `db:"address" json:"address"`
	Contact      string           `db:"contact" json:"contact"`
	Status       string           `db:"status" json:"status"`
	TotalPaid    float64          `db:"total_paid" json:"total_paid"`
	CreatedAt    string           `db:"created_at" json:"created_at"`
	UpdatedAt    string           `db:"updated_at" json:"updated_at"`
}

type TransferSlip struct {
	ID        string `json:"id"`
	FileName  string `json:"fileName"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type ProductsOrder struct {
	ID      string            `db:"id" json:"id"`
	Qty     int               `db:"qty" json:"qty"`
	Product *products.Product `db:"product" json:"product"`
}

type UpdateOrderReq struct {
	ID           string        `form:"id" json:"id"`
	TransferSlip *TransferSlip `form:"transfer_slip" json:"transfer_slip"`
	Status       string        `form:"status" json:"status"`
}
