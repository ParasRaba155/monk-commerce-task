package cart

import (
	"reflect"
	"testing"

	"github.com/ParasRaba155/monk-commerce-task/coupon"
)

func TestApplyBxGyWiseCoupon(t *testing.T) {
	const (
		productXID = 1
		productYID = 2
		productZID = 3
		productAID = 4
		productBID = 5
		productCID = 6
	)

	tests := []struct {
		name         string
		items        []PricedItem
		totalPrice   int
		coupon       coupon.BxGyDetails
		expectedCart DiscountedCart
	}{
		{
			name: "Valid B2G1 with one discount",
			items: []PricedItem{
				{ProductID: productXID, Quantity: 1, Price: 10},
				{ProductID: productYID, Quantity: 1, Price: 12},
				{ProductID: productAID, Quantity: 1, Price: 8},
			},
			totalPrice: 30,
			coupon: coupon.BxGyDetails{
				BuyProducts: []coupon.CouponProduct{
					{ProductID: productXID, Quantity: 1},
					{ProductID: productYID, Quantity: 1},
				},
				GetProducts:     []coupon.CouponProduct{{ProductID: productAID, Quantity: 1}},
				RepetitionLimit: 1,
			},
			expectedCart: DiscountedCart{
				Items: []DiscountedItem{
					{ProductID: productXID, Quantity: 1, Price: 10, Discount: 0},
					{ProductID: productYID, Quantity: 1, Price: 12, Discount: 0},
					{ProductID: productAID, Quantity: 1, Price: 8, Discount: 8},
				},
				TotalPrice:    30,
				TotalDiscount: 8,
				FinalPrice:    22,
			},
		},
		{
			name: "B2G1 not applicable due to insufficient 'buy' items",
			items: []PricedItem{
				{ProductID: productXID, Quantity: 1, Price: 10},
				{ProductID: productAID, Quantity: 1, Price: 8},
				{ProductID: productBID, Quantity: 1, Price: 9},
			},
			totalPrice: 27,
			coupon: coupon.BxGyDetails{
				BuyProducts: []coupon.CouponProduct{
					{ProductID: productXID, Quantity: 1},
					{ProductID: productYID, Quantity: 1},
				},
				GetProducts:     []coupon.CouponProduct{{ProductID: productAID, Quantity: 1}},
				RepetitionLimit: 1,
			},
			expectedCart: DiscountedCart{
				Items: []DiscountedItem{
					{ProductID: productXID, Quantity: 1, Price: 10, Discount: 0},
					{ProductID: productAID, Quantity: 1, Price: 8, Discount: 0},
					{ProductID: productBID, Quantity: 1, Price: 9, Discount: 0},
				},
				TotalPrice:    27,
				TotalDiscount: 0,
				FinalPrice:    27,
			},
		},
		{
			name: "B2G1 applied multiple times up to repetition limit",
			items: []PricedItem{
				{ProductID: productXID, Quantity: 6, Price: 10},
				{ProductID: productAID, Quantity: 1, Price: 8},
				{ProductID: productBID, Quantity: 1, Price: 9},
				{ProductID: productCID, Quantity: 1, Price: 11},
			},
			totalPrice: 88,
			coupon: coupon.BxGyDetails{
				BuyProducts: []coupon.CouponProduct{{ProductID: productXID, Quantity: 2}},
				GetProducts: []coupon.CouponProduct{
					{ProductID: productAID, Quantity: 1},
					{ProductID: productBID, Quantity: 1},
					{ProductID: productCID, Quantity: 1},
				},
				RepetitionLimit: 3,
			},
			expectedCart: DiscountedCart{
				Items: []DiscountedItem{
					{ProductID: productXID, Quantity: 6, Price: 10, Discount: 0},
					{ProductID: productAID, Quantity: 1, Price: 8, Discount: 8},
					{ProductID: productBID, Quantity: 1, Price: 9, Discount: 9},
					{ProductID: productCID, Quantity: 1, Price: 11, Discount: 11},
				},
				TotalPrice:    88,
				TotalDiscount: 28,
				FinalPrice:    60,
			},
		},
		{
			name: "Not enough 'get' items to apply coupon multiple times",
			items: []PricedItem{
				{ProductID: productXID, Quantity: 6, Price: 10},
				{ProductID: productAID, Quantity: 1, Price: 8},
			},
			totalPrice: 68,
			coupon: coupon.BxGyDetails{
				BuyProducts:     []coupon.CouponProduct{{ProductID: productXID, Quantity: 2}},
				GetProducts:     []coupon.CouponProduct{{ProductID: productAID, Quantity: 1}},
				RepetitionLimit: 3,
			},
			expectedCart: DiscountedCart{
				Items: []DiscountedItem{
					{ProductID: productXID, Quantity: 6, Price: 10, Discount: 0},
					{ProductID: productAID, Quantity: 1, Price: 8, Discount: 8},
				},
				TotalPrice:    68,
				TotalDiscount: 8,
				FinalPrice:    60,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			coup := coupon.Coupon{
				ID:      1,
				Type:    "bxgy",
				Details: tc.coupon,
			}
			gotCart := applyBxGyWiseCoupon(tc.items, tc.totalPrice, coup)
			if !reflect.DeepEqual(gotCart, tc.expectedCart) {
				t.Errorf("applyBxGyWiseCoupon() = %+v, want %+v", gotCart, tc.expectedCart)
			}
		})
	}
}

