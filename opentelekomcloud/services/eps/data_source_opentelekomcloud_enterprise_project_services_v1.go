package eps

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/eps/v1/providers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceEpsServicesV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEpsServicesV1Read,

		Schema: map[string]*schema.Schema{
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildEpsServicesSchema(),
			},
		},
	}
}

func buildEpsServicesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_i18n_display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildEpsResourceTypeSchema(),
			},
		},
	}
}

func buildEpsResourceTypeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type_i18n_display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceEpsServicesV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.EpsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	p, err := providers.List(client, providers.ListOpts{
		Provider: d.Get("service").(string),
		Locale:   d.Get("locale").(string),
		Limit:    200,
		Offset:   0,
	})
	if err != nil {
		return diag.Errorf("error retrieving enterprise project services: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("services", flattenListEpsServicesResponseBody(p)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListEpsServicesResponseBody(respArray []providers.Provider) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"service":                   v.Provider,
			"service_i18n_display_name": v.ProviderI18nDisplay,
			"resource_types":            flattenResourceTypeBody(v.ResourceTypes),
		})
	}
	return rst
}

func flattenResourceTypeBody(respArray []providers.ResourceTypes) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"resource_type":                   v.ResourceType,
			"resource_type_i18n_display_name": v.ResourceTypeI18nDisplayName,
			"global":                          v.Global,
			"regions":                         v.Regions,
		})
	}
	return rst
}
