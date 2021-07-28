package solus

import (
	"errors"
	"io"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.Default().SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestSchemaChainSetter_SetID(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s := schemaChainSetter{d: (&schema.Resource{}).Data(nil)}
		s.SetID(10)

		assert.Equal(t, "10", s.d.Id())
	})

	t.Run("negative", func(t *testing.T) {
		s := schemaChainSetter{
			err: errors.New("fake error"),
			d:   (&schema.Resource{}).Data(nil),
		}
		s.SetID(10)

		assert.Equal(t, "", s.d.Id())
	})
}

func TestSchemaChainSetter_Set(t *testing.T) {
	res := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"foo": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fizz": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}

	t.Run("positive", func(t *testing.T) {
		s := schemaChainSetter{d: res.Data(nil)}

		err := s.
			Set("foo", "bar").
			Set("fizz", "buzz").
			Error()
		require.NoError(t, err)

		v := s.d.Get("foo")
		require.IsType(t, "", v)
		assert.Equal(t, "bar", v.(string))

		v = s.d.Get("fizz")
		require.IsType(t, "", v)
		assert.Equal(t, "buzz", v.(string))
	})

	t.Run("negative", func(t *testing.T) {
		if os.Getenv("TF_ACC") != "" {
			assert.PanicsWithError(t,
				"foo: '' expected type 'string', got unconvertible type 'int'",
				func() {
					(&schemaChainSetter{d: res.Data(nil)}).Set("foo", 42)
				},
			)
		} else {
			s := schemaChainSetter{d: res.Data(nil)}

			err := s.
				Set("foo", 42).
				Set("fizz", "buzz").
				Error()
			assert.EqualError(t, err, "failed to set value for \"foo\" key: foo: '' expected type 'string', got unconvertible type 'int'") //nolint:lll

			assert.Equal(t, "", s.d.Get("foo"))
			assert.Equal(t, "", s.d.Get("fizz"))
		}
	})
}

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
