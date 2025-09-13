package coupon

import (
	"errors"
	"fmt"
)

var (
	errDoesNotExist = errors.New("no such entity")
)

// repository is the in memory db
// coupon id will be simply treated as the index in the array
//
// For simplicity I have initialized it with 100 cap
type repository struct {
	coupons []Coupon
}

func NewRepository() *repository {
	return &repository{
		coupons: make([]Coupon, 0, 100),
	}
}

func (r *repository) CreateCoupon(coupon Coupon) error {
	r.coupons = append(r.coupons, coupon)
	return nil
}

func (r *repository) GetAllCoupons() ([]Coupon, error) {
	return r.coupons, nil
}

// GetCouponByID will return the coupon at that index
//
// If the index is not reachable then it will throw the error of does not exist
func (r *repository) GetCouponByID(id int) (Coupon, error) {
	if id >= len(r.coupons) {
		return Coupon{}, fmt.Errorf("%w: no coupon with id %d", errDoesNotExist, id)
	}
	return r.coupons[id], nil
}

func (r *repository) UpdateCouponByID(id int, newCoupon Coupon) (Coupon, error) {
	if id >= len(r.coupons) {
		return Coupon{}, fmt.Errorf("%w: no coupon with id %d", errDoesNotExist, id)
	}
	r.coupons[id] = newCoupon
	return newCoupon, nil
}
