// Package utils for any additional utilities
package utils

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
	return genericResponse{
		Success: false,
		Error:   m,
	}
}