func TestApplyProductWiseCoupon(t *testing.T) {
	const (
		productAID = 1
		productBID = 2
		productCID = 3
	)

	tests := []struct {
		name         string
		items        []PricedItem
		totalPrice   int
		coupon       coupon.ProductWiseDetails
		expectedCart DiscountedCart
	}{
		{
			name: "Scenario: 20% off on a product",
			items: []PricedItem{
				{ProductID: productAID, Quantity: 1, Price: 100},
				{ProductID: productBID, Quantity: 1, Price: 50},
			},
			totalPrice: 150,
			coupon: coupon.ProductWiseDetails{
				ProductID: productAID, Discount: 20,
			},
			expectedCart: DiscountedCart{
				Items: []DiscountedItem{
					{ProductID: productAID, Quantity: 1, Price: 100, Discount: 20},
					{ProductID: productBID, Quantity: 1, Price: 50, Discount: 0},
				},
				TotalPrice:    150,
				TotalDiscount: 20,
				FinalPrice:    130,
			},
		},
		{
			name: "Scenario: No matching product in cart",
			items: []PricedItem{
				{ProductID: productBID, Quantity: 1, Price: 50},
				{ProductID: productCID, Quantity: 1, Price: 75},
			},
			totalPrice: 125,
			coupon: coupon.ProductWiseDetails{
				ProductID: productAID, Discount: 20,
			},
			expectedCart: DiscountedCart{
				Items: []DiscountedItem{
					{ProductID: productBID, Quantity: 1, Price: 50, Discount: 0},
					{ProductID: productCID, Quantity: 1, Price: 75, Discount: 0},
				},
				TotalPrice:    125,
				TotalDiscount: 0,
				FinalPrice:    125,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			coup := coupon.Coupon{
				ID:      1,
				Type:    "product-wise",
				Details: tc.coupon,
			}
			gotCart := applyProductWiseCoupon(tc.items, tc.totalPrice, coup)
			if !reflect.DeepEqual(gotCart, tc.expectedCart) {
				t.Errorf("applyProductWiseCoupon() = %v, want %v", gotCart, tc.expectedCart)
			}
		})
	}
}

func TestApplyCartWiseCoupon(t *testing.T) {
	const (
		productAID = 1
		productBID = 2
		productCID = 3
	)

	tests := []struct {
		name         string
		items        []PricedItem
		totalPrice   int
		coupon       coupon.CartWiseDetails
		expectedCart DiscountedCart
	}{
		{
			name: "Total above threshold, apply 10% discount",
			items: []PricedItem{
				{ProductID: productAID, Quantity: 1, Price: 200},
				{ProductID: productBID, Quantity: 1, Price: 100},
			},
			totalPrice: 300,
			coupon: coupon.CartWiseDetails{
				Threshold: 250, Discount: 10,
			},
			expectedCart: DiscountedCart{
				Items: []DiscountedItem{
					{ProductID: productAID, Quantity: 1, Price: 200, Discount: 0},
					{ProductID: productBID, Quantity: 1, Price: 100, Discount: 0},
				},
				TotalPrice:    300,
				TotalDiscount: 30,
				FinalPrice:    270,
			},
		},
		{
			name: "Total below threshold, no discount applied",
			items: []PricedItem{
				{ProductID: productAID, Quantity: 1, Price: 80},
				{ProductID: productCID, Quantity: 1, Price: 100},
			},
			totalPrice: 180,
			coupon: coupon.CartWiseDetails{
				Threshold: 200, Discount: 10,
			},
			expectedCart: DiscountedCart{
				Items: []DiscountedItem{
					{ProductID: productAID, Quantity: 1, Price: 80, Discount: 0},
					{ProductID: productCID, Quantity: 1, Price: 100, Discount: 0},
				},
				TotalPrice:    180,
				TotalDiscount: 0,
				FinalPrice:    180,
			},
		},
		{
			name: "Total equals threshold, discount applies",
			items: []PricedItem{
				{ProductID: productBID, Quantity: 2, Price: 100},
			},
			totalPrice: 200,
			coupon: coupon.CartWiseDetails{
				Threshold: 200, Discount: 25,
			},
			expectedCart: DiscountedCart{
				Items: []DiscountedItem{
					{ProductID: productBID, Quantity: 2, Price: 100, Discount: 0},
				},
				TotalPrice:    200,
				TotalDiscount: 50,
				FinalPrice:    150,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			coup := coupon.Coupon{
				ID:      1,
				Type:    "cart-wise",
				Details: tc.coupon,
			}
			gotCart := applyCartWiseCoupon(tc.items, tc.totalPrice, coup)
			if !reflect.DeepEqual(gotCart, tc.expectedCart) {
				t.Errorf("applyCartWiseCoupon() = %v, want %v", gotCart, tc.expectedCart)
			}
		})
	}
}
