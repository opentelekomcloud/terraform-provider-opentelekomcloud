package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The result may be not correct when the type of param is string and user config it to 'param=""'
// but, there is no other way.
func HasFilledOpt(d *schema.ResourceData, param string) bool {
	_, b := d.GetOkExists(param) // nolint:staticcheck
	return b
}
