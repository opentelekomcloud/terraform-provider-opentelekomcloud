package apigw

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	apis "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/api"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceApigwApiHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentsRead,

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"history": {
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
						"publish_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenHistory(vList []apis.VersionResp) ([]map[string]interface{}, []string) {
	if len(vList) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(vList))
	ids := make([]string, len(vList))
	for i, version := range vList {
		result[i] = map[string]interface{}{
			"id":           version.VersionID,
			"name":         version.Version,
			"description":  version.Description,
			"publish_time": version.PublishTime,
			"status":       version.Status,
		}
		ids[i] = version.VersionID
	}
	return result, ids
}

func dataSourceEnvironmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	gatewayId := d.Get("gateway_id").(string)
	apiId := d.Get("api_id").(string)
	opts := apis.ListHistoryOpts{
		EnvID:   d.Get("environment_id").(string),
		EnvName: d.Get("environment_name").(string),
	}

	history, err := apis.GetHistory(client, gatewayId, apiId, opts)
	if err != nil {
		return diag.FromErr(err)
	}
	hResult, ids := flattenHistory(history)
	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("history", hResult),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
