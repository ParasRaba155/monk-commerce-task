package cart

import (
	"errors"
	"fmt"

	"github.com/ParasRaba155/monk-commerce-task/coupon"
)

var errInvalidQuantity = errors.New("invalid quantity")

// Item in cart will have a price but that will be determined by BE
type Item struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type PricedItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	Price     int `json:"price"`
}

func (i Item) ToPricedItem(price int) PricedItem {
	return PricedItem{
		ProductID: i.ProductID,
		Quantity:  i.Quantity,
		Price:     price,
	}
}

type DiscountedItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	Price     int `json:"price"`
	Discount  int `json:"discount"`
}

func (i Item) ToDiscountedItem(price, discount int) DiscountedItem {
	return DiscountedItem{
		ProductID: i.ProductID,
		Quantity:  i.Quantity,
		Price:     price,
		Discount:  discount,
	}
}

func (i PricedItem) ToDiscountedItem(discount int) DiscountedItem {
	return DiscountedItem{
		ProductID: i.ProductID,
		Quantity:  i.Quantity,
		Price:     i.Price,
		Discount:  discount,
	}
}

type DiscountCoupon struct {
	CouponID int               `json:"coupon_id"`
	Type     coupon.CouponType `json:"type"`
	Discount int               `json:"discount"`
}

type Cart struct {
	Items []Item `json:"items"`
}

type DiscountedCart struct {
	Items         []DiscountedItem `json:"items"`
	TotalPrice    int              `json:"total_price"`
	TotalDiscount int              `json:"total_discount"`
	FinalPrice    int              `json:"final_price"`
}

// Validate will check for >= 1 quantity
func (c Cart) Validate() error {
	for _, item := range c.Items {
		if item.Quantity < 1 {
			return fmt.Errorf("%w: quantity should be positive", errInvalidQuantity)
		}
	}
	return nil
}
