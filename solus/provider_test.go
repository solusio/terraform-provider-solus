package solus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	t.Run("runs internal validation without error", func(t *testing.T) {
		assert.NoError(t, Provider().InternalValidate())
	})
}
