package csbs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/csbs/v1/backup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceCSBSBackupV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCSBSBackupV1Create,
		ReadContext:   resourceCSBSBackupV1Read,
		DeleteContext: resourceCSBSBackupV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
							Type:     schema.TypeInt,
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
	backupClient, err := config.CsbsV1Client(config.GetRegion(d))

	if err != nil {
		return diag.Errorf("Error creating csbs client: %s", err)
	}

	resourceID := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)

	queryOpts := backup.ResourceBackupCapOpts{
		CheckProtectable: []backup.ResourceCapQueryParams{
			{
				ResourceId:   resourceID,
				ResourceType: resourceType,
			},
		},
	}

	query, err := backup.QueryResourceBackupCapability(backupClient, queryOpts).ExtractQueryResponse()
	if err != nil {
		return diag.Errorf("Error querying resource backup capability: %s", err)
	}

	if query[0].Result {

		createOpts := backup.CreateOpts{
			BackupName:   d.Get("backup_name").(string),
			Description:  d.Get("description").(string),
			ResourceType: resourceType,
			Tags:         resourceCSBSTagsV1(d),
		}

		checkpoint, err := backup.Create(backupClient, resourceID, createOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating backup: %s", err)
		}

		backupOpts := backup.ListOpts{CheckpointId: checkpoint.Id}
		backupItems, err := backup.List(backupClient, backupOpts)

		if err != nil {
			return diag.Errorf("Error listing Backup: %s", err)
		}

		if len(backupItems) == 0 {
			return diag.Errorf("Not able to find created Backup: %s", err)
		}

		backupObject := backupItems[0]

		d.SetId(backupObject.Id)

		log.Printf("[INFO] Resource Backup %s created successfully", backupObject.Id)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"protecting"},
			Target:     []string{"available"},
			Refresh:    waitForCSBSBackupActive(backupClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      3 * time.Minute,
			MinTimeout: 3 * time.Minute,
		}
		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return diag.Errorf(
				"Error waiting for Backup (%s) to become available: %s",
				backupObject.Id, stateErr)
		}

	} else {
		return diag.Errorf("Error code: %s\n Error msg: %s", query[0].ErrorCode, query[0].ErrorMsg)
	}

	return resourceCSBSBackupV1Read(ctx, d, meta)

}

func resourceCSBSBackupV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(*cfg.Config)
	backupClient, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating csbs client: %s", err)
	}

	backupObject, err := backup.Get(backupClient, d.Id()).ExtractBackup()

	if err != nil {

		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Removing backup %s as it's already gone", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving backup: %s", err)

	}

	d.Set("resource_id", backupObject.ResourceId)
	d.Set("backup_name", backupObject.Name)
	d.Set("description", backupObject.Description)
	d.Set("resource_type", backupObject.ResourceType)
	d.Set("status", backupObject.Status)
	d.Set("volume_backups", flattenCSBSVolumeBackups(backupObject))
	d.Set("vm_metadata", flattenCSBSVMMetadata(backupObject))
	d.Set("backup_record_id", backupObject.CheckpointId)

	if err := d.Set("tags", flattenCSBSTags(backupObject)); err != nil {
		return diag.FromErr(err)
	}

	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceCSBSBackupV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	backupClient, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating csbs client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForCSBSBackupDelete(backupClient, d.Id(), d.Get("backup_record_id").(string)),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return diag.Errorf("Error deleting csbs backup: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCSBSBackupActive(backupClient *golangsdk.ServiceClient, backupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := backup.Get(backupClient, backupId).ExtractBackup()
		if err != nil {
			return nil, "", err
		}

		if n.Id == "error" {
			return nil, "", fmt.Errorf("Backup status: %s", n.Status)
		}

		return n, n.Status, nil
	}
}

func waitForCSBSBackupDelete(backupClient *golangsdk.ServiceClient, backupId string, backupRecordID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := backup.Get(backupClient, backupId).ExtractBackup()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted csbs backup %s", backupId)
				return r, "deleted", nil
			}
			return r, "deleting", err
		}

		err = backup.Delete(backupClient, backupRecordID).Err

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

func resourceCSBSTagsV1(d *schema.ResourceData) []backup.ResourceTag {
	rawTags := d.Get("tags").(*schema.Set).List()
	tags := make([]backup.ResourceTag, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tags[i] = backup.ResourceTag{
			Key:   rawMap["key"].(string),
			Value: rawMap["value"].(string),
		}
	}
	return tags
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

func flattenCSBSTags(backupObject *backup.Backup) []map[string]interface{} {
	var t []map[string]interface{}
	for _, tag := range backupObject.Tags {
		mapping := map[string]interface{}{
			"key":   tag.Key,
			"value": tag.Value,
		}
		t = append(t, mapping)
	}

	return t
}
