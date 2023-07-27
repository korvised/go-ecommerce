package mytests

import (
	"database/sql"
	"testing"
)

type testFindOneProduct struct {
	productID string
	isErr     bool
	expect    string
}

func TestFindOneProduct(t *testing.T) {
	tests := []testFindOneProduct{
		{
			productID: "P000099",
			isErr:     true,
			expect:    "get product failed: " + sql.ErrNoRows.Error(),
		},
		{
			productID: "P000001",
			isErr:     false,
			expect:    "",
		},
	}

	productsModule := SetupTest().ProductsModule()

	for _, test := range tests {
		if test.isErr {
			if _, err := productsModule.Usecase().FindOneProduct(test.productID); err.Error() != test.expect {
				t.Errorf("expect: %s, got: %s", test.expect, err.Error())
			}

		} else {
			product, err := productsModule.Usecase().FindOneProduct(test.productID)
			if err != nil {
				t.Errorf("expect: %v, got: %s", nil, err.Error())
			}

			if CompressToJson(&product) != test.expect {
				t.Errorf("expect: %v, got: %s", test.expect, CompressToJson(&product))
			}
		}

	}
}
