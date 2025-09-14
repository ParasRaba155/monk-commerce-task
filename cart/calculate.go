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
			if discount, ok := applicableCartWiseCoupons(totalPrice, coupon); ok {
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

func applicableCartWiseCoupons(totalPrice int, coup coupon.Coupon) (int, bool) {
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
