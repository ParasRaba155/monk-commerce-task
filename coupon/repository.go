package coupon

import (
	"errors"
	"fmt"
)

var (
	ErrDoesNotExist = errors.New("no such entity")
)

// repository is the in memory db
// coupon id will be simply treated as the index in the array
//
// For simplicity I have initialized it with 100 cap
type repository struct {
	coupons []Coupon
	deleted map[int]struct{} // set of deleted IDs
}

func NewRepository() *repository {
	return &repository{
		coupons: make([]Coupon, 0, 100),
		deleted: make(map[int]struct{}, 100),
	}
}

func (r *repository) CreateCoupon(coupon Coupon) error {
	r.coupons = append(r.coupons, coupon)
	return nil
}

func (r *repository) GetAllCoupons() ([]Coupon, error) {
	result := make([]Coupon, 0, len(r.coupons))
	for id, c := range r.coupons {
		if _, isDeleted := r.deleted[id]; !isDeleted {
			result = append(result, c)
		}
	}
	return result, nil
}

// GetCouponByID will return the coupon at that index
//
// If the index is not reachable then it will throw the error of does not exist
func (r *repository) GetCouponByID(id int) (Coupon, error) {
	if id >= len(r.coupons) {
		return Coupon{}, fmt.Errorf("%w: no coupon with id %d", ErrDoesNotExist, id)
	}
	if _, isDeleted := r.deleted[id]; isDeleted {
		return Coupon{}, fmt.Errorf("%w: no coupon with id %d", ErrDoesNotExist, id)
	}
	return r.coupons[id], nil
}

func (r *repository) UpdateCouponByID(id int, newCoupon Coupon) (Coupon, error) {
	if id >= len(r.coupons) {
		return Coupon{}, fmt.Errorf("%w: no coupon with id %d", ErrDoesNotExist, id)
	}
	if _, isDeleted := r.deleted[id]; isDeleted {
		return Coupon{}, fmt.Errorf("%w: no coupon with id %d", ErrDoesNotExist, id)
	}
	r.coupons[id] = newCoupon
	return newCoupon, nil
}

func (r *repository) DeleteCouponByID(id int) error {
	if id >= len(r.coupons) {
		return fmt.Errorf("%w: no coupon with id %d", ErrDoesNotExist, id)
	}
	if _, isDeleted := r.deleted[id]; isDeleted {
		return fmt.Errorf("%w: no coupon with id %d", ErrDoesNotExist, id)
	}
	r.deleted[id] = struct{}{}
	return nil
}
