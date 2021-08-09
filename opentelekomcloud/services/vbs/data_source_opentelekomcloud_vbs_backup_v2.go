package vbs

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vbs/v2/backups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vbs/v2/shares"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVBSBackupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVBSBackupV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"container": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_metadata": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to_project_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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

func dataSourceVBSBackupV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vbsClient, err := config.VbsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Vbs client: %s", err)
	}

	listBackupOpts := backups.ListOpts{
		Id:         d.Id(),
		Name:       d.Get("name").(string),
		Status:     d.Get("status").(string),
		VolumeId:   d.Get("volume_id").(string),
		SnapshotId: d.Get("snapshot_id").(string),
	}

	refinedBackups, err := backups.List(vbsClient, listBackupOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve backups: %s", err)
	}

	if len(refinedBackups) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedBackups) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Backup := refinedBackups[0]

	log.Printf("[INFO] Retrieved Backup using given filter %s: %+v", Backup.Id, Backup)
	d.SetId(Backup.Id)

	listShareOpts := shares.ListOpts{
		BackupID: d.Id(),
	}

	shareList, err := shares.List(vbsClient, listShareOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve shares: %s", err)
	}

	mErr := multierror.Append(
		d.Set("name", Backup.Name),
		d.Set("description", Backup.Description),
		d.Set("status", Backup.Status),
		d.Set("availability_zone", Backup.AvailabilityZone),
		d.Set("snapshot_id", Backup.SnapshotId),
		d.Set("service_metadata", Backup.ServiceMetadata),
		d.Set("size", Backup.Size),
		d.Set("container", Backup.Container),
		d.Set("volume_id", Backup.VolumeId),
		d.Set("region", config.GetRegion(d)),
		d.Set("to_project_ids", resourceToProjectIdsV2(shareList)),
		d.Set("share_ids", resourceShareIDsV2(shareList)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
