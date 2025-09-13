package coupon

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ParasRaba155/monk-commerce-task/utils"
)

type Repository interface {
	CreateCoupon(coupon Coupon) error
	GetAllCoupons() ([]Coupon, error)
	GetCouponByID(id int) (Coupon, error)
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

func (h Handler) Get(c echo.Context) error {
	coupons, err := h.Repo.GetAllCoupons()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.GenericFailure(err))
	}
	return c.JSON(http.StatusOK, utils.GenericSuccess(coupons))
}

func (h Handler) GetByID(c echo.Context) error {
	idstr := c.Param("id")
	if !utils.IsNonNegativeAlphaNumeric(idstr) {
		slog.Error("get coupon by id validation", slog.String("err", "id must be non negative number"), slog.String("idstr", idstr))
		return c.JSON(http.StatusBadRequest, utils.GenericFailure("id must be non negative number"))
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		slog.Error("get coupon by id parsing", slog.Any("err", err), slog.String("idstr", idstr))
		return c.JSON(http.StatusInternalServerError, utils.GenericFailure(err))
	}

	coupon, err := h.Repo.GetCouponByID(int(id))
	if err != nil {
		slog.Error("get coupon by id db", slog.Any("err", err), slog.Int64("id", id))
		if errors.Is(err, errDoesNotExist) {
			return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
		}
		return c.JSON(http.StatusInternalServerError, utils.GenericFailure(err))
	}
	return c.JSON(http.StatusOK, utils.GenericSuccess(coupon))
}
