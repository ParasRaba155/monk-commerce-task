package main

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ParasRaba155/monk-commerce-task/cart"
	"github.com/ParasRaba155/monk-commerce-task/coupon"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	repo := coupon.NewRepository()
	couponHandler := coupon.NewHandler(repo)
	cartHandler := cart.NewHandler(repo)

	e.POST("/coupons", couponHandler.Create)
	e.GET("/coupons", couponHandler.Get)
	e.GET("/coupons/:id", couponHandler.GetByID)
	e.PUT("/coupons/:id", couponHandler.UpdateByID)
	e.DELETE("/coupons/:id", couponHandler.DeleteByID)

	e.POST("/applicable-coupon", cartHandler.ApplicableCoupon)
	e.POST("/apply-coupon/:id", cartHandler.ApplyCoupon)

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
