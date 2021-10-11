package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
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
		RouteID:     d.Get("id").(string), // d.Id() will return an empty string
	}

	pages, err := routes.List(vpcRouteClient, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to list vpc Routes: %s", err)
	}
	refinedRoutes, err := routes.ExtractRoutes(pages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve vpc routes: %s", err)
	}

	if len(refinedRoutes) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedRoutes) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	route := refinedRoutes[0]

	log.Printf("[INFO] Retrieved Vpc Route using given filter %s: %+v", route.RouteID, route)
	d.SetId(route.RouteID)

	mErr := multierror.Append(
		d.Set("type", route.Type),
		d.Set("nexthop", route.NextHop),
		d.Set("destination", route.Destination),
		d.Set("tenant_id", route.Tenant_Id),
		d.Set("vpc_id", route.VPC_ID),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting VPC route attributes: %w", err)
	}

	return nil
}
