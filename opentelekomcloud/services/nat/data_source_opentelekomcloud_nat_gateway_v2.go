package nat

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/natgateways"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceNatGatewayV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNatGatewayV2Read,

		Schema: map[string]*schema.Schema{
			"nat_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tenant_id": {
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
				Optional: true,
				Computed: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"internal_network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceNatGatewayV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationClient, err)
	}

	listOpts := natgateways.ListOpts{
		ID:                d.Get("nat_id").(string),
		TenantId:          d.Get("tenant_id").(string),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Spec:              d.Get("spec").(string),
		RouterID:          d.Get("router_id").(string),
		InternalNetworkID: d.Get("internal_network_id").(string),
		Status:            d.Get("status").(string),
	}

	if adminState, ok := d.GetOk("admin_state_up"); ok {
		adminState := adminState.(bool)
		listOpts.AdminStateUp = &adminState
	}

	natGateways, err := natgateways.List(client, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to retrieve NAT gateway pages: %w", err)
	}

	refinedGateways, err := natgateways.ExtractNatGateways(natGateways)
	if err != nil {
		return fmterr.Errorf("error extracting NAT gateways: %w", err)
	}

	if len(refinedGateways) < 1 {
		return fmterr.Errorf("Your query returned no results. Please change your search criteria and try again.")
	} else if len(refinedGateways) > 1 {
		return fmterr.Errorf("your query returned more than one result. Please try a more " +
			"specific search criteria")
	}

	natGateway := refinedGateways[0]

	d.SetId(natGateway.ID)

	mErr := multierror.Append(
		d.Set("tenant_id", natGateway.TenantID),
		d.Set("name", natGateway.Name),
		d.Set("description", natGateway.Description),
		d.Set("spec", natGateway.Spec),
		d.Set("router_id", natGateway.RouterID),
		d.Set("internal_network_id", natGateway.InternalNetworkID),
		d.Set("status", natGateway.Status),
		d.Set("admin_state_up", natGateway.AdminStateUp),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
