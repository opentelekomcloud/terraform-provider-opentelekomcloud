package rms

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rms/region"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRmsRegionsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRmsRegionsRead,

		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRmsRegionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, rmsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.RmsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationRMSV1Client, err)
	}

	regions, err := region.GetRegions(client, GetRmsDomainId(client, config))

	if err != nil {
		return diag.Errorf("error getting the regions list from server: %s", err)
	}

	var flattenRegions []map[string]interface{}

	for _, item := range regions {
		query := map[string]interface{}{
			"region_id":    item.RegionId,
			"display_name": item.DisplayName,
		}
		flattenRegions = append(flattenRegions, query)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("regions", flattenRegions),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
