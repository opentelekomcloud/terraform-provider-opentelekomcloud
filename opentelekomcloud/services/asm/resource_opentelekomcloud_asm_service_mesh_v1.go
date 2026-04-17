package asm

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/asm/v1/servicemesh"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceASMServiceMeshV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASMServiceMeshV1Create,
		ReadContext:   resourceASMServiceMeshV1Read,
		DeleteContext: resourceASMServiceMeshV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"installation_nodes": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"injection_namespaces": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"proxy_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include_ip_ranges": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"exclude_ip_ranges": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"exclude_outbound_ports": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"exclude_inbound_ports": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"include_outbound_ports": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"include_inbound_ports": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"telemetry_config_tracing": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"random_sampling_percentage": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"default_providers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"extension_providers": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"zipkin_service_addr": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"zipkin_service_port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
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
	}
}

func resourceASMServiceMeshV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AsmV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	var meshConfig servicemesh.MeshConfig
	if d.Get("proxy_config.#").(int) != 0 {
		meshConfig.ProxyConfig = getProxyConfig(d)
	}
	if d.Get("telemetry_config_tracing.#").(int) != 0 {
		meshConfig.TelemetryConfig = getTelemetryConfig(d)
	}

	createOpts := servicemesh.CreateOpts{
		APIVersion: "v1",
		Kind:       "mesh",
		Metadata: servicemesh.MeshMetadata{
			Name: d.Get("name").(string),
		},
		Spec: servicemesh.MeshSpec{
			Type:    d.Get("type").(string),
			Version: d.Get("version").(string),
			ExtendParams: &servicemesh.MeshExtendParams{
				Clusters: getClusters(d),
			},
			IPv6Enable: d.Get("ipv6_enable").(bool),
			Config:     &meshConfig,
		},
	}

	createResp, err := servicemesh.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating asm service mesh: %w", err)
	}
	d.SetId(createResp.Metadata.UID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Creating"},
		Target:     []string{"Running"},
		Refresh:    instanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for ASM service mesh (%s) to become ready: %w", d.Id(), err)
	}

	log.Printf("Created ASM Service Mesh %s: %#v", d.Id(), createResp.Spec)

	return resourceASMServiceMeshV1Read(ctx, d, meta)
}

