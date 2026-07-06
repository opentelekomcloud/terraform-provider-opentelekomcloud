package cc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/policy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCcCentralNetworkPoliciesV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcCentralNetworkPoliciesV3Read,

		Schema: map[string]*schema.Schema{
			"central_network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_applied": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"central_network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document_template_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_applied": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataSourceCentralNetworkPolicyDocumentSchema(),
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

func dataSourceCentralNetworkPolicyDocumentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_plane": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"er_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCentralNetworkPolicyErInstanceSchema(),
			},
			"planes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associate_er_tables": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enterprise_router_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enterprise_router_table_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"exclude_er_connections": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exclude_er_instances": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     dataSourceCentralNetworkPolicyErInstanceSchema(),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCentralNetworkPolicyErInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCcCentralNetworkPoliciesV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CcV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	centralNetworkId := d.Get("central_network_id").(string)

	isApplied, err := boolFilter(d.Get("is_applied").(string))
	if err != nil {
		return diag.Errorf("invalid value for is_applied, expected \"true\" or \"false\": %s", err)
	}

	var policies []policy.CentralNetworkPolicy
	marker := ""
	for {
		resp, err := policy.List(client, policy.ListOpts{
			CentralNetworkId: centralNetworkId,
			Marker:           marker,
			ID:               stringFilter(d.Get("policy_id").(string)),
			State:            stringFilter(d.Get("state").(string)),
			IsApplied:        isApplied,
			Version:          common.ExpandToIntList(d.Get("version").([]interface{})),
		})
		if err != nil {
			return diag.Errorf("error retrieving OpenTelekomCloud CC central network policies: %s", err)
		}
		policies = append(policies, resp.CentralNetworkPolicies...)
		if resp.PageInfo.NextMarker == "" || len(resp.CentralNetworkPolicies) == 0 {
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
		d.Set("policies", flattenCentralNetworkPolicies(policies)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCentralNetworkPolicies(policies []policy.CentralNetworkPolicy) []map[string]interface{} {
	if len(policies) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(policies))
	for i, p := range policies {
		result[i] = map[string]interface{}{
			"id":                        p.ID,
			"central_network_id":        p.CentralNetworkId,
			"domain_id":                 p.DomainId,
			"state":                     p.State,
			"document_template_version": p.DocumentTemplateVersion,
			"is_applied":                p.IsApplied,
			"version":                   p.Version,
			"created_at":                p.CreatedAt,
			"document":                  flattenCentralNetworkPolicyDocument(p.Document),
		}
	}
	return result
}

func flattenCentralNetworkPolicyDocument(document policy.PolicyDocument) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"default_plane": document.DefaultPlane,
			"er_instances":  flattenCentralNetworkPolicyErInstances(document.ErInstances),
			"planes":        flattenDataSourcePolicyPlanes(document.Planes),
		},
	}
}

// flattenDataSourcePolicyPlanes mirrors the resource's flattenCentralNetworkPolicyPlanes but also
// exposes the plane name. The resource omits it because it always enforces a single default plane,
// whereas this data source reads arbitrary policies and reports the name as returned by the API.
func flattenDataSourcePolicyPlanes(planes []policy.PlaneDocument) []interface{} {
	if len(planes) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(planes))
	for _, plane := range planes {
		rst = append(rst, map[string]interface{}{
			"name":                   plane.Name,
			"associate_er_tables":    flattenCentralNetworkPolicyAssociateErTables(plane.AssociateErTables),
			"exclude_er_connections": flattenCentralNetworkPolicyExcludeErConnections(plane.ExcludeErConnections),
		})
	}
	return rst
}
