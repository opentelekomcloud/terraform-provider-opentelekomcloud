package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v2/peerings"

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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: common.ValidateName,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceVpcPeeringConnectionV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	peeringClient, err := config.VpcV2Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := peerings.ListOpts{
		ID:       d.Get("id").(string),
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		TenantID: d.Get("peer_tenant_id").(string),
		VpcID:    d.Get("vpc_id").(string),
	}

	refinedPeering, err := peerings.List(peeringClient, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve vpc peering connections: %s", err)
	}

	refinedPeering = filterVpcPeerings(
		refinedPeering,
		d.Get("vpc_id").(string),
		d.Get("vpc_tenant_id").(string),
		d.Get("peer_vpc_id").(string),
		d.Get("peer_tenant_id").(string),
	)

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

	if err := setVpcPeeringV2Fields(d, &Peering, config.GetRegion(d)); err != nil {
		return fmterr.Errorf("error setting VPC peering attributes: %w", err)
	}

	return nil
}
