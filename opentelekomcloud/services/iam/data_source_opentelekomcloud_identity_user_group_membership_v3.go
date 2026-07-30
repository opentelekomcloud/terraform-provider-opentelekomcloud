package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityUserGroupMembershipV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityUserGroupMembershipV3Read,

		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceIdentityUserGroupMembershipV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, iamClientKey, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	d.SetId(d.Get("user").(string))

	pages, err := users.ListGroups(client, d.Id()).AllPages()
	if err != nil {
		return fmterr.Errorf("error listing groups: %w", err)
	}
	gps, err := groups.ExtractGroups(pages)
	if err != nil {
		return fmterr.Errorf("error extracting group list: %w", err)
	}

	groupIDs := make([]string, len(gps))
	for i, g := range gps {
		groupIDs[i] = g.ID
	}
	if err := d.Set("groups", groupIDs); err != nil {
		return fmterr.Errorf("error setting group IDs: %w", err)
	}

	return nil
}
