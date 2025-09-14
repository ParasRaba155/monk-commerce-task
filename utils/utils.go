// Package utils for any additional utilities
package utils

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
)

var isAlphaNumericRegex = regexp.MustCompile(`^\d+$`)

type genericResponse struct {
	Success bool           `json:"success"`
	Data    any            `json:"data,omitempty"`
	Error   map[string]any `json:"error,omitempty"`
}

func GenericSuccess(data any) genericResponse {
	return genericResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	}
}

func GenericFailure(err any) genericResponse {
	errErr, ok := err.(error)
	m := map[string]any{}
	if ok {
		m["message"] = errErr.Error()
		return genericResponse{
			Success: false,
			Data:    nil,
			Error:   m,
		}
	}
	m["message"] = err
	return genericResponse{
		Success: false,
		Error:   m,
	}
}

// IsNonNegativeAlphaNumeric validates against the regex `^[\d+]$`
func IsNonNegativeAlphaNumeric(str string) bool {
	return isAlphaNumericRegex.MatchString(str)
}

// ParamIDHelper will check the param id, and make sure that it's a non-negative
// alphanumeric
func ParamIDHelper(c echo.Context) (int, error) {
	idstr := c.Param("id")
	if !IsNonNegativeAlphaNumeric(idstr) {
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
