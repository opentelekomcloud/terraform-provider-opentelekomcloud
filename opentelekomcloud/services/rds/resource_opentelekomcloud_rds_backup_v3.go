package rds

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	backups "github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/backups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceRdsBackupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRDSv3BackupCreate,
		ReadContext:   resourceRDSv3BackupRead,
		DeleteContext: resourceRDSv3BackupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"auto", "manual", "fragment", "incremental"},
					false,
				),
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceRDSv3BackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	opts := backups.CreateOpts{
		InstanceID:  d.Get("instance_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Databases:   resourceDatabaseExpand(d),
	}

	// check if rds instance exists
	rds, err := instances.List(client, instances.ListOpts{
		Id: opts.InstanceID,
	})
	if err != nil {
		return fmterr.Errorf("error getting RDSv3 instance: %w", err)
	}

	// wait until rds instance is active
	err = instances.WaitForStateAvailable(client, 120, opts.InstanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(rds.Instances) == 0 {
		return fmterr.Errorf("RDSv3 instance not found")
	}

	backup, err := backups.Create(client, opts)
	if err != nil {
		fmterr.Errorf("error creating new RDSv3 backup: %w", err)
	}

	err = backups.WaitForBackup(client, backup.InstanceID, backup.ID, "COMPLETED")
	if err != nil {
		return diag.FromErr(err)
	}

	if backup == nil {
		return diag.Errorf("backup not created")
	}
	// Check if backup.InstanceID and backup.ID are not empty
	if backup.InstanceID == "" || backup.ID == "" {
		return diag.Errorf("backup instance ID or backup ID is empty")
	}

	log.Printf("[DEBUG] RDSv3 backup created: %#v", backup)
	d.SetId(backup.ID)

	return resourceRDSv3BackupRead(ctx, d, meta)
}

func resourceRDSv3BackupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	err = backups.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud RDSv3 backup: %s", err)
	}
	d.SetId("")
	return nil
}

func resourceRDSv3BackupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	opts := backups.ListOpts{
		InstanceID: d.Get("instance_id").(string),
		BackupID:   d.Id(),
	}

	backupList, err := backups.List(client, opts)
	if err != nil {
		return fmterr.Errorf("error listing backups: %w", err)
	}
	if len(backupList) == 0 {
		return common.DataSourceTooFewDiag
	}
	backup := backupList[0]

	d.SetId(backup.ID)
	mErr := multierror.Append(
		d.Set("instance_id", backup.InstanceID),
		d.Set("name", backup.Name),
		d.Set("description", backup.Description),
		d.Set("databases", expandDatabases(backup.Databases)),
		d.Set("status", backup.Status),
		d.Set("type", backup.Type),
		d.Set("begin_time", backup.BeginTime),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting RDSv3 instance backup fields: %w", err)
	}
	return nil
}

func resourceDatabaseExpand(d *schema.ResourceData) []backups.BackupDatabase {
	var backupsDatabases []backups.BackupDatabase
	dbRaw := d.Get("databases").([]interface{})
	log.Printf("[DEBUG] dbRaw: %+v", dbRaw)
	for i := range dbRaw {
		db := dbRaw[i].(map[string]interface{})
		dbReq := backups.BackupDatabase{
			Name: db["name"].(string),
		}
		backupsDatabases = append(backupsDatabases, dbReq)
	}
	log.Printf("[DEBUG] backupsDatabases: %+v", backupsDatabases)
	return backupsDatabases
}
