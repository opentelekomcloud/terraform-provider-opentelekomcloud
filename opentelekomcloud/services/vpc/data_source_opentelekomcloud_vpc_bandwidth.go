package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/bandwidths"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceBandWidth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBandWidthRead,

		DeprecationMessage: "please use `opentelekomcloud_vpc_bandwidth_v2` data source instead",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(5, 2000),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourceBandWidthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud vpc client: %s", err)
	}

	listOpts := bandwidths.ListOpts{
		ShareType: "WHOLE",
	}

	allBWs, err := bandwidths.List(vpcClient, listOpts).Extract()
	if err != nil {
		return fmterr.Errorf("unable to list OpenTelekomCloud bandwidths: %s", err)
	}
	if len(allBWs) == 0 {
		return fmterr.Errorf("no OpenTelekomCloud bandwidth was found")
	}

	// Filter bandwidths by "name"
	var bandList []bandwidths.BandWidth
	name := d.Get("name").(string)
	for _, band := range allBWs {
		if name == band.Name {
			bandList = append(bandList, band)
		}
	}
	if len(bandList) == 0 {
		return fmterr.Errorf("no OpenTelekomCloud bandwidth was found by name: %s", name)
	}

	// Filter bandwidths by "size"
	result := bandList[0]
	if v, ok := d.GetOk("size"); ok {
		var found bool
		for _, band := range bandList {
			if v.(int) == band.Size {
				found = true
				result = band
				break
			}
		}
		if !found {
			return fmterr.Errorf("no OpenTelekomCloud bandwidth was found by size: %d", v.(int))
		}
	}

	log.Printf("[DEBUG] Retrieved OpenTelekomCloud bandwidth %s: %+v", result.ID, result)
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
