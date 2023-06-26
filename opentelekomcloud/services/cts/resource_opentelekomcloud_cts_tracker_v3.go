package cts

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cts/v3/tracker"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCTSTrackerV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCTSTrackerV3Create,
		ReadContext:   resourceCTSTrackerV3Read,
		UpdateContext: resourceCTSTrackerV3Update,
		DeleteContext: resourceCTSTrackerV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"is_lts_enabled": {
				Type:     schema.Bool,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_prefix_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: common.ValidateName,
			},
			"is_lts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_topic_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tracker_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tracker_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCTSTrackerV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CtsV1Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}

	ltsEnabled := d.Get("is_lts_enabled").(bool)

	createOpts := tracker.CreateOpts{
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		Lts: tracker.CreateLts{
			IsLtsEnabled: &ltsEnabled,
		},
	}

	ctsTracker, err := tracker.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating CTS tracker: %w", err)
	}

	d.SetId(ctsTracker.TrackerName)

	return resourceCTSTrackerRead(ctx, d, meta)
}

func resourceCTSTrackerV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CtsV1Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}

	ctsTracker, err := tracker.Get(client, trackerName)
	if err != nil {
		return fmterr.Errorf("error retrieving cts tracker: %w", err)
	}

	mErr := multierror.Append(
		d.Set("tracker_name", ctsTracker.TrackerName),
		d.Set("bucket_name", ctsTracker.BucketName),
		d.Set("status", ctsTracker.Status),
		d.Set("file_prefix_name", ctsTracker.FilePrefixName),
		d.Set("is_lts_enabled", ctsTracker.Lts.IsLtsEnabled),
		d.Set("log_group_name", ctsTracker.Lts.LogGroupName),
		d.Set("log_topic_name", ctsTracker.Lts.LogTopicName),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting CTS tracker fields: %w", err)
	}

	return nil
}

func resourceCTSTrackerV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CtsV1Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}
	ltsEnabled := d.Get("is_lts_enabled").(bool)
	updateOpts := tracker.UpdateOpts{
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		Lts: tracker.CreateLts{
			IsLtsEnabled: &ltsEnabled,
		},
	}

	if d.HasChange("file_prefix_name") {
		updateOpts.FilePrefixName = d.Get("file_prefix_name").(string)
	}
	if d.HasChange("status") {
		updateOpts.Status = d.Get("status").(string)
	}

	_, err = tracker.Update(client, updateOpts, trackerName)
	if err != nil {
		return fmterr.Errorf("error updating CTS tracker: %w", err)
	}
	return resourceCTSTrackerRead(ctx, d, meta)
}

func resourceCTSTrackerV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CtsV1Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}

	if err := tracker.Delete(client, trackerName); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Successfully deleted cts tracker %s", d.Id())

	d.SetId("")

	return nil
}
