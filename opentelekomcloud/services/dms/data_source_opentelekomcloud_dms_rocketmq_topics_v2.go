package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/topics"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceDmsRocketMQTopicsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketMQTopicsV2Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permission": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"sub", "pub", "all"}, false),
			},
			"message_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NORMAL", "FIFO", "DELAY", "TRANSACTION"}, false),
			},
			"topics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":                  {Type: schema.TypeString, Computed: true},
						"total_read_queue_num":  {Type: schema.TypeInt, Computed: true},
						"total_write_queue_num": {Type: schema.TypeInt, Computed: true},
						"permission":            {Type: schema.TypeString, Computed: true},
						"message_type":          {Type: schema.TypeString, Computed: true},
						"brokers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"broker_name":     {Type: schema.TypeString, Computed: true},
									"read_queue_num":  {Type: schema.TypeInt, Computed: true},
									"write_queue_num": {Type: schema.TypeInt, Computed: true},
								},
							},
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

func dataSourceDmsRocketMQTopicsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.DmsV2Client(region)
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	// The API maximum.
	opts := topics.ListOpts{Limit: 50}
	var all []topics.Topic
	for {
		resp, err := topics.List(client, instanceID, opts)
		if err != nil {
			return diag.Errorf("error getting topics of DMS RocketMQ instance (%s): %s", instanceID, err)
		}
		all = append(all, resp.Topics...)
		if len(resp.Topics) == 0 || len(all) >= resp.Total {
			break
		}
		opts.Offset += len(resp.Topics)
	}
	log.Printf("[DEBUG] DMS RocketMQ topic list: %#v", all)

	name := d.Get("name").(string)
	permission := d.Get("permission").(string)
	messageType := d.Get("message_type").(string)

	ids := []string{instanceID}
	result := make([]map[string]interface{}, 0, len(all))
	for _, v := range all {
		if (name != "" && v.Name != name) ||
			(permission != "" && v.Permission != permission) ||
			(messageType != "" && v.MessageType != messageType) {
			continue
		}
		brokers := make([]map[string]interface{}, len(v.Brokers))
		for i, b := range v.Brokers {
			brokers[i] = map[string]interface{}{
				"broker_name":     b.BrokerName,
				"read_queue_num":  b.ReadQueueNum,
				"write_queue_num": b.WriteQueueNum,
			}
		}
		ids = append(ids, v.Name)
		result = append(result, map[string]interface{}{
			"name":                  v.Name,
			"total_read_queue_num":  v.TotalReadQueueNum,
			"total_write_queue_num": v.TotalWriteQueueNum,
			"permission":            v.Permission,
			"message_type":          v.MessageType,
			"brokers":               brokers,
		})
	}
	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("topics", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
