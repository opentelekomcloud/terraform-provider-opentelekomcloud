package er

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fl "github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/flow-logs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceErFlowLogsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceErFlowLogsV3Read,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_log_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_logs": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_stream_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_store_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
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

func dataSourceErFlowLogsV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ErV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	resp, err := fl.List(client, d.Get("instance_id").(string), fl.ListOpts{
		ResourceType: d.Get("resource_type").(string),
		ResourceID:   common.StringSliceIgnoreEmpty(d.Get("resource_id").(string)),
		SortKey:      []string{"name"},
	})
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud ER v3 flow logs: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	filtered := filterListFlowLogsResponseBody(d, resp)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", config.GetRegion(d)),
		d.Set("flow_logs", flattenListTransitIpsResponseBody(filtered)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListTransitIpsResponseBody(logs []fl.FlowLogResponse) []map[string]interface{} {
	if len(logs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(logs))
	for i, item := range logs {
		result[i] = map[string]interface{}{
			"id":             item.ID,
			"name":           item.Name,
			"description":    item.Description,
			"resource_type":  item.ResourceType,
			"resource_id":    item.ResourceId,
			"log_group_id":   item.LogGroupId,
			"log_stream_id":  item.LogStreamId,
			"log_store_type": item.LogStoreType,
			"created_at":     common.FormatTimeStampRFC3339(common.ConvertTimeStrToNanoTimestamp(item.CreatedAt)/1000, false),
			"updated_at":     common.FormatTimeStampRFC3339(common.ConvertTimeStrToNanoTimestamp(item.UpdatedAt)/1000, false),
			"status":         item.Status,
			"enabled":        item.Enabled,
		}
	}
	return result
}

func filterListFlowLogsResponseBody(d *schema.ResourceData, logs []fl.FlowLogResponse) []fl.FlowLogResponse {
	rst := make([]fl.FlowLogResponse, 0, len(logs))

	for _, v := range logs {
		if param, ok := d.GetOk("flow_log_id"); ok &&
			param != v.ID {
			continue
		}
		if param, ok := d.GetOk("name"); ok &&
			param != v.Name {
			continue
		}
		if param, ok := d.GetOk("status"); ok &&
			param != v.Status {
			continue
		}
		if param, ok := d.GetOk("enabled"); ok &&
			param != v.Enabled {
			continue
		}
		if param, ok := d.GetOk("log_group_id"); ok &&
			param != v.LogGroupId {
			continue
		}
		if param, ok := d.GetOk("log_stream_id"); ok &&
			param != v.LogStreamId {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}
