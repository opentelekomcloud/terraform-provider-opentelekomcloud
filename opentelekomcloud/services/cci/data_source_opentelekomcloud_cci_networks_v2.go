package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/network"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCCINetworksV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCINetworksV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finalizers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_families": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"subnets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     networksStatusSchema(),
						},
					},
				},
			},
		},
	}
}

func networksStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_transition_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"subnet_attrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_v4_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_v6_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCCINetworksV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := common.ClientFromCtx(ctx, cciNetworkClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2NetworkClient(region)
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciNetworkClient, err)
	}

	ns := d.Get("namespace").(string)
	resp, err := network.List(client, network.ListOpts{Namespace: ns})
	if err != nil {
		return diag.Errorf("error querying CCI networks: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("networks", flattenNetworks(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNetworks(networks []network.Network) []map[string]interface{} {
	if len(networks) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(networks))
	for _, n := range networks {
		result = append(result, map[string]interface{}{
			"name":               n.Metadata.Name,
			"namespace":          n.Metadata.Namespace,
			"annotations":        n.Metadata.Annotations,
			"labels":             n.Metadata.Labels,
			"creation_timestamp": n.Metadata.CreationTimestamp,
			"resource_version":   n.Metadata.ResourceVersion,
			"finalizers":         n.Metadata.Finalizers,
			"uid":                n.Metadata.UID,
			"ip_families":        n.Spec.IPFamilies,
			"security_group_ids": n.Spec.SecurityGroups,
			"subnets":            flattenNetworkSubnets(n.Spec.Subnets),
			"status":             flattenNetworkStatus(n.Status),
		})
	}
	return result
}
