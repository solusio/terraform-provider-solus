package solus

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/guregu/null.v4"
)

type schemaChainSetter struct {
	d   *schema.ResourceData
	err error
}

func (s *schemaChainSetter) SetID(v int) *schemaChainSetter {
	if s.err == nil {
		s.d.SetId(strconv.Itoa(v))
	}
	return s
}

func (s *schemaChainSetter) Set(k string, v interface{}) *schemaChainSetter {
	if s.err != nil {
		return s
	}

	if err := s.d.Set(k, v); err != nil {
		s.err = fmt.Errorf("failed to set value for %q key: %w", k, err)
	}
	return s
}

func (s *schemaChainSetter) Error() error {
	return s.err
}

func newNullableIntForID(i int) null.Int {
	// Because valid ID can't be null.
	return null.NewInt(int64(i), i != 0)
}
