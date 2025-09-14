package cart

import (
	"fmt"

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

	// TODO: Create actual couponID in the coupon and modify the coupon repository
	// code accordingly, the current code can't handle the delete, as deleted entries will result
	// in old coupon ids if we rely on the index only
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
			if discount, ok := appliableBxGYCoupon(items, coupon); ok {
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

func appliableCartWiseCoupons(totalPrice int, coup coupon.Coupon) (int, bool) {
	detail := coup.Details.(coupon.CartWiseDetails)
	if detail.Threshold > totalPrice {
		return 0, false
	}
	return (detail.Discount * totalPrice) / 100, true
}

func appliableProductWiseCoupon(items []PricedItem, coup coupon.Coupon) (int, bool) {
	detail := coup.Details.(coupon.ProductWiseDetails)
	for _, item := range items {
		if item.ProductID == detail.ProductID {
			return (detail.Discount * item.Price * item.Quantity) / 100, true
		}
	}
	return 0, false
}

func appliableBxGYCoupon(items []PricedItem, coup coupon.Coupon) (int, bool) {
	detail := coup.Details.(coupon.BxGyDetails)

	cartMap := map[int]PricedItem{} // map of productID -> PricedItem
	for _, product := range items {
		cartMap[product.ProductID] = product
	}

	totalBuyRequired := detail.BuyProducts[0].Quantity
	totalBuyInCart := 0
	for _, product := range detail.BuyProducts {
		productInCart, ok := cartMap[product.ProductID]
		if !ok {
			continue
		}
		totalBuyInCart += productInCart.Quantity
	}

	if totalBuyRequired == 0 {
		return 0, false
	}

	actualRepetitions := min(totalBuyInCart/totalBuyRequired, detail.RepetitionLimit)
	if actualRepetitions == 0 {
		return 0, false
	}

	totalDiscount := 0
	for _, product := range detail.GetProducts {
		productInCart, ok := cartMap[product.ProductID]
		if !ok {
			continue
		}
		totalDiscount += product.Quantity * productInCart.Price * actualRepetitions
	}
	if totalDiscount == 0 {
		return 0, false
	}

	return totalDiscount, true
}

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
// along with the total discount
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
	detail := coup.Details.(coupon.BxGyDetails)

	discount, ok := appliableBxGYCoupon(items, coup)
	if !ok {
		discount = 0
	}

	getProductsMap := map[int]struct{}{}
	for _, prod := range detail.GetProducts {
		getProductsMap[prod.ProductID] = struct{}{}
	}

	for i, item := range items {
		_, inGetProduct := getProductsMap[item.ProductID]
		if !inGetProduct {
			discountedItems[i] = item.ToDiscountedItem(0)
			continue
		}
		discountedItems[i] = item.ToDiscountedItem(discount)
	}
	return DiscountedCart{
		Items:         discountedItems,
		TotalPrice:    totalPrice,
		TotalDiscount: discount,
		FinalPrice:    totalPrice - discount,
	}
}
