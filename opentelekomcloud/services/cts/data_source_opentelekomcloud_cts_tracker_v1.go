package cts

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cts/v1/tracker"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCTSTrackerV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCTSTrackerV1Read,

		Schema: map[string]*schema.Schema{
			"tracker_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"system",
				}, false),
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_prefix_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_lts_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"log_topic_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCTSTrackerV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	ctsClient, err := config.CtsV1Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}

	ctsTracker, err := tracker.Get(ctsClient, trackerName)
	if err != nil {
		return fmterr.Errorf("unable to retrieve cts tracker: %s", err)
	}

	log.Printf("[INFO] Retrieved cts tracker %s", ctsTracker.TrackerName)

	d.SetId(ctsTracker.TrackerName)

	mErr := multierror.Append(
		d.Set("tracker_name", ctsTracker.TrackerName),
		d.Set("bucket_name", ctsTracker.BucketName),
		d.Set("file_prefix_name", ctsTracker.FilePrefixName),
		d.Set("status", ctsTracker.Status),
		d.Set("is_lts_enabled", ctsTracker.Lts.IsLtsEnabled),
		d.Set("log_topic_name", ctsTracker.Lts.LogTopicName),
		d.Set("log_group_name", ctsTracker.Lts.LogGroupName),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
