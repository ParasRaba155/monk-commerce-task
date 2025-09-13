package cart

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ParasRaba155/monk-commerce-task/coupon"
	"github.com/ParasRaba155/monk-commerce-task/utils"
)

type Repository interface {
	GetAllCoupons() ([]coupon.Coupon, error)
}

type cartHandler struct {
	Repo Repository
}

func NewHandler(repo Repository) cartHandler {
	return cartHandler{Repo: repo}
}

func (h cartHandler) ApplicableCoupon(c echo.Context) error {
	var req Cart
	if err := c.Bind(&req); err != nil {
		slog.Error("applicable coupon bind error", slog.Any("err", err))
		return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
	}

	if err := req.Validate(); err != nil {
		slog.Error("applicable coupon validate error", slog.Any("err", err))
		return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
	}

	pricedItems := make([]PricedItem, 0, len(req.Items))
	for _, item := range req.Items {
		price, err := getProductPrice(item.ProductID)
		if err != nil {
			slog.Error("applicable coupon get product price", slog.Any("err", err))
			return c.JSON(http.StatusInternalServerError, utils.GenericFailure(err))
		}
		pricedItems = append(pricedItems, item.ToPricedItem(price))
	}

	coupons, err := h.Repo.GetAllCoupons()
	if err != nil {
		slog.Error("applicable coupon get all coupons", slog.Any("err", err))
		return c.JSON(http.StatusInternalServerError, utils.GenericFailure(err))
	}

	response := GetAppliableCoupons(pricedItems, coupons)
	if len(response) == 0 {
		return c.JSON(http.StatusOK, utils.GenericSuccess("Sorry! No coupons are available for you"))
	}
	return c.JSON(http.StatusOK, utils.GenericSuccess(response))
}
