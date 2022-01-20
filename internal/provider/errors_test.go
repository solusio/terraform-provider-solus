// Copyright 1999-2022. Plesk International GmbH. All rights reserved.

package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/solusio/solus-go-sdk"
	"github.com/stretchr/testify/assert"
)

func Test_normalizeAPIError(t *testing.T) {
	cc := map[string]struct {
		given    error
		expected error
	}{
		"nil": {},

		"not found": {
			solus.HTTPError{HTTPCode: http.StatusNotFound},
			errResourceNotFound,
		},

		"wrapped not found": {
			fmt.Errorf("fake error: %w", solus.HTTPError{HTTPCode: http.StatusNotFound}),
			errResourceNotFound,
		},

		"bad request": {
			solus.HTTPError{
				HTTPCode: http.StatusUnprocessableEntity,
				Errors: map[string][]string{
					"foo": {"bar"},
				},
			},
			errBadRequest,
		},
	}

	for name, c := range cc {
		t.Run(name, func(t *testing.T) {
			actual := normalizeAPIError(c.given)
			assert.ErrorIs(t, actual, c.expected)
		})
	}

	t.Run("over errors", func(t *testing.T) {
		err := solus.HTTPError{}

		actual := normalizeAPIError(err)
		assert.Equal(t, err, actual)
	})
}

func Test_marshalErrors(t *testing.T) {
	cc := map[string]map[string][]string{
		"": nil,
		"foo: bar": {
			"foo": {"bar"},
		},
		"foo: [fizz buzz]": {
			"foo": {"fizz", "buzz"},
		},
		"foo: 1, bar: 2": {
			"foo": {"1"},
			"bar": {"2"},
		},
	}

	for expected, given := range cc {
		t.Run(expected, func(t *testing.T) {
			actual := marshalErrors(given)
			assert.Equal(t, expected, actual)
		})
	}
}
