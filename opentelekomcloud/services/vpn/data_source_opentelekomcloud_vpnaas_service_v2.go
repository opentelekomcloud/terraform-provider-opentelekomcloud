package vpn

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVpnServiceV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnServiceV2Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_v6_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_v4_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVpnServiceV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}
	adminStateUp := d.Get("admin_state_up").(bool)
	listOpts := services.ListOpts{
		TenantID:     d.Get("tenant_id").(string),
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		AdminStateUp: &adminStateUp,
		Status:       d.Get("status").(string),
		SubnetID:     d.Get("subnet_id").(string),
		RouterID:     d.Get("router_id").(string),
		ProjectID:    d.Get("project_id").(string),
		ExternalV6IP: d.Get("external_v6_ip").(string),
		ExternalV4IP: d.Get("external_v4_ip").(string),
		FlavorID:     d.Get("flavor_id").(string),
	}
	var refinedVpns []services.Service

	pager := services.List(networkingClient, listOpts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		vpnList, err := services.ExtractServices(page)
		if err != nil {
			return false, err
		}
		refinedVpns = append(refinedVpns, vpnList...)
		return true, nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if len(refinedVpns) < 1 {
		return fmterr.Errorf("your query returned zero results. Please change your search criteria and try again.")
	}

	if len(refinedVpns) > 1 {
		return fmterr.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}
	Vpn := refinedVpns[0]

	log.Printf("[INFO] Retrieved Vpn using given filter %s: %+v", Vpn.ID, Vpn)
	d.SetId(Vpn.ID)

	mErr := multierror.Append(
		d.Set("id", Vpn.ID),
		d.Set("tenant_id", Vpn.TenantID),
		d.Set("name", Vpn.Name),
		d.Set("subnet_id", Vpn.SubnetID),
		d.Set("admin_state_up", Vpn.AdminStateUp),
		d.Set("external_v4_ip", Vpn.ExternalV4IP),
		d.Set("external_v6_ip", Vpn.ExternalV6IP),
		d.Set("project_id", Vpn.ProjectID),
		d.Set("router_id", Vpn.RouterID),
		d.Set("flavor_id", Vpn.FlavorID),
		d.Set("status", Vpn.Status),
		d.Set("description", Vpn.Description),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
