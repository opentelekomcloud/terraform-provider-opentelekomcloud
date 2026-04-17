package taurusdb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/instance"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceTaurusDBV3MysqlEngineVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTausurDBV3MysqlEngineVersionsRead,

		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datastores": {
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

func dataSourceTausurDBV3MysqlEngineVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	databaseName := d.Get("database_name").(string)

	resp, err := instance.ListDatastores(client, databaseName)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB MySQL engine versions: %s", err)
	}

	datastores := make([]map[string]interface{}, len(resp.Datastores))
	for i, ds := range resp.Datastores {
		datastores[i] = map[string]interface{}{
			"id":   ds.Id,
			"name": ds.Name,
		}
	}

	if err := d.Set("datastores", datastores); err != nil {
		return diag.Errorf("error setting datastores: %s", err)
	}

	if err := d.Set("region", config.GetRegion(d)); err != nil {
		return diag.Errorf("error setting region: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return nil
}
