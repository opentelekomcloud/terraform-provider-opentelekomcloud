package fmterr

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Errorf wraps fmt.Errorf into diag.Diagnostics
func Errorf(format string, a ...interface{}) diag.Diagnostics {
	return diag.FromErr(
		fmt.Errorf(format, a...),
	)
}
