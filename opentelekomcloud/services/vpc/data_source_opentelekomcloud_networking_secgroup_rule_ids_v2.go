package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/groups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceNetworkingSecGroupRuleIdsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupIdsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceNetworkingSecGroupIdsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	secGroupID := d.Get("security_group_id").(string)
	listOpts := groups.ListOpts{
		ID: secGroupID,
	}

	secGroupPages, err := groups.List(client, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to retrieve security group pages: %w", err)
	}
	secGroupList, err := groups.ExtractGroups(secGroupPages)
	if err != nil {
		return fmterr.Errorf("unable to extract security groups: %w", err)
	}

	if len(secGroupList) == 0 {
		return fmterr.Errorf("no matching security groups found for security group with ID: %s", secGroupID)
	}

	if len(secGroupList) > 1 {
		return fmterr.Errorf("more than 1 security groups found for security group with ID: %s", secGroupID)
	}

	foundSecGroup := secGroupList[0]

	secGroupRules := make([]string, len(foundSecGroup.Rules))
	for _, rule := range foundSecGroup.Rules {
		secGroupRules = append(secGroupRules, rule.ID)
	}

	d.SetId(secGroupID)
	mErr := multierror.Append(
		d.Set("ids", secGroupRules),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
