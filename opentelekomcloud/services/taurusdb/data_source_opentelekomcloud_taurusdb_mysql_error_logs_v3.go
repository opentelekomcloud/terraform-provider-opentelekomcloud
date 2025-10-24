package taurusdb

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/logs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceTaurusDBV3MysqlErrorLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBMysqlErrorLogsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"error_log_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content": {
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

func dataSourceTaurusDBMysqlErrorLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	opts := logs.GetErrorLogsOpts{
		InstanceId: d.Get("instance_id").(string),
		StartDate:  d.Get("start_time").(string),
		EndDate:    d.Get("end_time").(string),
		NodeId:     d.Get("node_id").(string),
	}

	if v, ok := d.GetOk("level"); ok {
		opts.Level = v.(string)
	}

	errorLogs, err := logs.GetErrorLogs(client, opts)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB MySQL error logs: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	result := make([]map[string]interface{}, len(errorLogs))
	for i, errorLog := range errorLogs {
		result[i] = map[string]interface{}{
			"node_id": errorLog.NodeId,
			"time":    errorLog.Time,
			"level":   errorLog.Level,
			"content": errorLog.Content,
		}
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("error_log_list", result),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
