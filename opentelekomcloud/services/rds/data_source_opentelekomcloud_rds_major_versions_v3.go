package rds

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/upgrade"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRdsAvailableVersionV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsAvailableVersionV3Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"available_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRdsAvailableVersionV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	instanceId := d.Get("instance_id").(string)

	opts := upgrade.GetAvailableVersionOpts{
		InstanceId: instanceId,
	}

	versions, err := upgrade.GetAvailableVersion(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("available_versions", versions.AvailableVersions); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceId)
	return nil
}
