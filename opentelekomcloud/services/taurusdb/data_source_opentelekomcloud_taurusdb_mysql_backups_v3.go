package taurusdb

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/backup"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceTaurusDBV3MysqlBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBMysqlBackupsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
			},
			"backup_id": {
				Type:        schema.TypeString,
				Optional:    true,
			},
			"backup_type": {
				Type:        schema.TypeString,
				Optional:    true,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
			},
			"backups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
						},
						"begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
						},
						"take_up_time": {
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
						},
						"size": {
							Type:        schema.TypeFloat,
							Computed:    true,
						},
						"datastore": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaurusDBMysqlBackupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	opts := backup.ListOpts{
		InstanceId: d.Get("instance_id").(string),
		BackupId:   d.Get("backup_id").(string),
		BackupType: d.Get("backup_type").(string),
		BeginTime:  d.Get("begin_time").(string),
		EndTime:    d.Get("end_time").(string),
	}

	backupList, err := backup.List(client, opts)
	if err != nil {
		return fmterr.Errorf("error retrieving TaurusDB MySQL backups: %w", err)
	}

	log.Printf("[DEBUG] Retrieved %d TaurusDB MySQL backups", len(backupList))

	backupsMapped := make([]map[string]interface{}, len(backupList))
	for i, b := range backupList {
		backupsMapped[i] = map[string]interface{}{
			"id":           b.Id,
			"name":         b.Name,
			"instance_id":  b.InstanceId,
			"begin_time":   b.BeginTime,
			"end_time":     b.EndTime,
			"take_up_time": b.TakeUpTime,
			"type":         b.Type,
			"size":         b.Size,
			"datastore": []map[string]interface{}{
				{
					"type":    b.Datastore.Type,
					"version": b.Datastore.Version,
				},
			},
			"status":      b.Status,
			"description": b.Description,
		}
	}

	d.SetId(config.GetRegion(d))
	mErr := multierror.Append(nil,
		d.Set("backups", backupsMapped),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting TaurusDB MySQL backups fields: %w", err)
	}

	return nil
}
