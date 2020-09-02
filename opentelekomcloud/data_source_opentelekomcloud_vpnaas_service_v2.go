package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/huaweicloud/golangsdk/pagination"
)

func dataSourceVpnServiceV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpnServiceV2Read,

		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type: schema.TypeString,
			},
			"name": {
				Type: schema.TypeString,
			},
			"description": {
				Type: schema.TypeString,
			},
			"admin_state_up": {
				Type: schema.TypeBool,
			},
			"status": {
				Type: schema.TypeString,
			},
			"subnet_id": {
				Type: schema.TypeString,
			},
			"router_id": {
				Type: schema.TypeString,
			},
			"project_id": {
				Type: schema.TypeString,
			},
			"flavor_id": {
				Type: schema.TypeString,
			},
			"external_v6_ip": {
				Type: schema.TypeString,
			},
			"external_v4_ip": {
				Type: schema.TypeString,
			},
		},
	}
}

func dataSourceVpnServiceV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
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
	var vpns []services.Service

	pager := services.List(networkingClient, listOpts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		vpnList, err := services.ExtractServices(page)
		if err != nil {
			return false, err
		}
		for _, vpn := range vpnList {
			if vpn.ID != "" {
				vpns = append(vpns, vpn)
			}
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	if len(vpns) < 1 {
		return fmt.Errorf("Your query returned zero results. Please change your search criteria and try again.")
	}

	d.SetId(d.Get("cluster_id").(string))
	d.Set("ids", vpns)
	d.Set("region", GetRegion(d, config))
	return nil
}
