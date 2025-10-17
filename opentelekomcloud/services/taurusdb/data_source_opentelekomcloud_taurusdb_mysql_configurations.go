package taurusdb

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/template"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceTaurusDBV3MysqlConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusdbMysqlConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"configurations": {
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
						"user_defined": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore_version_name": {
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

func dataSourceTaurusdbMysqlConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	opts := template.ListConfigurationsOpts{
		Offset: 0,
		Limit:  100,
	}

	allConfigurations := make([]template.ConfigurationSummary, 0)
	for {
		configurations, err := template.ListConfigurations(client, opts)
		if err != nil {
			return diag.Errorf("error listing TaurusDB MySQL configurations: %s", err)
		}

		if len(configurations) == 0 {
			break
		}

		allConfigurations = append(allConfigurations, configurations...)

		if len(configurations) < opts.Limit {
			break
		}
		opts.Offset += opts.Limit
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("configurations", flattenConfigurations(allConfigurations)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting TaurusDB MySQL configurations attributes: %s", err)
	}

	return nil
}

func flattenConfigurations(configurations []template.ConfigurationSummary) []map[string]interface{} {
	result := make([]map[string]interface{}, len(configurations))
	for i, conf := range configurations {
		result[i] = map[string]interface{}{
			"id":                     conf.Id,
			"name":                   conf.Name,
			"user_defined":           conf.UserDefined,
			"description":            conf.Description,
			"datastore_name":         conf.DatastoreName,
			"datastore_version_name": conf.DatastoreVersionName,
			"created_at":             conf.Created,
			"updated_at":             conf.Updated,
		}
	}
	return result
}
