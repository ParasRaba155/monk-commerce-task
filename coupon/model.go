// Package coupon to handle everything related to the coupon entity
//
// Including DB and endpoints
package coupon

import (
	"errors"
	"fmt"
)

type CouponType string

var (
	errInvalidThreshold   = errors.New("invalid threshold")
	errInvalidDiscount    = errors.New("invalid discount")
	errInvalidProductList = errors.New("invalid product list")
	errInvalidRepition    = errors.New("invalid repetition limit")
)

// couponTypes for all the possible couponTypes
//
// NOTE: We have added limited couponTypes here, we can add more types in this slice
// In real world this would have been kept it in a separate db table/collection
// And slice would be better choice since we want it to be extensible
var couponTypes = [...]CouponType{
	"cart-wise",
	"product-wise",
	"bxgy",
}

type CouponDetails interface {
	GetCouponType() CouponType
	ValidateCoupon() error
}

type CartWiseDetails struct {
	Threshold int `json:"threshold"`
	Discount  int `json:"discount"`
}

func (CartWiseDetails) GetCouponType() CouponType {
	return couponTypes[0]
}

func (c CartWiseDetails) ValidateCoupon() error {
	if c.Threshold < 0 {
		return fmt.Errorf("%w, threshold must be positive", errInvalidThreshold)
	}
	if c.Discount < 0 || c.Discount > 100 {
		return fmt.Errorf("%w, discount must be between 0 and 100%%", errInvalidDiscount)
	}
	return nil
}

type ProductWiseDetails struct {
	ProductID int `json:"product_id"`
	Discount  int `json:"discount"`
}

func (ProductWiseDetails) GetCouponType() CouponType {
	return couponTypes[1]
}

func (c ProductWiseDetails) ValidateCoupon() error {
	if c.Discount < 0 || c.Discount > 100 {
		return fmt.Errorf("%w, discount must be between 0 and 100%%", errInvalidDiscount)
	}
	return nil
}

type BxGyDetails struct {
	BuyProducts     []CouponProduct `json:"buy_products"`
	GetProducts     []CouponProduct `json:"get_products"`
	RepetitionLimit int             `json:"repition_limit"`
}

func (BxGyDetails) GetCouponType() CouponType {
	return couponTypes[2]
}

func (c BxGyDetails) ValidateCoupon() error {
	if len(c.BuyProducts) == 0 || len(c.GetProducts) == 0 {
		return fmt.Errorf("%w: buy or get product list can not be empty", errInvalidProductList)
	}

	for _, prod := range c.BuyProducts {
		if prod.Quantity < 1 {
			return fmt.Errorf("%w: quantity should be positive", errInvalidProductList)
		}
	}

	for _, prod := range c.GetProducts {
		if prod.Quantity < 1 {
			return fmt.Errorf("%w: quantity should be positive", errInvalidProductList)
		}
	}

	if c.RepetitionLimit < 1 {
		return fmt.Errorf("%w: repetition limit should be greater than 0", errInvalidRepition)
	}
	return nil
}

type CouponProduct struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Coupon struct {
	ID      int
	Type    CouponType
	Details CouponDetails
}
