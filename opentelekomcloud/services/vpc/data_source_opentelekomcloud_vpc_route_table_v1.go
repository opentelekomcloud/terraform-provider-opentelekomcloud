package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/routetables"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceVPCRouteTableV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcRouteTableV1Read,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcRouteTableV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := routetables.ListOpts{
		VpcID: d.Get("vpc_id").(string),
		ID:    d.Get("id").(string),
	}
	allRouteTables, err := routetables.List(client, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve OpenTelekomCloud VPC route tables: %s", err)
	}

	if len(allRouteTables) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var rtbID string
	if v, ok := d.GetOk("name"); ok {
		filterName := v.(string)
		for _, rtb := range allRouteTables {
			if filterName == rtb.Name {
				rtbID = rtb.ID
				break
			}
		}
	} else {
		// find the default route table if name was not specified
		for _, rtb := range allRouteTables {
			if rtb.Default {
				rtbID = rtb.ID
				break
			}
		}
	}

	if rtbID == "" {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	// call Get API to retrieve more details about the route table
	routeTable, err := routetables.Get(client, rtbID)
	if err != nil {
		return diag.Errorf("unable to retrieve OpenTelekomCloud route VPC table %s: %s", rtbID, err)
	}

	log.Printf("[DEBUG] Retrieved VPC route table %s: %+v", rtbID, routeTable)
	d.SetId(rtbID)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("vpc_id", routeTable.VpcID),
		d.Set("name", routeTable.Name),
		d.Set("description", routeTable.Description),
		d.Set("default", routeTable.Default),
		d.Set("subnets", expandRouteTableSubnets(routeTable.Subnets)),
		d.Set("route", expandRouteTableRoutes(routeTable.Routes)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud VPC route table: %s", err)
	}

	return nil
}
