package vpc

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/networks"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceNetworkingNetworkV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingNetworkV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"matching_subnet_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: common.Descriptions["tenant_id"],
			},
			"admin_state_up": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingNetworkV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating networking client: %w", err)
	}

	_, id := ExtractValSFromNid(d.Get("network_id").(string))
	listOpts := networks.ListOpts{
		ID:       id,
		Name:     d.Get("name").(string),
		TenantID: d.Get("tenant_id").(string),
		Status:   "ACTIVE",
	}

	pages, err := networks.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allNetworks, err := networks.ExtractNetworks(pages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve networks: %s", err)
	}

	var refinedNetworks []networks.Network
	if cidr := d.Get("matching_subnet_cidr").(string); cidr != "" {
		for _, n := range allNetworks {
			for _, s := range n.Subnets {
				subnet, err := subnets.Get(networkingClient, s).Extract()
				if err != nil {
					if _, ok := err.(golangsdk.ErrDefault404); ok {
						continue
					}
					return fmterr.Errorf("unable to retrieve network subnet: %s", err)
				}
				if cidr == subnet.CIDR {
					refinedNetworks = append(refinedNetworks, n)
				}
			}
		}
	} else {
		refinedNetworks = allNetworks
	}

	if len(refinedNetworks) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNetworks) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	network := refinedNetworks[0]

	log.Printf("[DEBUG] Retrieved Network %s: %+v", network.ID, network)
	d.SetId(network.ID)

	mErr := multierror.Append(
		d.Set("name", network.Name),
		d.Set("admin_state_up", strconv.FormatBool(network.AdminStateUp)),
		d.Set("shared", strconv.FormatBool(network.Shared)),
		d.Set("tenant_id", network.TenantID),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
