package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityGroupMembershipV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityGroupMembershipV3Read,

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"xuser_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"xuser_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityGroupMembershipV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}
	d.SetId(d.Get("group").(string))

	allPages, err := users.ListInGroup(identityClient, d.Id(), users.ListOpts{}).AllPages()
	if err != nil {
		if _, b := err.(golangsdk.ErrDefault404); b {
			return nil
		}
		return fmterr.Errorf("unable to query groups: %s", err)
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve users: %s", err)
	}

	var ul []map[string]interface{}
	for _, u := range allUsers {
		user := map[string]interface{}{
			"id":          u.ID,
			"name":        u.Name,
			"description": u.Description,
			"domain_id":   u.DomainID,
			"email":       u.Email,
			"xuser_type":  u.XUserType,
			"xuser_id":    u.XUserID,
		}
		ul = append(ul, user)
	}

	if err := d.Set("users", ul); err != nil {
		return fmterr.Errorf("error setting user list from IAM (%s), error: %s", d.Id(), err)
	}

	return nil
}
