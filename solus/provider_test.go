package solus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvider(t *testing.T) {
	t.Run("runs internal validation without error", func(t *testing.T) {
		assert.NoError(t, Provider().InternalValidate())
	})
}
