package vpc

import (
	"context"
	"log"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/routes"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVPCRouteV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcRouteV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nexthop": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVpcRouteV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcRouteClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := routes.ListOpts{
		Type:        d.Get("type").(string),
		Destination: d.Get("destination").(string),
		VPC_ID:      d.Get("vpc_id").(string),
		Tenant_Id:   d.Get("tenant_id").(string),
		RouteID:     d.Id(),
	}

	pages, err := routes.List(vpcRouteClient, listOpts).AllPages()
	refinedRoutes, err := routes.ExtractRoutes(pages)

	if err != nil {
		return fmterr.Errorf("Unable to retrieve vpc routes: %s", err)
	}

	if len(refinedRoutes) < 1 {
		return fmterr.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedRoutes) > 1 {
		return fmterr.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Route := refinedRoutes[0]

	log.Printf("[INFO] Retrieved Vpc Route using given filter %s: %+v", Route.RouteID, Route)
	d.SetId(Route.RouteID)

	d.Set("type", Route.Type)
	d.Set("nexthop", Route.NextHop)
	d.Set("destination", Route.Destination)
	d.Set("tenant_id", Route.Tenant_Id)
	d.Set("vpc_id", Route.VPC_ID)
	d.Set("id", Route.RouteID)
	d.Set("region", config.GetRegion(d))

	return nil
}
