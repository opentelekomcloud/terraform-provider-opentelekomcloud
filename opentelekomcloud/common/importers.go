package common

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ImportAsManaged(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	_ = d.Set("shared", false)
	return []*schema.ResourceData{d}, nil
}
