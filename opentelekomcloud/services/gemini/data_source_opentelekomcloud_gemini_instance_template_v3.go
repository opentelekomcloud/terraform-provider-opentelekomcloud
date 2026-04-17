package gemini

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/template"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceGeminiDBV3InstanceTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBInstanceTemplateRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datastore_version_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_name": {
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
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configuration_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restart_required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"readonly": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"value_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
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

func dataSourceGeminiDBInstanceTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)

	parameters, err := template.GetInstanceParameters(client, instanceId)
	if err != nil {
		return diag.Errorf("error getting GeminiDB instance template: %s", err)
	}

	d.SetId(parameters.Id)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("datastore_version_name", parameters.DataStoreVersionName),
		d.Set("datastore_name", parameters.DataStoreName),
		d.Set("created_at", parameters.Created),
		d.Set("updated_at", parameters.Updated),
		d.Set("mode", parameters.Mode),
		d.Set("configuration_parameters", flattenConfigurationParameters(parameters.ConfigurationParameters)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting GeminiDB instance template attributes: %s", err)
	}

	return nil
}

func flattenConfigurationParameters(parameters []template.InstanceParameterResult) []map[string]interface{} {
	result := make([]map[string]interface{}, len(parameters))
	for i, param := range parameters {
		result[i] = map[string]interface{}{
			"name":             param.Name,
			"value":            param.Value,
			"restart_required": param.RestartRequired,
			"readonly":         param.Readonly,
			"value_range":      param.ValueRange,
			"type":             param.Type,
			"description":      param.Description,
		}
	}
	return result
}
