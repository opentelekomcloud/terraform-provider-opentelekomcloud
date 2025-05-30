package lts

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/streams"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceLtsStreamsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLtsStreamsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"streams": {
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
						"ttl_in_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLtsStreamsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := streams.ListStreamsOpts{
		GroupName:  d.Get("log_group_name").(string),
		StreamName: d.Get("name").(string),
	}
	requestResp, err := streams.ListStreams(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	var ltsStreams []map[string]interface{}
	for _, s := range requestResp {
		ltsStream := map[string]interface{}{
			"id":                    s.LogStreamId,
			"name":                  s.LogStreamName,
			"ttl_in_days":           s.TTLInDays,
			"tags":                  ignoreSysEpsTag(s.Tag),
			"created_at":            common.FormatTimeStampRFC3339(s.CreationTime/1000, false),
			"enterprise_project_id": s.Tag["_sys_enterprise_project_id"],
		}
		ltsStreams = append(ltsStreams, ltsStream)
	}

	mErr := multierror.Append(
		d.Set("streams", ltsStreams),
		d.Set("region", config.GetRegion(d)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}
