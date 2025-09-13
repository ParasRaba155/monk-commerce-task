// Package cart handles the everything related to card and overall product list
package cart

import (
	"errors"
	"fmt"
)

var errInvalidProductID = errors.New("invalid product id")

// getProductPrice: This would act as our db layer for product metadata
// We are assuming that we will have unlimited quantity, so user can specify
// but we only have 10 products from id 1 to 10
// and price of each product is productID * 10
func getProductPrice(productID int) (int, error) {
	if productID < 1 || productID > 10 {
		return 0, fmt.Errorf("%w: invalid id %d", errInvalidProductID, productID)
	}

	return productID * 10, nil
}
