package opentelekomcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// The result may be not correct when the type of param is string and user config it to 'param=""'
// but, there is no other way.
func hasFilledOpt(d *schema.ResourceData, param string) bool {
	_, b := d.GetOkExists(param)
	return b
}
