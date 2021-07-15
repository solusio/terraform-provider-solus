package solus

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func diagErr(summary string, detailsFormat string, args ...interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  summary,
			Detail:   fmt.Sprintf(detailsFormat, args...),
		},
	}
}

type schemaChainSetter struct {
	d   *schema.ResourceData
	err error
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
