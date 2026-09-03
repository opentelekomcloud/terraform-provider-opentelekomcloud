package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/bandwidths"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceBandWidthV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBandWidthV2Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
				Optional: true,
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
			"billing_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicip_info": bandwidthPublicIPInfoSchema(),
		},
	}
}

func dataSourceBandWidthV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcClient, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	enterpriseProjectID := config.GetEnterpriseProjectID(d)
	if enterpriseProjectID == "" {
		enterpriseProjectID = allGrantedEnterpriseProjects
	}
	bandSlice, err := bandwidths.List(vpcClient, bandwidths.ListOpts{EnterpriseProjectId: enterpriseProjectID})
	if err != nil {
		return fmterr.Errorf("error listing bandwidths: %w", err)
	}

	results := filterBandwidthsV1(bandSlice, bandwidthFilters{
		ID:                d.Get("id").(string),
		Name:              d.Get("name").(string),
		Size:              d.Get("size").(int),
		ShareType:         d.Get("share_type").(string),
		PublicBorderGroup: d.Get("public_border_group").(string),
	})

	if len(results) < 1 {
		return common.DataSourceTooFewDiag
	}
	if len(results) > 1 {
		return common.DataSourceTooManyDiag
	}
	result := results[0]

	d.SetId(result.ID)
	mErr := multierror.Append(nil,
		d.Set("id", result.ID),
		d.Set("name", result.Name),
		d.Set("size", result.Size),
		d.Set("share_type", result.ShareType),
		d.Set("bandwidth_type", result.BandwidthType),
		d.Set("charge_mode", result.ChargeMode),
		d.Set("status", result.Status),
		d.Set("billing_info", result.BillingInfo),
		d.Set("tenant_id", result.TenantId),
		d.Set("enterprise_project_id", result.EnterpriseProjectID),
		d.Set("public_border_group", result.PublicBorderGroup),
		d.Set("created_at", result.CreatedAt),
		d.Set("updated_at", result.UpdatedAt),
		d.Set("publicip_info", flattenBandwidthPublicIPs(result.PublicipInfo)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

type bandwidthFilters struct {
	ID                string
	Name              string
	Size              int
	ShareType         string
	PublicBorderGroup string
}

func filterBandwidthsV1(allBandwidths []bandwidths.BandWidth, filters bandwidthFilters) []bandwidths.BandWidth {
	result := make([]bandwidths.BandWidth, 0, len(allBandwidths))
	for _, bandwidth := range allBandwidths {
		if filters.ID != "" && bandwidth.ID != filters.ID {
			continue
		}
		if filters.Name != "" && bandwidth.Name != filters.Name {
			continue
		}
		if filters.Size != 0 && bandwidth.Size != filters.Size {
			continue
		}
		if filters.ShareType != "" && bandwidth.ShareType != filters.ShareType {
			continue
		}
		if filters.PublicBorderGroup != "" && bandwidth.PublicBorderGroup != filters.PublicBorderGroup {
			continue
		}
		result = append(result, bandwidth)
	}
	return result
}
