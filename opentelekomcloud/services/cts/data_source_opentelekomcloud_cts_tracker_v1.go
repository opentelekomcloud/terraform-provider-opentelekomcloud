package cts

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cts/v1/tracker"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCTSTrackerV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCTSTrackerV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"file_prefix_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tracker_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_support_smn": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"topic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"is_send_all_key_operation": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"need_notify_user_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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

	listOpts := tracker.Tracker{
		TrackerName:    d.Get("tracker_name").(string),
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		Status:         d.Get("status").(string),
	}

	refinedTrackers, err := tracker.Get(ctsClient)
	if err != nil {
		return fmterr.Errorf("unable to retrieve cts tracker: %s", err)
	}

	if len(refinedTrackers) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedTrackers) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	trackers := refinedTrackers[0]

	log.Printf("[INFO] Retrieved cts tracker %s using given filter", trackers.TrackerName)

	d.SetId(trackers.TrackerName)

	mErr := multierror.Append(
		d.Set("tracker_name", trackers.TrackerName),
		d.Set("bucket_name", trackers.BucketName),
		d.Set("file_prefix_name", trackers.FilePrefixName),
		d.Set("status", trackers.Status),
		d.Set("is_support_smn", trackers.SimpleMessageNotification.IsSupportSMN),
		d.Set("topic_id", trackers.SimpleMessageNotification.TopicID),
		d.Set("is_send_all_key_operation", trackers.SimpleMessageNotification.IsSendAllKeyOperation),
		d.Set("operations", trackers.SimpleMessageNotification.Operations),
		d.Set("need_notify_user_list", trackers.SimpleMessageNotification.NeedNotifyUserList),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
