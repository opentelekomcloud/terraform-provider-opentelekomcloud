package csbs

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/csbs/v1/backup"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCSBSBackupV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCSBSBackupV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
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
			"backup_name": {
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
			"backup_record_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_trigger": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"average_speed": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vm_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func dataSourceCSBSBackupV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	backupClient, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating CSBSv1 client: %w", err)
	}

	var refinedBackups []backup.Backup
	if v, ok := d.GetOk("id"); ok {
		bk, err := backup.Get(backupClient, v.(string))
		if err != nil {
			return fmterr.Errorf("unable to retrieve backup: %s", err)
		}
		refinedBackups = append(refinedBackups, *bk)
	} else {
		listOpts := backup.ListOpts{
			Name:         d.Get("backup_name").(string),
			Status:       d.Get("status").(string),
			ResourceName: d.Get("resource_name").(string),
			CheckpointId: d.Get("backup_record_id").(string),
			ResourceType: d.Get("resource_type").(string),
			ResourceId:   d.Get("resource_id").(string),
			PolicyId:     d.Get("policy_id").(string),
			VmIp:         d.Get("vm_ip").(string),
		}
		backUps, err := backup.List(backupClient, listOpts)
		if err != nil {
			return fmterr.Errorf("unable to retrieve backup: %s", err)
		}
		refinedBackups = backUps
	}

	if len(refinedBackups) < 1 {
		return common.DataSourceTooFewDiag
	}

	if len(refinedBackups) > 1 {
		return common.DataSourceTooManyDiag
	}

	backupObject := refinedBackups[0]
	log.Printf("[INFO] Retrieved backup %s using given filter", backupObject.Id)

	d.SetId(backupObject.Id)

	mErr := multierror.Append(
		d.Set("backup_record_id", backupObject.CheckpointId),
		d.Set("backup_name", backupObject.Name),
		d.Set("resource_id", backupObject.ResourceId),
		d.Set("status", backupObject.Status),
		d.Set("description", backupObject.Description),
		d.Set("resource_type", backupObject.ResourceType),
		d.Set("auto_trigger", backupObject.ExtendInfo.AutoTrigger),
		d.Set("average_speed", backupObject.ExtendInfo.AverageSpeed),
		d.Set("resource_name", backupObject.ExtendInfo.ResourceName),
		d.Set("size", backupObject.ExtendInfo.Size),
		d.Set("volume_backups", flattenCSBSVolumeBackups(&backupObject)),
		d.Set("vm_metadata", flattenCSBSVMMetadata(&backupObject)),
		d.Set("region", config.GetRegion(d)),
		d.Set("tags", flattenCSBSTags(backupObject.Tags)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
