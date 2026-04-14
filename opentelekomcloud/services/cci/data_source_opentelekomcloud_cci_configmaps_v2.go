package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/configmap"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCCIConfigMapsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCIConfigMapsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_maps": {
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
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
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
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"binary_data": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"immutable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCCIConfigMapsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(region)
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)

	cmList := make([]configmap.ConfigMap, 0)
	if name, ok := d.GetOk("name"); ok {
		resp, err := configmap.Get(client, ns, name.(string))
		if err != nil {
			return diag.Errorf("error getting the CCI v2 ConfigMap (%s/%s) from the server: %s", ns, name.(string), err)
		}
		cmList = append(cmList, *resp)
	} else {
		resp, err := configmap.List(client, ns, configmap.ListOpts{})
		if err != nil {
			return diag.Errorf("error querying CCI v2 ConfigMaps under namespace %s: %s", ns, err)
		}
		cmList = resp
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("config_maps", flattenConfigMaps(cmList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConfigMaps(cmList []configmap.ConfigMap) []map[string]interface{} {
	if len(cmList) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(cmList))
	for _, cm := range cmList {
		result = append(result, map[string]interface{}{
			"name":               cm.Metadata.Name,
			"namespace":          cm.Metadata.Namespace,
			"api_version":        cm.APIVersion,
			"kind":               cm.Kind,
			"annotations":        cm.Metadata.Annotations,
			"labels":             cm.Metadata.Labels,
			"creation_timestamp": cm.Metadata.CreationTimestamp,
			"resource_version":   cm.Metadata.ResourceVersion,
			"uid":                cm.Metadata.UID,
			"data":               cm.Data,
			"binary_data":        cm.BinaryData,
			"immutable":          cm.Immutable,
		})
	}
	return result
}
