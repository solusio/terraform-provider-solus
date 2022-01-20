// Copyright 1999-2022. Plesk International GmbH. All rights reserved.

package provider

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/solusio/solus-go-sdk"
)

var (
	errResourceNotFound = errors.New("not found")
	errTooManyResults   = errors.New("too many results")
	errBadRequest       = errors.New("bad request")
)

func normalizeAPIError(err error) error {
	if err == nil {
		return nil
	}

	var httpErr solus.HTTPError
	if !errors.As(err, &httpErr) {
		return err
	}

	switch httpErr.HTTPCode {
	case http.StatusNotFound:
		return errResourceNotFound

	case http.StatusUnprocessableEntity:
		return fmt.Errorf("%w: %v", errBadRequest, marshalErrors(httpErr.Errors))
	}

	return err
}

func marshalErrors(errs map[string][]string) string {
	buf := bytes.Buffer{}

	for k, v := range errs {
		if len(v) == 1 {
			buf.WriteString(fmt.Sprintf("%s: %s", k, v[0]))
		} else {
			buf.WriteString(fmt.Sprintf("%s: %v", k, v))
		}
		buf.WriteString(", ")
	}

	return strings.TrimSuffix(buf.String(), ", ")
}
