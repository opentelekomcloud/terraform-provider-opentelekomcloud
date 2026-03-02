package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/peerings"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
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
						"peer_vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_tenant_id": {
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
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := peerings.ListOpts{
		Name:       d.Get("name").(string),
		Status:     d.Get("status").(string),
		VpcId:      d.Get("vpc_id").(string),
		Peer_VpcId: d.Get("peer_vpc_id").(string),
		TenantId:   d.Get("peer_tenant_id").(string),
	}

	peeringList, err := peerings.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve VPC peering connections: %s", err)
	}

	uID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uID)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("peering_connections", flattenPeeringConnections(peeringList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPeeringConnections(peeringList []peerings.Peering) []interface{} {
	if peeringList == nil {
		return nil
	}

	result := make([]interface{}, len(peeringList))
	for i, p := range peeringList {
		result[i] = map[string]interface{}{
			"id":             p.ID,
			"name":           p.Name,
			"description":    p.Description,
			"status":         p.Status,
			"vpc_id":         p.RequestVpcInfo.VpcId,
			"peer_vpc_id":    p.AcceptVpcInfo.VpcId,
			"peer_tenant_id": p.AcceptVpcInfo.TenantId,
		}
	}
	return result
}
