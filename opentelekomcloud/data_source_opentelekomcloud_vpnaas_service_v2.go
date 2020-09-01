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

		Schema: map[string]*schema.Schema{},
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
	pager := services.List(networkingClient, listOpts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		vpnList, err := services.ExtractServices(page)
		if err != nil {
			return false, err
		}
		for _, policy := range policyList {
			for _, rule := range policy.Rules {
				if rule == ruleID {
					policyID = policy.ID
					return false, nil
				}
			}
		}
		return true, nil
	})

	if len(refinedVpns) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

}
