package csbs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/csbs/v1/backup"
	res "github.com/opentelekomcloud/gophertelekomcloud/openstack/csbs/v1/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCSBSBackupV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCSBSBackupV1Create,
		ReadContext:   resourceCSBSBackupV1Read,
		DeleteContext: resourceCSBSBackupV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"backup_record_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "OS::Nova::Server",
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_backups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"space_saving_ratio": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"average_speed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incremental": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_volume_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"vm_metadata": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCSBSBackupV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating csbs client: %s", err)
	}

	resourceID := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)

	query, err := res.GetResBackupCapabilities(client, []res.ResourceBackupCapOpts{
		{
			ResourceId:   resourceID,
			ResourceType: resourceType,
		},
	})
	if err != nil {
		return fmterr.Errorf("error querying resource backup capability: %s", err)
	}

	if query[0].Result {
		createOpts := backup.CreateOpts{
			BackupName:   d.Get("backup_name").(string),
			Description:  d.Get("description").(string),
			ResourceType: resourceType,
			Tags:         resourceCSBSTagsV1(d),
		}

		checkpoint, err := backup.Create(client, resourceID, createOpts)
		if err != nil {
			return fmterr.Errorf("error creating backup: %s", err)
		}

		backupOpts := backup.ListOpts{CheckpointId: checkpoint.Id}
		backupItems, err := backup.List(client, backupOpts)
		if err != nil {
			return fmterr.Errorf("error listing Backup: %s", err)
		}

		if len(backupItems) == 0 {
			return fmterr.Errorf("not able to find created Backup: %s", err)
		}

		backupObject := backupItems[0]

		d.SetId(backupObject.Id)

		log.Printf("[INFO] Resource Backup %s created successfully", backupObject.Id)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"protecting"},
			Target:     []string{"available"},
			Refresh:    waitForCSBSBackupActive(client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      3 * time.Minute,
			MinTimeout: 3 * time.Minute,
		}
		_, stateErr := stateConf.WaitForStateContext(ctx)
		if stateErr != nil {
			return fmterr.Errorf(
				"Error waiting for Backup (%s) to become available: %s",
				backupObject.Id, stateErr)
		}
	} else {
		return fmterr.Errorf("error code: %s\n Error msg: %s", query[0].ErrorCode, query[0].ErrorMsg)
	}

	return resourceCSBSBackupV1Read(ctx, d, meta)
}

func resourceCSBSBackupV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating csbs client: %s", err)
	}

	backupObject, err := backup.Get(client, d.Id())

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Removing backup %s as it's already gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmterr.Errorf("error retrieving backup: %s", err)
	}

	mErr := multierror.Append(
		d.Set("resource_id", backupObject.ResourceId),
		d.Set("backup_name", backupObject.Name),
		d.Set("description", backupObject.Description),
		d.Set("resource_type", backupObject.ResourceType),
		d.Set("status", backupObject.Status),
		d.Set("volume_backups", flattenCSBSVolumeBackups(backupObject)),
		d.Set("vm_metadata", flattenCSBSVMMetadata(backupObject)),
		d.Set("backup_record_id", backupObject.CheckpointId),
		d.Set("region", config.GetRegion(d)),
		d.Set("tags", flattenCSBSTags(backupObject.Tags)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCSBSBackupV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating csbs client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForCSBSBackupDelete(client, d.Id(), d.Get("backup_record_id").(string)),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting csbs backup: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCSBSBackupActive(client *golangsdk.ServiceClient, backupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := backup.Get(client, backupId)
		if err != nil {
			return nil, "", err
		}

		if n.Id == "error" {
			return nil, "", fmt.Errorf("backup status: %s", n.Status)
		}

		return n, n.Status, nil
	}
}

func waitForCSBSBackupDelete(client *golangsdk.ServiceClient, backupId string, backupRecordID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := backup.Get(client, backupId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted csbs backup %s", backupId)
				return r, "deleted", nil
			}
			return r, "deleting", err
		}

		err = backup.Delete(client, backupRecordID)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud Backup %s", backupId)
				return r, "deleted", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "deleting", nil
				}
			}
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				return r, "deleting", nil
			}

			return r, "deleting", err
		}

		return r, r.Status, nil
	}
}

func flattenCSBSVolumeBackups(backupObject *backup.Backup) []map[string]interface{} {
	var volumeBackups []map[string]interface{}

	for _, volume := range backupObject.ExtendInfo.VolumeBackups {
		mapping := map[string]interface{}{
			"status":             volume.Status,
			"space_saving_ratio": volume.SpaceSavingRatio,
			"name":               volume.Name,
			"bootable":           volume.Bootable,
			"average_speed":      volume.AverageSpeed,
			"source_volume_size": volume.SourceVolumeSize,
			"source_volume_id":   volume.SourceVolumeId,
			"snapshot_id":        volume.SnapshotID,
			"incremental":        volume.Incremental,
			"source_volume_name": volume.SourceVolumeName,
			"image_type":         volume.ImageType,
			"id":                 volume.Id,
			"size":               volume.Size,
		}
		volumeBackups = append(volumeBackups, mapping)
	}

	return volumeBackups
}

func flattenCSBSVMMetadata(backupObject *backup.Backup) []map[string]interface{} {
	var vmMetadata []map[string]interface{}

	mapping := map[string]interface{}{
		"name":               backupObject.ExtendInfo.ResourceName,
		"eip":                backupObject.VMMetadata.Eip,
		"cloud_service_type": backupObject.VMMetadata.CloudServiceType,
		"ram":                backupObject.VMMetadata.Ram,
		"vcpus":              backupObject.VMMetadata.Vcpus,
		"private_ip":         backupObject.VMMetadata.PrivateIp,
		"disk":               backupObject.VMMetadata.Disk,
		"image_type":         backupObject.VMMetadata.ImageType,
	}

	vmMetadata = append(vmMetadata, mapping)

	return vmMetadata
}