func resourceASMServiceMeshV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AsmV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	getResp, err := servicemesh.Get(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching asm service mesh : %w", err)
	}

	var clusterIds []string
	for _, v := range getResp.Spec.ExtendParams.Clusters {
		clusterIds = append(clusterIds, v.ClusterID)
	}

	mErr := multierror.Append(
		d.Set("name", getResp.Metadata.Name),
		d.Set("type", getResp.Spec.Type),
		d.Set("version", getResp.Spec.Version),
		d.Set("ipv6_enable", getResp.Spec.IPv6Enable),
		d.Set("proxy_config", setProxyConfig(getResp.Spec.Config.ProxyConfig)),
		d.Set("telemetry_config_tracing", setTelemetryConfig(getResp.Spec.Config.TelemetryConfig)),
		d.Set("cluster_ids", clusterIds),
		d.Set("creation_timestamp", getResp.Metadata.CreationTimestamp),
		d.Set("status", getResp.Status.Phase),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceASMServiceMeshV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AsmV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	err = servicemesh.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting ASM Service Mesh: %w", err)
	}

	err = WaitForDeleteServiceMesh(client, 600, 5, d.Id())
	if err != nil {
		return fmterr.Errorf("error waiting for ASM service Mesh (%s) to be deleted: %w", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getClusters(d *schema.ResourceData) []servicemesh.MeshCluster {
	clustersInput := d.Get("clusters").([]interface{})
	result := make([]servicemesh.MeshCluster, 0, len(clustersInput))

	for _, clusterInputRaw := range clustersInput {
		clusterInput := clusterInputRaw.(map[string]interface{})
		installationNodes := clusterInput["installation_nodes"].([]interface{})
		injectionNamespaces := clusterInput["injection_namespaces"].([]interface{})
		cluster := servicemesh.MeshCluster{
			ClusterID: clusterInput["cluster_id"].(string),
			Installation: &servicemesh.InstallationConfig{
				Nodes: &servicemesh.Selector{
					FieldSelector: &servicemesh.FieldSelector{
						Key:      "UID",
						Operator: "In",
						Values:   common.ExpandToStringList(installationNodes),
					},
				},
			},
		}
		if len(injectionNamespaces) > 0 {
			cluster.Injection = &servicemesh.InjectionConfig{
				Namespaces: &servicemesh.Selector{
					FieldSelector: &servicemesh.FieldSelector{
						Key:      "Name",
						Operator: "In",
						Values:   common.ExpandToStringList(injectionNamespaces),
					},
				},
			}
		}
		result = append(result, cluster)
	}
	return result
}

func getProxyConfig(d *schema.ResourceData) *servicemesh.ProxyConfig {
	proxyConfigList := d.Get("proxy_config").([]interface{})
	proxyConfig := proxyConfigList[0].(map[string]interface{})
	result := servicemesh.ProxyConfig{
		IncludeIPRanges:      proxyConfig["include_ip_ranges"].(string),
		ExcludeIPRanges:      proxyConfig["exclude_ip_ranges"].(string),
		ExcludeOutboundPorts: proxyConfig["exclude_outbound_ports"].(string),
		ExcludeInboundPorts:  proxyConfig["exclude_inbound_ports"].(string),
		IncludeOutboundPorts: proxyConfig["include_outbound_ports"].(string),
		IncludeInboundPorts:  proxyConfig["include_inbound_ports"].(string),
	}

	return &result
}

func getTelemetryConfig(d *schema.ResourceData) *servicemesh.TelemetryConfig {
	telemetryConfigList := d.Get("telemetry_config_tracing").([]interface{})
	telemetryConfig := telemetryConfigList[0].(map[string]interface{})
	result := servicemesh.TelemetryConfig{
		Tracing: &servicemesh.Tracing{
			RandomSamplingPercentage: telemetryConfig["random_sampling_percentage"].(float64),
			DefaultProviders:         common.ExpandToStringList(telemetryConfig["default_providers"].([]interface{})),
			ExtensionProviders:       getExtensionProviders(telemetryConfig["extension_providers"].([]interface{})),
		},
	}

	return &result
}

func getExtensionProviders(extensionProvidersInput []interface{}) []servicemesh.TracingExtensionProvider {
	result := make([]servicemesh.TracingExtensionProvider, 0, len(extensionProvidersInput))
	for _, val := range extensionProvidersInput {
		extensionProviderInput := val.(map[string]interface{})
		extensionProvider := servicemesh.TracingExtensionProvider{
			Name: extensionProviderInput["name"].(string),
			Zipkin: servicemesh.ZipkinTracingProvider{
				Service: extensionProviderInput["zipkin_service_addr"].(string),
				Port:    extensionProviderInput["zipkin_service_port"].(int),
			},
		}
		result = append(result, extensionProvider)
	}
	return result
}

func setProxyConfig(proxyConfigInResp servicemesh.ProxyConfigResponse) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"include_ip_ranges":      proxyConfigInResp.IncludeIPRanges,
			"exclude_ip_ranges":      proxyConfigInResp.ExcludeIPRanges,
			"exclude_outbound_ports": proxyConfigInResp.ExcludeOutboundPorts,
			"exclude_inbound_ports":  proxyConfigInResp.ExcludeInboundPorts,
			"include_outbound_ports": proxyConfigInResp.IncludeOutboundPorts,
			"include_inbound_ports":  proxyConfigInResp.IncludeInboundPorts,
		},
	}
}

func setTelemetryConfig(telemetryConfigInResp servicemesh.TelemetryConfigResponse) []map[string]interface{} {
	var extensionProviders []map[string]interface{}
	for _, extensionProviderInResp := range telemetryConfigInResp.Tracing.ExtensionProviders {
		extensionProvider := map[string]interface{}{
			"name":                extensionProviderInResp.Name,
			"zipkin_service_addr": extensionProviderInResp.Zipkin.Service,
			"zipkin_service_port": extensionProviderInResp.Zipkin.Port,
		}
		extensionProviders = append(extensionProviders, extensionProvider)
	}

	return []map[string]interface{}{
		{
			"random_sampling_percentage": telemetryConfigInResp.Tracing.RandomSamplingPercentage,
			"default_providers":          telemetryConfigInResp.Tracing.DefaultProviders,
			"extension_providers":        extensionProviders,
		},
	}
}
