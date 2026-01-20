package cce

import (
	"context"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/clusters"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodepools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCCENodePoolConfigV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCCENodePoolConfigV3Create,
		ReadContext:   resourceCCENodePoolConfigV3Read,
		UpdateContext: resourceCCENodePoolConfigV3Update,
		DeleteContext: resourceCCENodePoolConfigV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),

			// used for cluster waiting
			Default: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nodepool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"packages": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"configurations": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
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

func resourceCCENodePoolConfigV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CceV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(cceClientError, err)
	}

	clusterId := d.Get("cluster_id").(string)
	nodepool_id := d.Get("nodepool_id").(string)

	cfgs := d.Get("packages").([]interface{})
	pkgs, err := getPackageConfiguration(cfgs)
	if err != nil {
		return diag.FromErr(err)
	}
	configOpts := nodepools.UpdateConfigurationOpts{
		Kind:       "Configuration",
		APIVersion: "v3",
		Metadata: nodepools.ConfigurationMetadata{
			Name:   d.Get("name").(string),
			Labels: getLabels(d),
		},
		Spec: nodepools.ClusterConfigurationsSpec{
			Packages: pkgs,
		},
	}
	_, err = nodepools.UpdateConfiguration(client, clusterId, nodepool_id, configOpts)
	if err != nil {
		return diag.Errorf("error updating CCE cluster node pool configurations: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceCCENodePoolConfigV3Read(clientCtx, d, meta)
}

func resourceCCENodePoolConfigV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceCCENodePoolConfigV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CceV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(cceClientError, err)
	}

	clusterId := d.Get("cluster_id").(string)
	nodepool_id := d.Get("nodepool_id").(string)

	cfgs := d.Get("configuration").([]interface{})
	pkgs, err := getPackageConfiguration(cfgs)
	if err != nil {
		return diag.FromErr(err)
	}
	configOpts := nodepools.UpdateConfigurationOpts{
		Kind:       "Configuration",
		APIVersion: "v3",
		Metadata: nodepools.ConfigurationMetadata{
			Name:   d.Get("name").(string),
			Labels: getLabels(d),
		},
		Spec: nodepools.ClusterConfigurationsSpec{
			Packages: pkgs,
		},
	}
	_, err = nodepools.UpdateConfiguration(client, clusterId, nodepool_id, configOpts)
	if err != nil {
		return diag.Errorf("error updating CCE cluster node pool configurations: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceCCENodePoolConfigV3Read(clientCtx, d, meta)
}

func resourceCCENodePoolConfigV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The action resource is only removed from the state the task remains in the cloud.",
		},
	}
}

func getLabels(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func getPackageConfiguration(cfgs []interface{}) ([]clusters.PackageConfiguration, error) {
	result := make([]clusters.PackageConfiguration, 0, len(cfgs))
	for _, p := range cfgs {
		pm := p.(map[string]interface{})
		pkg := clusters.PackageConfiguration{
			Name: pm["name"].(string),
		}
		if v, ok := pm["configurations"]; ok {
			itemsRaw := v.([]interface{})
			items := make([]clusters.Configuration, 0, len(itemsRaw))
			for _, it := range itemsRaw {
				im := it.(map[string]interface{})
				name := im["name"].(string)
				valStr := im["value"].(string)
				items = append(items, clusters.Configuration{
					Name:  name,
					Value: common.ParseAnyType(valStr),
				})
			}
			pkg.Configurations = items
		}
		result = append(result, pkg)
	}
	return result, nil
}
