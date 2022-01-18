// Copyright 1999-2022. Plesk International GmbH. All rights reserved.

package provider

import (
	"errors"
	"fmt"
	"net/http"

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
		return fmt.Errorf("%w: %v", errBadRequest, httpErr.Errors)
	}

	return err
}
