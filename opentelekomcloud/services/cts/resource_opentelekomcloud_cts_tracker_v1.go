package cts

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cts/v1/tracker"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCTSTrackerV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCTSTrackerCreate,
		ReadContext:   resourceCTSTrackerRead,
		UpdateContext: resourceCTSTrackerUpdate,
		DeleteContext: resourceCTSTrackerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tracker_name": {
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
			"is_support_smn": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"topic_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operations": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"is_send_all_key_operation": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"need_notify_user_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}

}

func resourceCTSTrackerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CtsV1Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}

	createOpts := tracker.CreateOptsWithSMN{
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		SimpleMessageNotification: tracker.SimpleMessageNotification{
			IsSupportSMN:          d.Get("is_support_smn").(bool),
			TopicID:               d.Get("topic_id").(string),
			Operations:            resourceCTSOperations(d),
			IsSendAllKeyOperation: d.Get("is_send_all_key_operation").(bool),
			NeedNotifyUserList:    resourceCTSNeedNotifyUserList(d),
		},
	}

	trackers, err := tracker.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating CTS tracker: %w", err)
	}

	d.SetId(trackers.TrackerName)

	stateConf := &resource.StateChangeConf{
		Refresh: refreshCTSTrackerState(d, client),
		Target:  []string{"exists"},
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for CTS tracker to appear in the list: %w", err)
	}

	return resourceCTSTrackerRead(ctx, d, meta)
}

func refreshCTSTrackerState(d *schema.ResourceData, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	listOpts := tracker.ListOpts{
		TrackerName:    d.Get("tracker_name").(string),
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		Status:         d.Get("status").(string),
	}
	return func() (interface{}, string, error) {
		trackers, err := tracker.List(client, listOpts)
		if err != nil {
			return nil, "", err
		}
		if len(trackers) == 0 {
			return trackers, "", nil
		}
		return trackers, "exists", nil
	}
}

func resourceCTSTrackerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	ctsClient, err := config.CtsV1Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}

	listOpts := tracker.ListOpts{
		TrackerName:    d.Get("tracker_name").(string),
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		Status:         d.Get("status").(string),
	}
	trackers, err := tracker.List(ctsClient, listOpts)
	if err != nil {
		return fmterr.Errorf("error retrieving cts tracker: %w", err)
	}

	if len(trackers) == 0 {
		log.Printf("[WARN] Removing cts tracker %s as it's already gone", d.Id())
		d.SetId("")
		return nil
	}

	ctsTracker := trackers[0]

	mErr := multierror.Append(
		d.Set("tracker_name", ctsTracker.TrackerName),
		d.Set("bucket_name", ctsTracker.BucketName),
		d.Set("status", ctsTracker.Status),
		d.Set("file_prefix_name", ctsTracker.FilePrefixName),
		d.Set("is_support_smn", ctsTracker.SimpleMessageNotification.IsSupportSMN),
		d.Set("topic_id", ctsTracker.SimpleMessageNotification.TopicID),
		d.Set("is_send_all_key_operation", ctsTracker.SimpleMessageNotification.IsSendAllKeyOperation),
		d.Set("operations", ctsTracker.SimpleMessageNotification.Operations),
		d.Set("need_notify_user_list", ctsTracker.SimpleMessageNotification.NeedNotifyUserList),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting CTS tracker fields: %w", err)
	}

	return nil
}

func resourceCTSTrackerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	ctsClient, err := config.CtsV1Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}
	updateOpts := tracker.UpdateOptsWithSMN{
		BucketName: d.Get("bucket_name").(string),
		SimpleMessageNotification: tracker.SimpleMessageNotification{
			IsSupportSMN:          d.Get("is_support_smn").(bool),
			TopicID:               d.Get("topic_id").(string),
			Operations:            resourceCTSOperations(d),
			IsSendAllKeyOperation: d.Get("is_send_all_key_operation").(bool),
			NeedNotifyUserList:    resourceCTSNeedNotifyUserList(d),
		},
	}

	if d.HasChange("file_prefix_name") {
		updateOpts.FilePrefixName = d.Get("file_prefix_name").(string)
	}
	if d.HasChange("status") {
		updateOpts.Status = d.Get("status").(string)
	}

	_, err = tracker.Update(ctsClient, updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating cts tracker: %w", err)
	}
	return resourceCTSTrackerRead(ctx, d, meta)
}

func resourceCTSTrackerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CtsV1Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(clientError, err)
	}

	if err := tracker.Delete(client).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Successfully deleted cts tracker %s", d.Id())

	stateConf := &resource.StateChangeConf{
		Refresh: refreshCTSTrackerState(d, client),
		Target:  []string{""},
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for CTS tracker to appear in the list: %w", err)
	}
	d.SetId("")

	return nil
}

func resourceCTSOperations(d *schema.ResourceData) []string {
	rawOperations := d.Get("operations").(*schema.Set)
	operation := make([]string, (rawOperations).Len())
	for i, raw := range rawOperations.List() {
		operation[i] = raw.(string)
	}
	return operation
}

func resourceCTSNeedNotifyUserList(d *schema.ResourceData) []string {
	rawNotify := d.Get("need_notify_user_list").(*schema.Set)
	notify := make([]string, (rawNotify).Len())
	for i, raw := range rawNotify.List() {
		notify[i] = raw.(string)
	}
	return notify
}
