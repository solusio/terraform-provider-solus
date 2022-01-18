package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_validationIsDomainName(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ww, ee := validationIsDomainName("example.com", "foo")
		assert.Nil(t, ww)
		assert.Nil(t, ee)
	})

	t.Run("negative", func(t *testing.T) {
		cc := map[string]interface{}{
			`expected type of "foo" to be string`: 42,
			"invalid domain name":                 "192.0.2.1",
		}

		for expected, val := range cc {
			t.Run(expected, func(t *testing.T) {
				ww, ee := validationIsDomainName(val, "foo")
				assert.Nil(t, ww)
				require.Len(t, ee, 1)
				assert.EqualError(t, ee[0], expected)
			})
		}
	})
}

func Test_isDomainName(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := []string{
			"example.com",
			"foo.example.com",
			"example",
			strings.Repeat("a", 63),
			fmt.Sprintf(
				"%[1]s.%[1]s.%[1]s.%[1]s.%[2]s",
				strings.Repeat("a", 63),
				strings.Repeat("a", 3),
			),
		}

		for _, c := range cc {
			t.Run(c, func(t *testing.T) {
				b := isDomainName(c)
				assert.True(t, b)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		cc := []string{
			"192.0.2.1",
			"example.com/foo",
			"https://example.com",
			"example$com",
			strings.Repeat("a", 64),
			fmt.Sprintf(
				"%[1]s.%[1]s.%[1]s.%[1]s.%[2]s",
				strings.Repeat("a", 63),
				strings.Repeat("a", 4),
			),
		}

		for _, c := range cc {
			t.Run(c, func(t *testing.T) {
				b := isDomainName(c)
				assert.False(t, b)
			})
		}
	})
}

func Test_validationIsVirtualizationType(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ww, ee := validationIsVirtualizationType("kvm", "foo")
		assert.Nil(t, ww)
		assert.Nil(t, ee)
	})

	t.Run("negative", func(t *testing.T) {
		cc := map[string]interface{}{
			`expected type of "foo" to be string`: 42,
			`invalid virtualization type "bar"`:   "bar",
		}

		for expected, val := range cc {
			t.Run(expected, func(t *testing.T) {
				ww, ee := validationIsVirtualizationType(val, "foo")
				assert.Nil(t, ww)
				require.Len(t, ee, 1)
				assert.EqualError(t, ee[0], expected)
			})
		}
	})
}
