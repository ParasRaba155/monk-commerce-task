package cart

import (
	"fmt"
	"log/slog"

	"github.com/ParasRaba155/monk-commerce-task/coupon"
)

// TODO: remove logs

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
			slog.Info("cart-wise discount", slog.Int("total_price", totalPrice))
			if discount, ok := applicableCartWiseCoupons(totalPrice, coupon); ok {
				slog.Info("cart-wise discount", slog.Int("discount", discount))
				result = append(result, DiscountCoupon{
					CouponID: couponID,
					Type:     coupon.Type,
					Discount: discount,
				})
			}
		case "product-wise":
			slog.Info("product-wise discount", slog.Int("total_price", totalPrice))
			if discount, ok := appliableProductWiseCoupon(items, coupon); ok {
				slog.Info("product-wise discount", slog.Int("discount", discount))
				result = append(result, DiscountCoupon{
					CouponID: couponID,
					Type:     coupon.Type,
					Discount: discount,
				})
			}
		case "bxgy":
			slog.Info("bxgy-wise discount", slog.Int("total_price", totalPrice))
			if discount, ok := appliableBxGYCoupon(items, coupon); ok {
				slog.Info("bxgy discount", slog.Int("discount", discount))
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
	slog.Info("coupon detail cart wise", slog.Any("coupon", detail))
	if detail.Threshold > totalPrice {
		return 0, false
	}
	return (detail.Discount * totalPrice) / 100, true
}

func appliableProductWiseCoupon(items []PricedItem, coup coupon.Coupon) (int, bool) {
	detail := coup.Details.(coupon.ProductWiseDetails)
	slog.Info("coupon detail product wise", slog.Any("coupon", detail))
	for _, item := range items {
		if item.ProductID == detail.ProductID {
			return (detail.Discount * item.Price * item.Quantity) / 100, true
		}
	}
	return 0, false
}

func appliableBxGYCoupon(items []PricedItem, coup coupon.Coupon) (int, bool) {
	detail := coup.Details.(coupon.BxGyDetails)
	slog.Info("coupon detail bxgy", slog.Any("coupon", detail))

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
		prodcutInCart, ok := cartMap[product.ProductID]
		if !ok {
			return 0, false
		}
		totalDiscount += product.Quantity * prodcutInCart.Price * actualRepetitions
	}

	return totalDiscount, true
}
