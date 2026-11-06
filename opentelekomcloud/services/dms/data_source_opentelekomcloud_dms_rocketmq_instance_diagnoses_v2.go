package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/diagnosis"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceDmsRocketMQInstanceDiagnosesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketMQInstanceDiagnosesV2Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"report_id":         {Type: schema.TypeString, Computed: true},
						"group_name":        {Type: schema.TypeString, Computed: true},
						"consumer_nums":     {Type: schema.TypeInt, Computed: true},
						"status":            {Type: schema.TypeString, Computed: true},
						"created_at":        {Type: schema.TypeString, Computed: true},
						"abnormal_item_sum": {Type: schema.TypeInt, Computed: true},
						"faulted_node_sum":  {Type: schema.TypeInt, Computed: true},
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

func dataSourceDmsRocketMQInstanceDiagnosesV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.DmsV2Client(region)
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	opts := diagnosis.ListOpts{Limit: 50}
	var all []diagnosis.ReportSummary
	for {
		resp, err := diagnosis.List(client, instanceID, opts)
		if err != nil {
			return diag.Errorf("error getting diagnosis reports of DMS RocketMQ instance (%s): %s", instanceID, err)
		}
		all = append(all, resp.Reports...)
		if len(resp.Reports) == 0 || len(all) >= resp.TotalNum {
			break
		}
		opts.Offset += len(resp.Reports)
	}

	ids := []string{instanceID}
	result := make([]map[string]interface{}, len(all))
	for i, v := range all {
		ids = append(ids, v.ReportID)
		result[i] = map[string]interface{}{
			"report_id":         v.ReportID,
			"group_name":        v.GroupName,
			"consumer_nums":     v.ConsumerNums,
			"status":            v.Status,
			"created_at":        v.CreatedAt,
			"abnormal_item_sum": v.AbnormalItemSum,
			"faulted_node_sum":  v.FaultedNodeSum,
		}
	}
	log.Printf("[DEBUG] DMS RocketMQ diagnosis reports: %v", ids[1:])
	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("reports", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
