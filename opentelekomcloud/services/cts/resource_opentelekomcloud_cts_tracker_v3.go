package cts

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
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
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled", "disabled",
				}, false),
			},
			"file_prefix_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: common.ValidateName,
			},
			"is_obs_created": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
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
			"detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_lifecycle": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceCTSTrackerV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CtsV3Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}

	createOpts := tracker.CreateOpts{
		IsLtsEnabled: d.Get("is_lts_enabled").(bool),
		TrackerName:  trackerName,
		TrackerType:  trackerName,
		ObsInfo: tracker.ObsInfo{
			BucketName:     d.Get("bucket_name").(string),
			FilePrefixName: d.Get("file_prefix_name").(string),
			IsObsCreated:   pointerto.Bool(d.Get("is_obs_created").(bool)),
		},
	}

	ctsTracker, err := tracker.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating CTS tracker: %w", err)
	}

	d.SetId(ctsTracker.TrackerName)

	if d.Get("status").(string) == "disabled" {
		_, err = tracker.Update(client, tracker.UpdateOpts{
			TrackerType: trackerName,
			TrackerName: trackerName,
			Status:      "disabled",
		})
		if err != nil {
			return fmterr.Errorf("error setting CTS tracker status: %w", err)
		}
	}

	return resourceCTSTrackerV3Read(ctx, d, meta)
}

func resourceCTSTrackerV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CtsV3Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}

	ctsTracker, err := tracker.List(client, trackerName)
	if err != nil {
		return fmterr.Errorf("error retrieving cts tracker: %w", err)
	}

	if len(ctsTracker) > 1 {
		return fmterr.Errorf("tracker query returned 2 or more results")
	}

	mErr := multierror.Append(
		d.Set("id", ctsTracker[0].Id),
		d.Set("tracker_name", ctsTracker[0].TrackerName),
		d.Set("tracker_type", ctsTracker[0].TrackerType),
		d.Set("domain_id", ctsTracker[0].DomainId),
		d.Set("project_id", ctsTracker[0].ProjectId),
		d.Set("status", ctsTracker[0].Status),
		d.Set("detail", ctsTracker[0].Detail),
		d.Set("log_group_name", ctsTracker[0].Lts.LogGroupName),
		d.Set("log_topic_name", ctsTracker[0].Lts.LogTopicName),
		d.Set("is_lts_enabled", ctsTracker[0].Lts.IsLtsEnabled),
		d.Set("bucket_name", ctsTracker[0].ObsInfo.BucketName),
		d.Set("file_prefix_name", ctsTracker[0].ObsInfo.FilePrefixName),
		d.Set("is_obs_created", ctsTracker[0].ObsInfo.IsObsCreated),
		d.Set("bucket_lifecycle", ctsTracker[0].ObsInfo.BucketLifecycle),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting CTS tracker fields: %w", err)
	}

	return nil
}

func resourceCTSTrackerV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CtsV3Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}
	updateOpts := tracker.UpdateOpts{
		TrackerType: trackerName,
		TrackerName: trackerName,
	}

	if d.HasChange("status") {
		updateOpts.Status = d.Get("status").(string)
		if updateOpts.Status == "enabled" {
			// tracker needs to be enabled for other parameters to be applied
			_, err = tracker.Update(client, updateOpts)
			if err != nil {
				return fmterr.Errorf("error updating CTS tracker: %w", err)
			}
		}
	}

	if d.HasChange("is_lts_enabled") {
		updateOpts.IsLtsEnabled = pointerto.Bool(d.Get("is_lts_enabled").(bool))
	}

	if d.HasChange("bucket_name") {
		updateOpts.ObsInfo.BucketName = d.Get("bucket_name").(string)
	}

	if d.HasChange("file_prefix_name") {
		updateOpts.ObsInfo.FilePrefixName = d.Get("file_prefix_name").(string)
		updateOpts.ObsInfo.BucketName = d.Get("bucket_name").(string)
	}

	if d.HasChange("is_obs_created") {
		updateOpts.ObsInfo.IsObsCreated = pointerto.Bool(d.Get("is_obs_created").(bool))
	}

	_, err = tracker.Update(client, updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating CTS tracker: %w", err)
	}

	return resourceCTSTrackerV3Read(ctx, d, meta)
}

func resourceCTSTrackerV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CtsV3Client(config.GetProjectName(d))
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
