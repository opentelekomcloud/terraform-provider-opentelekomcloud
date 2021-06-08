package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ImportAsManaged(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	_ = d.Set("shared", false)
	return []*schema.ResourceData{d}, nil
}
