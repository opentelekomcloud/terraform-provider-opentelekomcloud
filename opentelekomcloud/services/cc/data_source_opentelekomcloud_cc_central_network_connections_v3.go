package cc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/connection"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCcCentralNetworkConnectionsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcCentralNetworkConnectionsV3Read,

		Schema: map[string]*schema.Schema{
			"central_network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connection_id": {
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
			"bandwidth_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"global_connection_bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_cross_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
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
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"central_network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"central_network_plane_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"global_connection_bandwidth_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_frozen": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"connection_type": {
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
						"connection_point_pair": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"site_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parent_instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceCcCentralNetworkConnectionsV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CcV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	centralNetworkId := d.Get("central_network_id").(string)

	isCrossRegion, err := boolFilter(d.Get("is_cross_region").(string))
	if err != nil {
		return diag.Errorf("invalid value for is_cross_region, expected \"true\" or \"false\": %s", err)
	}

	var connections []connection.CentralNetworkConnection
	marker := ""
	for {
		resp, err := connection.List(client, connection.ListOpts{
			CentralNetworkId:            centralNetworkId,
			Marker:                      marker,
			ID:                          stringFilter(d.Get("connection_id").(string)),
			Name:                        stringFilter(d.Get("name").(string)),
			State:                       stringFilter(d.Get("state").(string)),
			BandwidthType:               d.Get("bandwidth_type").(string),
			ConnectionType:              d.Get("connection_type").(string),
			GlobalConnectionBandwidthId: stringFilter(d.Get("global_connection_bandwidth_id").(string)),
			IsCrossRegion:               isCrossRegion,
		})
		if err != nil {
			return diag.Errorf("error retrieving OpenTelekomCloud CC central network connections: %s", err)
		}
		connections = append(connections, resp.CentralNetworkConnections...)
		if resp.PageInfo.NextMarker == "" || len(resp.CentralNetworkConnections) == 0 {
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
		d.Set("connections", flattenCentralNetworkConnections(connections)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCentralNetworkConnections(connections []connection.CentralNetworkConnection) []map[string]interface{} {
	if len(connections) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(connections))
	for i, c := range connections {
		result[i] = map[string]interface{}{
			"id":                             c.ID,
			"name":                           c.Name,
			"description":                    c.Description,
			"domain_id":                      c.DomainId,
			"enterprise_project_id":          c.EnterpriseProjectId,
			"central_network_id":             c.CentralNetworkId,
			"central_network_plane_id":       c.CentralNetworkPlaneId,
			"global_connection_bandwidth_id": c.GlobalConnectionBandwidthId,
			"bandwidth_type":                 c.BandwidthType,
			"bandwidth_size":                 c.BandwidthSize,
			"state":                          c.State,
			"is_frozen":                      c.IsFrozen,
			"connection_type":                c.ConnectionType,
			"created_at":                     c.CreatedAt,
			"updated_at":                     c.UpdatedAt,
			"connection_point_pair":          flattenConnectionPointPair(c.ConnectionPointPair),
		}
	}
	return result
}

func flattenConnectionPointPair(points []connection.ConnectionPoint) []map[string]interface{} {
	if len(points) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(points))
	for i, p := range points {
		result[i] = map[string]interface{}{
			"id":                 p.ID,
			"project_id":         p.ProjectId,
			"region_id":          p.RegionId,
			"site_code":          p.SiteCode,
			"instance_id":        p.InstanceId,
			"parent_instance_id": p.ParentInstanceId,
			"type":               p.Type,
		}
	}
	return result
}
