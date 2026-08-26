package vpc

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/flow_logs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceVpcFlowLogsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcFlowLogsV1Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: common.ValidateName,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"port", "vpc", "network",
				}, true),
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"traffic_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"all", "accept", "reject",
				}, true),
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ACTIVE", "DOWN", "ERROR",
				}, false),
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2000,
				ValidateFunc: validation.IntBetween(0, 2000),
			},
			"marker": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: flowLogDataSourceSchema(),
				},
			},
		},
	}
}

func dataSourceVpcFlowLogsV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud vpc client: %s", err)
	}

	opts := flow_logs.ListOpts{
		ID:           d.Get("id").(string),
		Name:         d.Get("name").(string),
		TenantID:     d.Get("tenant_id").(string),
		Description:  d.Get("description").(string),
		ResourceType: d.Get("resource_type").(string),
		ResourceID:   d.Get("resource_id").(string),
		TrafficType:  d.Get("traffic_type").(string),
		LogGroupID:   d.Get("log_group_id").(string),
		LogTopicID:   d.Get("log_topic_id").(string),
		Status:       d.Get("status").(string),
		Marker:       d.Get("marker").(string),
	}
	limit := d.Get("limit").(int)
	opts.Limit = &limit

	flowLogs, err := flow_logs.List(client, opts)
	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud VPC flow logs: %s", err)
	}
	sort.Slice(flowLogs, func(i, j int) bool {
		return flowLogs[i].ID < flowLogs[j].ID
	})

	items := make([]map[string]interface{}, len(flowLogs))
	stateParts := []string{
		config.GetRegion(d),
		opts.ID,
		opts.Name,
		opts.TenantID,
		opts.Description,
		opts.ResourceType,
		opts.ResourceID,
		opts.TrafficType,
		opts.LogGroupID,
		opts.LogTopicID,
		opts.Status,
		opts.Marker,
	}
	if opts.Limit != nil {
		stateParts = append(stateParts, fmt.Sprintf("%d", *opts.Limit))
	}
	for i, flowLog := range flowLogs {
		items[i] = flattenVpcFlowLog(flowLog)
		stateParts = append(stateParts, flowLog.ID)
	}

	d.SetId(fmt.Sprintf("vpc-flow-logs-%s", hashcode.Strings(stateParts)))
	mErr := multierror.Append(nil, d.Set("flow_logs", items))
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flowLogDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id":            {Type: schema.TypeString, Computed: true},
		"name":          {Type: schema.TypeString, Computed: true},
		"tenant_id":     {Type: schema.TypeString, Computed: true},
		"description":   {Type: schema.TypeString, Computed: true},
		"resource_type": {Type: schema.TypeString, Computed: true},
		"resource_id":   {Type: schema.TypeString, Computed: true},
		"traffic_type":  {Type: schema.TypeString, Computed: true},
		"log_group_id":  {Type: schema.TypeString, Computed: true},
		"log_topic_id":  {Type: schema.TypeString, Computed: true},
		"enabled":       {Type: schema.TypeBool, Computed: true},
		"status":        {Type: schema.TypeString, Computed: true},
		"created_at":    {Type: schema.TypeString, Computed: true},
		"updated_at":    {Type: schema.TypeString, Computed: true},
	}
}

func flattenVpcFlowLog(flowLog flow_logs.FlowLog) map[string]interface{} {
	return map[string]interface{}{
		"id":            flowLog.ID,
		"name":          flowLog.Name,
		"tenant_id":     flowLog.TenantID,
		"description":   flowLog.Description,
		"resource_type": flowLog.ResourceType,
		"resource_id":   flowLog.ResourceID,
		"traffic_type":  flowLog.TrafficType,
		"log_group_id":  flowLog.LogGroupID,
		"log_topic_id":  flowLog.LogTopicID,
		"enabled":       flowLog.AdminState,
		"status":        flowLog.Status,
		"created_at":    flowLog.CreatedAt,
		"updated_at":    flowLog.UpdatedAt,
	}
}
