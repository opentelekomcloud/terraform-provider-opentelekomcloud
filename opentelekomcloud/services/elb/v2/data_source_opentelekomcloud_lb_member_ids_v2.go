package v2

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/pools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceLBMemberIDsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBMemberIDsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"pool_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
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

func dataSourceLBMemberIDsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	poolID := d.Get("pool_id").(string)

	memberPages, err := pools.ListMembers(client, poolID, pools.ListMembersOpts{}).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to retrieve ELBv2 member pages: %w", err)
	}

	refinedMembers, err := pools.ExtractMembers(memberPages)
	if err != nil {
		return fmterr.Errorf("error extracting ELBv2 members: %w", err)
	}

	if len(refinedMembers) < 1 {
		log.Println("[WARN] Your query returned no results. Please change your search criteria and try again.")
	}

	var memberList []string
	for _, member := range refinedMembers {
		memberList = append(memberList, member.ID)
	}

	d.SetId(poolID)

	mErr := multierror.Append(
		d.Set("ids", memberList),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
