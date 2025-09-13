package coupon

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ParasRaba155/monk-commerce-task/utils"
)

func paramIDHelper(c echo.Context) (int, error) {
	idstr := c.Param("id")
	if !utils.IsNonNegativeAlphaNumeric(idstr) {
		slog.Error("get coupon by id validation", slog.String("err", "id must be non negative number"), slog.String("idstr", idstr))
		return 0, fmt.Errorf("id must be non negative number")
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		slog.Error("get coupon by id parsing", slog.Any("err", err), slog.String("idstr", idstr))
		return 0, fmt.Errorf("invalid id: %w", err)
	}
	return int(id), nil
}
