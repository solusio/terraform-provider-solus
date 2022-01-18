package solus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// HTTPError represents errors occurred when some action failed due to some problem
// with request.
type HTTPError struct {
	Method   string
	Path     string
	HTTPCode int                 `json:"http_code"`
	Message  string              `json:"message"`
	Errors   map[string][]string `json:"errors"`
}

func (e HTTPError) Error() string {
	buf := bytes.NewBufferString(fmt.Sprintf("HTTP %s %s returns %d status code", e.Method, e.Path, e.HTTPCode))

	if len(e.Errors) > 0 {
		buf.WriteString(" with errors")
	}

	if e.Message != "" {
		//goland:noinspection GrazieInspection
		buf.WriteString(fmt.Sprintf(": %s", e.Message))
	}

	return buf.String()
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

// IsNotFound returns true if specified error is produced 'cause requested resource
// is not found.
func IsNotFound(err error) bool {
	var httpErr HTTPError
	if !errors.As(err, &httpErr) {
		return false
	}

	return httpErr.HTTPCode == http.StatusNotFound
}
