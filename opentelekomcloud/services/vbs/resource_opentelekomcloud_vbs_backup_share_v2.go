package vbs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vbs/v2/shares"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVBSBackupShareV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVBSBackupShareV2Create,
		ReadContext:   resourceVBSBackupShareV2Read,
		DeleteContext: resourceVBSBackupShareV2Delete,
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
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"to_project_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_metadata": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"share_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceBackupShareToProjectIdsV2(d *schema.ResourceData) []string {
	rawProjectIDs := d.Get("to_project_ids").(*schema.Set)
	projectids := make([]string, rawProjectIDs.Len())
	for i, raw := range rawProjectIDs.List() {
		projectids[i] = raw.(string)
	}
	return projectids
}

func resourceVBSBackupShareV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vbsClient, err := config.VbsV2Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud vbs client: %s", err)
	}

	createOpts := shares.CreateOpts{
		ToProjectIDs: resourceBackupShareToProjectIdsV2(d),
		BackupID:     d.Get("backup_id").(string),
	}

	n, err := shares.Create(vbsClient, createOpts).Extract()

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VBS Backup Share: %s", err)
	}

	share := n[0]
	d.SetId(share.BackupID)

	log.Printf("[INFO] VBS Backup Share ID: %s", d.Id())

	return resourceVBSBackupShareV2Read(ctx, d, meta)
}

func resourceVBSBackupShareV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vbsClient, err := config.VbsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Vbs client: %s", err)
	}

	backups, err := shares.List(vbsClient, shares.ListOpts{BackupID: d.Id()})
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Vbs: %s", err)
	}

	n := backups[0]

	mErr := multierror.Append(
		d.Set("backup_id", n.BackupID),
		d.Set("backup_name", n.Backup.Name),
		d.Set("backup_status", n.Backup.Status),
		d.Set("description", n.Backup.Description),
		d.Set("availability_zone", n.Backup.AvailabilityZone),
		d.Set("volume_id", n.Backup.VolumeID),
		d.Set("size", n.Backup.Size),
		d.Set("service_metadata", n.Backup.ServiceMetadata),
		d.Set("container", n.Backup.Container),
		d.Set("snapshot_id", n.Backup.SnapshotID),
		d.Set("region", config.GetRegion(d)),
		d.Set("to_project_ids", resourceToProjectIdsV2(backups)),
		d.Set("share_ids", resourceShareIDsV2(backups)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVBSBackupShareV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vbsClient, err := config.VbsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Vbs client: %s", err)
	}

	deleteopts := shares.DeleteOpts{IsBackupID: true}

	err = shares.Delete(vbsClient, d.Id(), deleteopts).ExtractErr()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[INFO] Successfully deleted OpenTelekomCloud Vbs Backup Share %s", d.Id())
		}
		if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
			if errCode.Actual == 409 {
				log.Printf("[INFO] Error deleting OpenTelekomCloud Vbs Backup Share %s", d.Id())
			}
		}
		log.Printf("[INFO] Successfully deleted OpenTelekomCloud Vbs Backup Share %s", d.Id())
	}

	d.SetId("")
	return nil
}

func resourceToProjectIdsV2(s []shares.Share) []string {
	projectids := make([]string, len(s))
	for i, raw := range s {
		projectids[i] = raw.ToProjectID
	}
	return projectids
}

func resourceShareIDsV2(s []shares.Share) []string {
	shareids := make([]string, len(s))
	for i, raw := range s {
		shareids[i] = raw.ID
	}
	return shareids
}
