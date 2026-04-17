package tms

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/tms/v1/provider"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceResourceTypesV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceTypesV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_global": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func filterResourceTypes(serviceName, region string, providerList []provider.Providers) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(providerList))
	for _, val := range providerList {
		if serviceName != "" && val.Provider != serviceName {
			continue
		}
		for _, resource := range val.ResourceTypes {
			if region != "" && !common.StrSliceContains(resource.Regions, region) {
				continue
			}
			result = append(result, map[string]interface{}{
				"name":         resource.ResourceType,
				"is_global":    resource.Global,
				"display_name": resource.DisplayName,
				"service_name": val.Provider,
			})
		}
	}
	return result
}

func dataSourceResourceTypesV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, tmsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.TmsV1Client()
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	var (
		region      = d.Get("region").(string)
		serviceName = d.Get("service_name").(string)
		opts        = provider.ListOpts{
			Provider: serviceName,
			Limit:    pointerto.Int(100),
		}
	)
	resp, err := provider.List(client, opts)
	if err != nil {
		return diag.Errorf("error querying resource types: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	err = d.Set("types", filterResourceTypes(serviceName, region, resp))
	if err != nil {
		return diag.Errorf("error setting resource types field: %s", err)
	}
	return nil
}
