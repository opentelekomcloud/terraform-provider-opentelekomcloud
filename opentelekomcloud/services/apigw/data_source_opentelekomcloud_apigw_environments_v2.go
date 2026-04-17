package apigw

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceApigwEnvironmentsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApigwEnvironmentsV2Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environments": {
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
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func flattenEnvironments(envList []env.EnvResp) ([]map[string]interface{}, []string) {
	if len(envList) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(envList))
	ids := make([]string, len(envList))
	for i, e := range envList {
		result[i] = map[string]interface{}{
			"id":          e.ID,
			"name":        e.Name,
			"description": e.Description,
			"created_at":  e.CreateTime,
		}
		ids[i] = e.ID
	}
	return result, ids
}

func dataSourceApigwEnvironmentsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := env.ListOpts{
		Name:      d.Get("name").(string),
		GatewayID: d.Get("instance_id").(string),
	}
	resp, err := env.List(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	envResult, ids := flattenEnvironments(resp)
	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("environments", envResult),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
