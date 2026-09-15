package vpc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v2/peerings"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceVpcPeeringConnectionsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcPeeringConnectionsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: common.ValidateName,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"vpc_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peer_vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"peer_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_connections": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_tenant_id": {
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
				},
			},
		},
	}
}

func dataSourceVpcPeeringConnectionsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.VpcV2Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := peerings.ListOpts{
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		VpcID:    d.Get("vpc_id").(string),
		TenantID: d.Get("peer_tenant_id").(string),
	}

	peeringList, err := peerings.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve VPC peering connections: %s", err)
	}

	peeringList = filterVpcPeerings(
		peeringList,
		d.Get("vpc_id").(string),
		d.Get("vpc_tenant_id").(string),
		d.Get("peer_vpc_id").(string),
		d.Get("peer_tenant_id").(string),
	)

	stateParts := []string{
		config.GetRegion(d),
		d.Get("name").(string),
		d.Get("status").(string),
		d.Get("vpc_id").(string),
		d.Get("vpc_tenant_id").(string),
		d.Get("peer_vpc_id").(string),
		d.Get("peer_tenant_id").(string),
	}
	d.SetId(fmt.Sprintf("vpc-peerings-%s", hashcode.Strings(stateParts)))

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("peering_connections", flattenPeeringConnections(peeringList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
