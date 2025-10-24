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

func DataSourceTaurusDBV3MysqlSlowLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBMysqlSlowLogsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_date": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_date": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slow_log_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rows_sent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rows_examined": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"users": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_sample": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_ip": {
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

func dataSourceTaurusDBMysqlSlowLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	var mErr *multierror.Error

	opts := logs.GetSlowLogsOpts{
		InstanceId: d.Get("instance_id").(string),
		StartDate:  d.Get("start_date").(string),
		EndDate:    d.Get("end_date").(string),
		NodeId:     d.Get("node_id").(string),
	}

	if v, ok := d.GetOk("type"); ok {
		opts.Type = v.(string)
	}

	slowLogs, err := logs.GetSlowLogs(client, opts)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB MySQL slow logs: %s", err)
	}

	result := flattenTaurusDBMysqlSlowLogs(slowLogs)

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", config.GetRegion(d)),
		d.Set("slow_log_list", result),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTaurusDBMysqlSlowLogs(slowLogs []logs.MysqlSlowLogList) []map[string]interface{} {
	if len(slowLogs) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(slowLogs))
	for _, slowLog := range slowLogs {
		result = append(result, map[string]interface{}{
			"node_id":       slowLog.NodeId,
			"count":         slowLog.Count,
			"time":          slowLog.Time,
			"lock_time":     slowLog.LockTime,
			"rows_sent":     slowLog.RowsSent,
			"rows_examined": slowLog.RowsExamined,
			"database":      slowLog.Database,
			"users":         slowLog.Users,
			"query_sample":  slowLog.QuerySample,
			"type":          slowLog.Type,
			"start_time":    slowLog.StartTime,
			"client_ip":     slowLog.ClientIp,
		})
	}
	return result
}
