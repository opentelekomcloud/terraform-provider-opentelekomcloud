package cc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/central_network"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCcCentralNetworksV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcCentralNetworksV3Read,

		Schema: map[string]*schema.Schema{
			"central_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"central_networks": {
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
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_plane_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
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
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceCcCentralNetworksV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CcV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var networks []central_network.CentralNetwork
	marker := ""
	for {
		resp, err := central_network.List(client, central_network.ListOpts{
			Marker:              marker,
			ID:                  stringFilter(d.Get("central_network_id").(string)),
			Name:                stringFilter(d.Get("name").(string)),
			State:               stringFilter(d.Get("state").(string)),
			EnterpriseProjectId: stringFilter(d.Get("enterprise_project_id").(string)),
		})
		if err != nil {
			return diag.Errorf("error retrieving OpenTelekomCloud CC central networks: %s", err)
		}
		networks = append(networks, resp.CentralNetworks...)
		if resp.PageInfo.NextMarker == "" || len(resp.CentralNetworks) == 0 {
			break
		}
		marker = resp.PageInfo.NextMarker
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("central_networks", flattenCentralNetworks(networks)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCentralNetworks(networks []central_network.CentralNetwork) []map[string]interface{} {
	if len(networks) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(networks))
	for i, cn := range networks {
		result[i] = map[string]interface{}{
			"id":                    cn.ID,
			"name":                  cn.Name,
			"description":           cn.Description,
			"state":                 cn.State,
			"enterprise_project_id": cn.EnterpriseProjectId,
			"default_plane_id":      cn.DefaultPlaneId,
			"domain_id":             cn.DomainId,
			"created_at":            cn.CreatedAt,
			"updated_at":            cn.UpdatedAt,
		}
	}
	return result
}
