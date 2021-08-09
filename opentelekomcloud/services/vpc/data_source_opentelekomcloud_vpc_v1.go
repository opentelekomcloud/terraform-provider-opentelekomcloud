package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"routes": {
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
			},
		},
	}
}

func dataSourceVirtualPrivateCloudV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := vpcs.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
		CIDR:   d.Get("cidr").(string),
	}

	refinedVPCs, err := vpcs.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve VPCs: %w", err)
	}

	if len(refinedVPCs) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedVPCs) > 1 {
		return fmterr.Errorf("your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	singleVpc := refinedVPCs[0]

	var routes []map[string]interface{}
	for _, route := range singleVpc.Routes {
		mapping := map[string]interface{}{
			"destination": route.DestinationCIDR,
			"nexthop":     route.NextHop,
		}
		routes = append(routes, mapping)
	}

	log.Printf("[INFO] Retrieved Vpc using given filter %s: %+v", singleVpc.ID, singleVpc)
	d.SetId(singleVpc.ID)

	mErr := multierror.Append(
		d.Set("name", singleVpc.Name),
		d.Set("cidr", singleVpc.CIDR),
		d.Set("status", singleVpc.Status),
		d.Set("shared", singleVpc.EnableSharedSnat),
		d.Set("region", config.GetRegion(d)),
		d.Set("routes", routes),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
