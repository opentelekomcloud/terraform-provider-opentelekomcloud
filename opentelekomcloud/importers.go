package opentelekomcloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func importAsManaged(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	_ = d.Set("shared", false)
	return []*schema.ResourceData{d}, nil
}
