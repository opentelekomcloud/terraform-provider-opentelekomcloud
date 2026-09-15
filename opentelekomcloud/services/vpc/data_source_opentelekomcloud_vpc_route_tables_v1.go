package vpc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/routetables"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceVpcRouteTablesV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRouteTablesV1Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routetables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
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
						"subnets": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"routes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination": {
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
				},
			},
		},
	}
}

func dataSourceRouteTablesV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := routetables.ListOpts{
		ID:       d.Get("id").(string),
		VpcID:    d.Get("vpc_id").(string),
		SubnetID: d.Get("subnet_id").(string),
	}

	routeTablesList, err := routetables.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve route tables: %w", err)
	}

	routeTables := flattenVpcRouteTablesV1(routeTablesList)

	v, e := d.GetOk("id")
	if e {
		d.SetId(v.(string))
	} else {
		d.SetId(
			fmt.Sprintf("routetables-%s",
				hashcode.Strings(
					[]string{
						config.GetRegion(d),
						d.Get("vpc_id").(string),
						d.Get("subnet_id").(string),
					},
				),
			),
		)
	}

	mErr := multierror.Append(
		d.Set("routetables", routeTables),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}
