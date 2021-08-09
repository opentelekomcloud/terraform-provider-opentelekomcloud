package vbs

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vbs/v2/backups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVBSBackupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVBSBackupV2Create,
		ReadContext:   resourceVBSBackupV2Read,
		DeleteContext: resourceVBSBackupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateVBSBackupName,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 64),
					validation.StringMatch(
						regexp.MustCompile(`^[^<>]+$`),
						"description doesn't comply with restrictions",
					),
				),
			},
			"container": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_metadata": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: common.ValidateVBSTagKey,
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: common.ValidateVBSTagValue,
						},
					},
				},
			},
		},
	}
}

func resourceVBSBackupTagsV2(d *schema.ResourceData) []backups.Tag {
	rawTags := d.Get("tags").(*schema.Set).List()
	tags := make([]backups.Tag, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tags[i] = backups.Tag{
			Key:   rawMap["key"].(string),
			Value: rawMap["value"].(string),
		}
	}
	return tags
}

func resourceVBSBackupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vbsClient, err := config.VbsV2Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud vbs client: %s", err)
	}

	createOpts := backups.CreateOpts{
		Name:        d.Get("name").(string),
		VolumeId:    d.Get("volume_id").(string),
		SnapshotId:  d.Get("snapshot_id").(string),
		Description: d.Get("description").(string),
		Tags:        resourceVBSBackupTagsV2(d),
	}

	n, err := backups.Create(vbsClient, createOpts).ExtractJobResponse()

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VBS Backup: %s", err)
	}

	if err := backups.WaitForJobSuccess(vbsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return diag.FromErr(err)
	}

	entity, err := backups.GetJobEntity(vbsClient, n.JobID, "backup_id")
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := entity.(string); ok {
		d.SetId(id)
		return resourceVBSBackupV2Read(ctx, d, meta)
	}

	return fmterr.Errorf("unexpected conversion error in resourceVBSBackupV2Create.")
}

func resourceVBSBackupV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vbsClient, err := config.VbsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating Vbs client: %s", err)
	}

	n, err := backups.Get(vbsClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving VBS Backup: %s", err)
	}

	mErr := multierror.Append(
		d.Set("name", n.Name),
		d.Set("description", n.Description),
		d.Set("status", n.Status),
		d.Set("availability_zone", n.AvailabilityZone),
		d.Set("snapshot_id", n.SnapshotId),
		d.Set("service_metadata", n.ServiceMetadata),
		d.Set("size", n.Size),
		d.Set("container", n.Container),
		d.Set("volume_id", n.VolumeId),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceVBSBackupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vbsClient, err := config.VbsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating  vbs: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForBackupDelete(vbsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting VBS Backup: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForBackupDelete(client *golangsdk.ServiceClient, backupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := backups.Get(client, backupID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return r, "deleted", nil
			}
			return nil, "available", err
		}

		if r.Status != "deleting" {
			err := backups.Delete(client, backupID).ExtractErr()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					log.Printf("[INFO] Successfully deleted VBS backup %s", backupID)
					return r, "deleted", nil
				}
				return r, r.Status, err
			}
		}
		return r, r.Status, nil
	}
}
