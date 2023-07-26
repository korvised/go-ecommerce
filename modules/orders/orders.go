package orders

import "github.com/korvised/go-ecommerce/modules/products"

type Orders struct {
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
	Qty     string            `db:"qty" json:"qty"`
	Product *products.Product `db:"product" json:"product"`
}
