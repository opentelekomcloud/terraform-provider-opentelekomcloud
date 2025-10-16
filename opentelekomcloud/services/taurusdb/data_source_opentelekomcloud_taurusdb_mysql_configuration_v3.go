package taurusdb

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/template"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceTaurusDBV3MysqlConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBMysqlConfigurationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBMysqlConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	opts := template.ListConfigurationsOpts{
		Limit: 100,
	}

	configsList, err := template.ListConfigurations(client, opts)
	if err != nil {
		return diag.Errorf("unable to retrieve configurations: %s", err)
	}

	if len(configsList) < 1 {
		return diag.Errorf("your query returned no results. " +
			"please change your search criteria and try again.")
	}

	if name, ok := d.GetOk("name"); ok {
		var filteredConfigs []template.ConfigurationSummary
		for _, conf := range configsList {
			if conf.Name == name.(string) {
				filteredConfigs = append(filteredConfigs, conf)
			}
		}
		configsList = filteredConfigs
	}

	if len(configsList) < 1 {
		return diag.Errorf("your query returned no results. " +
			"please change your search criteria and try again.")
	}

	configuration := configsList[0]
	d.SetId(configuration.Id)

	mErr := multierror.Append(
		d.Set("name", configuration.Name),
		d.Set("description", configuration.Description),
		d.Set("datastore_version", configuration.DatastoreVersionName),
		d.Set("datastore_name", configuration.DatastoreName),
		d.Set("region", config.GetRegion(d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
