package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/routes"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVPCRouteIdsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcRouteIdsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
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

func dataSourceVpcRouteIdsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcRouteClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := routes.ListOpts{
		VPC_ID: d.Get("vpc_id").(string),
	}

	pages, err := routes.List(vpcRouteClient, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to list vpc Routes: %s", err)
	}

	refinedRoutes, err := routes.ExtractRoutes(pages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve vpc Routes: %s", err)
	}

	if len(refinedRoutes) == 0 {
		return fmterr.Errorf("no matching route found for vpc with id %s", d.Get("vpc_id").(string))
	}

	listRoutes := make([]string, 0)

	for _, route := range refinedRoutes {
		listRoutes = append(listRoutes, route.RouteID)
	}

	d.SetId(d.Get("vpc_id").(string))
	mErr := multierror.Append(
		d.Set("ids", listRoutes),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
