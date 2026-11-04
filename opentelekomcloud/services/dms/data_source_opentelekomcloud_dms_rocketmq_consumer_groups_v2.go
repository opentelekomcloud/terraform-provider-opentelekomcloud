package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceDmsRocketMQConsumerGroupsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketMQConsumerGroupsV2Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":            {Type: schema.TypeString, Computed: true},
						"description":     {Type: schema.TypeString, Computed: true},
						"enabled":         {Type: schema.TypeBool, Computed: true},
						"broadcast":       {Type: schema.TypeBool, Computed: true},
						"brokers":         {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"retry_max_time":  {Type: schema.TypeInt, Computed: true},
						"permissions":     {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"consume_orderly": {Type: schema.TypeBool, Computed: true},
						"group_online":    {Type: schema.TypeBool, Computed: true},
						"created_at":      {Type: schema.TypeInt, Computed: true},
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

func dataSourceDmsRocketMQConsumerGroupsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.DmsV2Client(region)
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	opts := groups.ListOpts{Group: name, Limit: 50}
	var all []groups.ConsumerGroup
	for {
		resp, err := groups.List(client, instanceID, opts)
		if err != nil {
			return diag.Errorf("error getting consumer groups of DMS RocketMQ instance (%s): %s", instanceID, err)
		}
		all = append(all, resp.Groups...)
		if len(resp.Groups) == 0 || len(all) >= resp.Total {
			break
		}
		opts.Offset += len(resp.Groups)
	}
	log.Printf("[DEBUG] DMS RocketMQ consumer group list: %#v", all)

	ids := []string{instanceID}
	result := make([]map[string]interface{}, 0, len(all))
	for _, v := range all {
		// The group query parameter may match by substring.
		if name != "" && v.Name != name {
			continue
		}
		ids = append(ids, v.Name)
		result = append(result, map[string]interface{}{
			"name":            v.Name,
			"description":     v.Description,
			"enabled":         v.Enabled,
			"broadcast":       v.Broadcast,
			"brokers":         v.Brokers,
			"retry_max_time":  v.RetryMaxTime,
			"permissions":     v.Permissions,
			"consume_orderly": v.ConsumeOrderly,
			"group_online":    v.GroupOnline,
			"created_at":      v.CreatedAt,
		})
	}
	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
