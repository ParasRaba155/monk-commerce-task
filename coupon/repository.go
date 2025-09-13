package coupon

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
