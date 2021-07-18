package solus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newNullableIntForID(t *testing.T) {
	t.Run("not null", func(t *testing.T) {
		i := newNullableIntForID(42)
		assert.Equal(t, int64(42), i.Int64)
		assert.True(t, i.Valid)
	})

	t.Run("null", func(t *testing.T) {
		i := newNullableIntForID(0)
		assert.Equal(t, int64(0), i.Int64)
		assert.False(t, i.Valid)
	})
}
