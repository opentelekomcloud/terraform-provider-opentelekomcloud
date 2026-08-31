package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/vpcs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVirtualPrivateCloudVpcV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVirtualPrivateCloudV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"routes": vpcRoutesSchema(),
			"tenant_id": {
				Type:     schema.TypeString,
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
		},
	}
}

func dataSourceVirtualPrivateCloudV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	listOpts := vpcs.ListOpts{
		ID:                  d.Get("id").(string),
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
	}

	allVpcs, err := vpcs.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve VPCs: %w", err)
	}
	filters := vpcFilters{
		Name:   d.Get("name").(string),
		CIDR:   d.Get("cidr").(string),
		Status: d.Get("status").(string),
	}
	shared := d.GetRawConfig().GetAttr("shared")
	if shared.IsKnown() && !shared.IsNull() {
		value := d.Get("shared").(bool)
		filters.Shared = &value
	}
	refinedVPCs := filterVpcs(allVpcs, filters)

	if len(refinedVPCs) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedVPCs) > 1 {
		return fmterr.Errorf("your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	singleVpc := refinedVPCs[0]

	log.Printf("[INFO] Retrieved Vpc using given filter %s: %+v", singleVpc.ID, singleVpc)
	d.SetId(singleVpc.ID)

	mErr := multierror.Append(
		d.Set("name", singleVpc.Name),
		d.Set("description", singleVpc.Description),
		d.Set("cidr", singleVpc.CIDR),
		d.Set("status", singleVpc.Status),
		d.Set("shared", singleVpc.EnableSharedSnat),
		d.Set("enterprise_project_id", singleVpc.EnterpriseProjectID),
		d.Set("tenant_id", singleVpc.TenantId),
		d.Set("created_at", singleVpc.CreatedAt),
		d.Set("updated_at", singleVpc.UpdatedAt),
		d.Set("region", config.GetRegion(d)),
		d.Set("routes", flattenVpcRoutes(singleVpc.Routes)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

type vpcFilters struct {
	Name   string
	CIDR   string
	Status string
	Shared *bool
}

func filterVpcs(allVpcs []vpcs.Vpc, filters vpcFilters) []vpcs.Vpc {
	refined := make([]vpcs.Vpc, 0, len(allVpcs))
	for _, vpc := range allVpcs {
		if filters.Name != "" && vpc.Name != filters.Name {
			continue
		}
		if filters.CIDR != "" && vpc.CIDR != filters.CIDR {
			continue
		}
		if filters.Status != "" && vpc.Status != filters.Status {
			continue
		}
		if filters.Shared != nil && vpc.EnableSharedSnat != *filters.Shared {
			continue
		}
		refined = append(refined, vpc)
	}
	return refined
}

func vpcRoutesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"destination": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"nexthop": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenVpcRoutes(routes []vpcs.Route) []map[string]interface{} {
	result := make([]map[string]interface{}, len(routes))
	for i, route := range routes {
		result[i] = map[string]interface{}{
			"destination": route.DestinationCIDR,
			"nexthop":     route.NextHop,
		}
	}
	return result
}
