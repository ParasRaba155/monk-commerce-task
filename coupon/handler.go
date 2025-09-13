package coupon

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ParasRaba155/monk-commerce-task/utils"
)

type Repository interface {
	CreateCoupon(coupon Coupon) error
}

type Handler struct {
	// Repo will give us a abstraction over db/repository layer
	// mostly the handler directly is not bulky and instead a additional service layer
	// is created to handle the business logic, however we will have bulky Handler methods for this case
	Repo Repository
}

func NewHandler(repo Repository) Handler {
	return Handler{
		Repo: repo,
	}
}

func (h Handler) Create(c echo.Context) error {
	var req CreateCouponReq
	if err := c.Bind(&req); err != nil {
		slog.Error("create coupon bind error", slog.Any("err", err))
		return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
	}

	if err := req.Validate(); err != nil {
		slog.Error("create coupon validate error", slog.Any("err", err))
		return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
	}

	h.Repo.CreateCoupon(Coupon{
		Type:    CouponType(req.Type),
		Details: req.Details,
	})
	return c.JSON(http.StatusCreated, utils.GenericSuccess("coupon created"))
}
