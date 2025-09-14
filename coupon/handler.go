package coupon

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ParasRaba155/monk-commerce-task/utils"
)

type Repository interface {
	CreateCoupon(coupon Coupon) error
	GetAllCoupons() ([]Coupon, error)
	GetCouponByID(id int) (Coupon, error)
	UpdateCouponByID(id int, newCoupon Coupon) (Coupon, error)
	DeleteCouponByID(id int) error
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
	id, err := utils.ParamIDHelper(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
	}

	coupon, err := h.Repo.GetCouponByID(id)
	if err != nil {
		slog.Error("get coupon by id db", slog.Any("err", err), slog.Int("id", id))
		if errors.Is(err, ErrDoesNotExist) {
			return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
		}
		return c.JSON(http.StatusInternalServerError, utils.GenericFailure(err))
	}
	return c.JSON(http.StatusOK, utils.GenericSuccess(coupon))
}

func (h Handler) UpdateByID(c echo.Context) error {
	id, err := utils.ParamIDHelper(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
	}

	// ideally we might have different request body for update and create
	// but for simplicity we will have the same request body and update it as a whole
	// after validation
	var req CreateCouponReq
	if err := c.Bind(&req); err != nil {
		slog.Error("update coupon bind error", slog.Any("err", err))
		return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
	}

	if err := req.Validate(); err != nil {
		slog.Error("update coupon validate error", slog.Any("err", err))
		return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
	}

	updated, err := h.Repo.UpdateCouponByID(id, Coupon{
		Type:    CouponType(req.Type),
		Details: req.Details,
	})
	if err != nil {
		slog.Error("update coupon by id db", slog.Any("err", err), slog.Int("id", id))
		if errors.Is(err, ErrDoesNotExist) {
			return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
		}
		return c.JSON(http.StatusInternalServerError, utils.GenericFailure(err))
	}

	return c.JSON(http.StatusOK, utils.GenericSuccess(updated))
}

func (h Handler) DeleteByID(c echo.Context) error {
	id, err := utils.ParamIDHelper(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
	}
	err = h.Repo.DeleteCouponByID(id)

	if err != nil {
		slog.Error("delete coupon by id db", slog.Any("err", err), slog.Int("id", id))
		if errors.Is(err, ErrDoesNotExist) {
			return c.JSON(http.StatusBadRequest, utils.GenericFailure(err))
		}
		return c.JSON(http.StatusInternalServerError, utils.GenericFailure(err))
	}
	return c.JSON(http.StatusNoContent, nil)
}
