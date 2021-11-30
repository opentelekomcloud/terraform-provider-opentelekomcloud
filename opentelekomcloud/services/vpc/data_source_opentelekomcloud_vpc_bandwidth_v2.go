package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/bandwidths"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceBandWidthV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBandWidthV2Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBandWidthV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	pages, err := bandwidths.List(vpcClient).AllPages()
	if err != nil {
		return fmterr.Errorf("error listing bandwidths v2: %s", err)
	}
	bandSlice, err := bandwidths.ExtractBandwidths(pages)
	if err != nil {
		return fmterr.Errorf("error extracting bandwidth list: %w", err)
	}

	var results []bandwidths.Bandwidth
	expectedName, nameOk := d.GetOk("name")
	expectedSize, sizeOk := d.GetOk("size")
	for _, v := range bandSlice {
		if nameOk && v.Name != expectedName {
			continue
		}
		if sizeOk && v.Size != expectedSize {
			continue
		}
		results = append(results, v)
	}

	if len(results) < 1 {
		return common.DataSourceTooFewDiag
	}
	if len(results) > 1 {
		return common.DataSourceTooManyDiag
	}
	result := results[0]

	d.SetId(result.ID)
	mErr := multierror.Append(nil,
		d.Set("name", result.Name),
		d.Set("size", result.Size),
		d.Set("share_type", result.ShareType),
		d.Set("bandwidth_type", result.BandwidthType),
		d.Set("charge_mode", result.ChargeMode),
		d.Set("status", result.Status),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
