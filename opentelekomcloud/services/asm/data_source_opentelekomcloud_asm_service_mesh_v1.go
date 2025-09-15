package asm

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/asm/v1/servicemesh"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceASMServiceMeshV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceASMServiceMeshV1Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_meshes": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"proxy_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"include_ip_ranges": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"exclude_ip_ranges": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"exclude_outbound_ports": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"exclude_inbound_ports": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"include_outbound_ports": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"include_inbound_ports": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"telemetry_config_tracing": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"random_sampling_percentage": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"default_providers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"extension_providers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"zipkin_service_addr": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"zipkin_service_port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"cluster_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceASMServiceMeshV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AsmV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	listResp, err := servicemesh.List(client)
	if err != nil {
		return fmterr.Errorf("error fetching asm service mesh : %w", err)
	}

	var meshes []servicemesh.ServiceMesh
	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
		mesh, err := filterMesh(listResp, d.Id())
		if err != nil {
			return diag.Errorf("error in mesh filter: %s", err)
		}
		meshes = append(meshes, mesh)
	} else {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return diag.Errorf("unable to generate ID: %s", err)
		}
		d.SetId(id)
		meshes = listResp
	}

	mErr := multierror.Append(
		d.Set("service_meshes", flattenMeshes(meshes)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func filterMesh(listResp []servicemesh.ServiceMesh, id string) (servicemesh.ServiceMesh, error) {
	for _, mesh := range listResp {
		if mesh.Metadata.UID == id {
			return mesh, nil
		}
	}
	return servicemesh.ServiceMesh{}, fmt.Errorf("error finding service mesh with id: %s", id)
}

func flattenMeshes(meshesInResp []servicemesh.ServiceMesh) []map[string]interface{} {
	var meshes []map[string]interface{}
	for _, meshInResp := range meshesInResp {
		var clusterIds []string
		for _, v := range meshInResp.Spec.ExtendParams.Clusters {
			clusterIds = append(clusterIds, v.ClusterID)
		}
		mesh := map[string]interface{}{
			"id":                       meshInResp.Metadata.UID,
			"name":                     meshInResp.Metadata.Name,
			"type":                     meshInResp.Spec.Type,
			"version":                  meshInResp.Spec.Version,
			"ipv6_enable":              meshInResp.Spec.IPv6Enable,
			"proxy_config":             setProxyConfig(meshInResp.Spec.Config.ProxyConfig),
			"telemetry_config_tracing": setTelemetryConfig(meshInResp.Spec.Config.TelemetryConfig),
			"cluster_ids":              clusterIds,
			"creation_timestamp":       meshInResp.Metadata.CreationTimestamp,
			"status":                   meshInResp.Status.Phase,
		}
		meshes = append(meshes, mesh)
	}
	return meshes
}
