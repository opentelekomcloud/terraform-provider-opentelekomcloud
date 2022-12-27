package rds

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/backups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRDSv3Backup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRDSv3BackupRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				Computed: true,
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
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"databases": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRDSv3BackupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	opts := backups.ListOpts{
		InstanceID: d.Get("instance_id").(string),
		BackupID:   d.Get("backup_id").(string),
		BackupType: d.Get("type").(string),
	}

	backupList, err := backups.List(client, opts)
	if err != nil {
		return fmterr.Errorf("error listing backups: %w", err)
	}
	if len(backupList) < 1 {
		return common.DataSourceTooFewDiag
	}
	backup := backupList[0]

	d.SetId(backup.ID)
	mErr := multierror.Append(
		d.Set("name", backup.Name),
		d.Set("status", backup.Status),
		d.Set("type", backup.Type),
		d.Set("size", backup.Size),
		d.Set("begin_time", backup.BeginTime),
		d.Set("end_time", backup.EndTime),
		d.Set("databases", expandDatabases(backup.Databases)),
		d.Set("db_version", backup.Datastore.Version),
		d.Set("db_type", backup.Datastore.Type),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting RDSv3 instance backup fields: %w", err)
	}
	return nil
}

func expandDatabases(dbs []backups.BackupDatabase) []string {
	res := make([]string, len(dbs))
	for i, db := range dbs {
		res[i] = db.Name
	}
	return res
}
