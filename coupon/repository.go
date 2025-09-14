package coupon

import (
	"errors"
	"fmt"
)

var (
	ErrDoesNotExist = errors.New("no such entity")
)

// repository is the in-memory db
// coupons are stored by coupon.ID
type repository struct {
	coupons map[int]Coupon
	nextID  int // auto-incrementing ID counter
}

func NewRepository() *repository {
	return &repository{
		coupons: make(map[int]Coupon, 100),
		nextID:  0,
	}
}

// CreateCoupon assigns a new ID and stores the coupon.
func (r *repository) CreateCoupon(coupon Coupon) error {
	coupon.ID = r.nextID
	r.coupons[coupon.ID] = coupon
	r.nextID++
	return nil
}

// GetAllCoupons returns all coupons currently in the repository.
func (r *repository) GetAllCoupons() ([]Coupon, error) {
	result := make([]Coupon, 0, len(r.coupons))
	for _, c := range r.coupons {
		result = append(result, c)
	}
	return result, nil
}

// GetCouponByID returns the coupon with the given ID.
func (r *repository) GetCouponByID(id int) (Coupon, error) {
	c, ok := r.coupons[id]
	if !ok {
		return Coupon{}, fmt.Errorf("%w: no coupon with id %d", ErrDoesNotExist, id)
	}
	return c, nil
}

// UpdateCouponByID replaces the coupon with the new details.
func (r *repository) UpdateCouponByID(id int, newCoupon Coupon) (Coupon, error) {
	_, ok := r.coupons[id]
	if !ok {
		return Coupon{}, fmt.Errorf("%w: no coupon with id %d", ErrDoesNotExist, id)
	}
	newCoupon.ID = id // enforce correct ID
	r.coupons[id] = newCoupon
	return newCoupon, nil
}

// DeleteCouponByID removes the coupon from the repository.
func (r *repository) DeleteCouponByID(id int) error {
	_, ok := r.coupons[id]
	if !ok {
		return fmt.Errorf("%w: no coupon with id %d", ErrDoesNotExist, id)
	}
	delete(r.coupons, id)
	return nil
}
