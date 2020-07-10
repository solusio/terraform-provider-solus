package solus

import (
	"encoding/json"
	"fmt"
)

type HTTPError struct {
	Method   string
	Path     string
	HTTPCode int    `json:"http_code"`
	Message  string `json:"message"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP %s %s returns %d status code: %s", e.Method, e.Path, e.HTTPCode, e.Message)
}

func newHTTPError(method, path string, httpCode int, body []byte) error {
	e := HTTPError{
		Method:   method,
		Path:     path,
		HTTPCode: httpCode,
	}

	if err := json.Unmarshal(body, &e); err != nil {
		e.Message = string(body)
		return e
	}

	return e
}
