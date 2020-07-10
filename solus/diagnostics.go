package solus

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
