package cbr

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/backups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCBRBackupsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCBRBackupsV3Read,

		Schema: map[string]*schema.Schema{
			"checkpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expired_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_az": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_trigger": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"bootable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"incremental": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"support_lld": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"supported_restore_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"contain_system_disk": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"system_disk": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceCBRBackupsV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	backupClient, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating CbrV3Client: %w", err)
	}

	listOpts := backups.ListOpts{
		ID:           d.Get("id").(string),
		CheckpointID: d.Get("checkpoint_id").(string),
		Status:       d.Get("status").(string),
		ResourceName: d.Get("resource_name").(string),
		ImageType:    d.Get("image_type").(string),
		ResourceType: d.Get("resource_type").(string),
		ResourceID:   d.Get("resource_id").(string),
		Name:         d.Get("name").(string),
		ParentID:     d.Get("parent_id").(string),
		ResourceAZ:   d.Get("resource_az").(string),
		VaultID:      d.Get("vault_id").(string),
	}

	extractedBackups, err := backups.List(backupClient, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to list all backups pages: %s", err)
	}

	if len(extractedBackups) < 1 {
		return common.DataSourceTooFewDiag
	}

	if len(extractedBackups) > 1 {
		return common.DataSourceTooManyDiag
	}

	backup := extractedBackups[0]

	log.Printf("[INFO] Retrieved backup policy %s using given filter", backup.ID)

	d.SetId(backup.ID)
	mErr := multierror.Append(
		d.Set("checkpoint_id", backup.CheckpointID),
		d.Set("created_at", backup.CreatedAt),
		d.Set("description", backup.Description),
		d.Set("expired_at", backup.ExpiredAt),
		d.Set("image_type", backup.ImageType),
		d.Set("name", backup.Name),
		d.Set("parent_id", backup.ParentID),
		d.Set("project_id", backup.ProjectID),
		d.Set("provider_id", backup.ProviderID),
		d.Set("resource_az", backup.ResourceAZ),
		d.Set("resource_id", backup.ResourceID),
		d.Set("resource_name", backup.ResourceName),
		d.Set("resource_size", backup.ResourceSize),
		d.Set("resource_type", backup.ResourceType),
		d.Set("status", backup.Status),
		d.Set("updated_at", backup.UpdatedAt),
		d.Set("vault_id", backup.VaultId),
		d.Set("auto_trigger", backup.ExtendInfo.AutoTrigger),
		d.Set("bootable", backup.ExtendInfo.Bootable),
		d.Set("incremental", backup.ExtendInfo.Incremental),
		d.Set("snapshot_id", backup.ExtendInfo.SnapshotID),
		d.Set("support_lld", backup.ExtendInfo.SupportLld),
		d.Set("supported_restore_mode", backup.ExtendInfo.SupportedRestoreMode),
		d.Set("contain_system_disk", backup.ExtendInfo.ContainSystemDisk),
		d.Set("encrypted", backup.ExtendInfo.Encrypted),
		d.Set("system_disk", backup.ExtendInfo.SystemDisk),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
