package cart

import (
	"fmt"
	"slices"

	"github.com/ParasRaba155/monk-commerce-task/coupon"
)

// GetAppliableCoupons will return the list of all the applicable coupons
func GetAppliableCoupons(items []PricedItem, coupons []coupon.Coupon) []DiscountCoupon {
	if len(items) == 0 || len(coupons) == 0 {
		return nil
	}
	totalPrice := 0
	for _, item := range items {
		totalPrice += item.Price * item.Quantity
	}

	result := make([]DiscountCoupon, 0, len(coupons))

	for couponID, coupon := range coupons {
		switch coupon.Type {
		case "cart-wise":
			if discount, ok := appliableCartWiseCoupons(totalPrice, coupon); ok {
				result = append(result, DiscountCoupon{
					CouponID: couponID,
					Type:     coupon.Type,
					Discount: discount,
				})
			}
		case "product-wise":
			if discount, ok := appliableProductWiseCoupon(items, coupon); ok {
				result = append(result, DiscountCoupon{
					CouponID: couponID,
					Type:     coupon.Type,
					Discount: discount,
				})
			}
		case "bxgy":
			if discount, _, ok := appliableBxGYCoupon(items, coupon); ok {
				result = append(result, DiscountCoupon{
					CouponID: couponID,
					Type:     coupon.Type,
					Discount: discount,
				})
			}
		default:
			panic(fmt.Errorf("unsupported coupon type %s", coupon.Type))
		}
	}
	return result
}

// appliableCartWiseCoupons for handling the cart wise coupon
func appliableCartWiseCoupons(totalPrice int, coup coupon.Coupon) (int, bool) {
	detail := coup.Details.(coupon.CartWiseDetails)
	if detail.Threshold > totalPrice {
		return 0, false
	}
	return (detail.Discount * totalPrice) / 100, true
}

// appliableProductWiseCoupon for handling the product wise coupon
func appliableProductWiseCoupon(items []PricedItem, coup coupon.Coupon) (int, bool) {
	detail := coup.Details.(coupon.ProductWiseDetails)
	itemIdx := slices.IndexFunc(items, func(item PricedItem) bool {
		return item.ProductID == detail.ProductID
	})
	if itemIdx == -1 {
		return 0, false
	}
	item := items[itemIdx]
	return (detail.Discount * item.Price * item.Quantity) / 100, true
}

// appliableBxGYCoupon for handling the bxgy coupon
func appliableBxGYCoupon(items []PricedItem, coup coupon.Coupon) (int, map[int]int, bool) {
	detail := coup.Details.(coupon.BxGyDetails)

	cartMap := map[int]PricedItem{} // map of productID -> PricedItem
	for _, product := range items {
		cartMap[product.ProductID] = product
	}

	// we work on the assumption that each buy item has equal quantity
	// i.e. each item in the buy product will have same quantity
	totalBuyRequired := detail.BuyProducts[0].Quantity
	totalBuyInCart := 0
	for _, product := range detail.BuyProducts {
		// every product of the buy array should be in our cart for discount to be applicable
		productInCart, ok := cartMap[product.ProductID]
		if !ok {
			return 0, nil, false
		}
		totalBuyInCart += productInCart.Quantity
	}

	actualRepetitions := min(totalBuyInCart/totalBuyRequired, detail.RepetitionLimit)
	if actualRepetitions == 0 {
		return 0, nil, false
	}

	totalDiscount := 0
	productDiscounts := map[int]int{}
	for _, getProduct := range detail.GetProducts {
		productInCart, ok := cartMap[getProduct.ProductID]
		if !ok {
			continue
		}
		// how many times can this item by multiplied, say
		// we only have 1 item in the cart but repetation is 3 then we should only allow 1
		maxTimesByCart := productInCart.Quantity / getProduct.Quantity
		times := min(actualRepetitions, maxTimesByCart)

		discount := getProduct.Quantity * productInCart.Price * times
		productDiscounts[getProduct.ProductID] = discount
		totalDiscount += discount
	}
	if totalDiscount == 0 {
		return 0, nil, false
	}

	return totalDiscount, productDiscounts, true
}

// ApplyCoupon will apply the given coupon to the cart and return the discounted cart
// It will panic if the coupon is invalid
func ApplyCoupon(items []PricedItem, coupon coupon.Coupon) DiscountedCart {
	totalPrice := 0
	for _, item := range items {
		totalPrice += item.Price * item.Quantity
	}

	switch coupon.Type {
	case "cart-wise":
		return applyCartWiseCoupon(items, totalPrice, coupon)
	case "product-wise":
		return applyProductWiseCoupon(items, totalPrice, coupon)
	case "bxgy":
		return applyBxGyWiseCoupon(items, totalPrice, coupon)
	default:
		panic(fmt.Errorf("unsupported coupon type %s", coupon.Type))
	}
}

// applyCartWiseCoupon will apply the cart wise coupon
// since the coupon is on whole cart the discount on individual item will be zero
// and the total discount will be the calculated discount
func applyCartWiseCoupon(items []PricedItem, totalPrice int, coupon coupon.Coupon) DiscountedCart {
	discountedItems := make([]DiscountedItem, len(items))
	for i := range items {
		discountedItems[i] = items[i].ToDiscountedItem(0)
	}
	discount, ok := appliableCartWiseCoupons(totalPrice, coupon)
	if !ok {
		return DiscountedCart{
			Items:         discountedItems,
			TotalPrice:    totalPrice,
			TotalDiscount: 0,
			FinalPrice:    totalPrice,
		}
	}
	return DiscountedCart{
		Items:         discountedItems,
		TotalPrice:    totalPrice,
		TotalDiscount: discount,
		FinalPrice:    totalPrice - discount,
	}
}

// applyProductWiseCoupon will return the cart list with discount against the product
// along with the total discount and remaining products having zero discount
func applyProductWiseCoupon(items []PricedItem, totalPrice int, coup coupon.Coupon) DiscountedCart {
	discountedItems := make([]DiscountedItem, len(items))
	detail := coup.Details.(coupon.ProductWiseDetails)

	discount, ok := appliableProductWiseCoupon(items, coup)
	if !ok {
		discount = 0
	}

	for i, item := range items {
		if item.ProductID != detail.ProductID {
			discountedItems[i] = items[i].ToDiscountedItem(0)
			continue
		}
		discountedItems[i] = items[i].ToDiscountedItem(discount)
	}
	return DiscountedCart{
		Items:         discountedItems,
		TotalPrice:    totalPrice,
		TotalDiscount: discount,
		FinalPrice:    totalPrice - discount,
	}
}

// applyBxGyWiseCoupon will return the cart list with discount against the products
// in the get products from the bxgy along with the total discount
func applyBxGyWiseCoupon(items []PricedItem, totalPrice int, coup coupon.Coupon) DiscountedCart {
	discountedItems := make([]DiscountedItem, len(items))

	discount, productDiscounts, ok := appliableBxGYCoupon(items, coup)
	if !ok {
		discount = 0
	}

	// discount map has the associated discount against the get item
	// and default of zero for other items
	for i, item := range items {
		discount := productDiscounts[item.ProductID]
		discountedItems[i] = item.ToDiscountedItem(discount)
	}
	return DiscountedCart{
		Items:         discountedItems,
		TotalPrice:    totalPrice,
		TotalDiscount: discount,
		FinalPrice:    totalPrice - discount,
	}
}
