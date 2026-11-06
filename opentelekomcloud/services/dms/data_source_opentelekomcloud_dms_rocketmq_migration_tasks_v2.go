package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/migration"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceDmsRocketMQMigrationTasksV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketMQMigrationTasksV2Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"rocketmq", "rabbitToRocket"}, false),
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":         {Type: schema.TypeString, Computed: true},
						"name":       {Type: schema.TypeString, Computed: true},
						"start_date": {Type: schema.TypeString, Computed: true},
						"status":     {Type: schema.TypeString, Computed: true},
						"type":       {Type: schema.TypeString, Computed: true},
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

func dataSourceDmsRocketMQMigrationTasksV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.DmsV2Client(region)
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	// offset of this API is a page number starting from 1
	const limit = 50
	var all []migration.Task
	for page := 1; ; page++ {
		resp, err := migration.List(client, instanceID, migration.ListOpts{Offset: page, Limit: limit})
		if err != nil {
			return diag.Errorf("error getting migration tasks of DMS RocketMQ instance (%s): %s", instanceID, err)
		}
		all = append(all, resp.Tasks...)
		if len(resp.Tasks) == 0 || page*limit >= resp.Total {
			break
		}
	}

	taskID := d.Get("task_id").(string)
	name := d.Get("name").(string)
	taskType := d.Get("type").(string)
	ids := []string{instanceID}
	result := make([]map[string]interface{}, 0, len(all))
	for _, v := range all {
		if (taskID != "" && v.ID != taskID) || (name != "" && v.Name != name) || (taskType != "" && v.Type != taskType) {
			continue
		}
		ids = append(ids, v.ID)
		result = append(result, map[string]interface{}{
			"id":         v.ID,
			"name":       v.Name,
			"start_date": v.StartDate,
			"status":     v.Status,
			"type":       v.Type,
		})
	}
	log.Printf("[DEBUG] DMS RocketMQ migration tasks: %v", ids[1:])
	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
