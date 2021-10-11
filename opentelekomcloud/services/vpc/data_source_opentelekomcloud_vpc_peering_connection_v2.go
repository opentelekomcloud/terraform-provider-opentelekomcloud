package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/peerings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVpcPeeringConnectionV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcPeeringConnectionV2Read,

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
		},
	}
}

func dataSourceVpcPeeringConnectionV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	peeringClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := peerings.ListOpts{
		ID:         d.Get("id").(string),
		Name:       d.Get("name").(string),
		Status:     d.Get("status").(string),
		TenantId:   d.Get("peer_tenant_id").(string),
		VpcId:      d.Get("vpc_id").(string),
		Peer_VpcId: d.Get("peer_vpc_id").(string),
	}

	refinedPeering, err := peerings.List(peeringClient, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve vpc peering connections: %s", err)
	}

	if len(refinedPeering) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedPeering) > 1 {
		return fmterr.Errorf("multiple VPC peering connections matched." +
			" Use additional constraints to reduce matches to a single VPC peering connection")
	}

	Peering := refinedPeering[0]

	log.Printf("[INFO] Retrieved Vpc peering Connections using given filter %s: %+v", Peering.ID, Peering)
	d.SetId(Peering.ID)

	mErr := multierror.Append(
		d.Set("name", Peering.Name),
		d.Set("status", Peering.Status),
		d.Set("vpc_id", Peering.RequestVpcInfo.VpcId),
		d.Set("peer_vpc_id", Peering.AcceptVpcInfo.VpcId),
		d.Set("peer_tenant_id", Peering.AcceptVpcInfo.TenantId),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting VPC peering attributes: %w", err)
	}

	return nil
}
