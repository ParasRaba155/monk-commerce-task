// Package utils for any additional utilities
package utils

import "regexp"

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
