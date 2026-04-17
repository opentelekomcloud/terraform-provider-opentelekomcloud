package taurusdb

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/backup"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceTaurusDbMysqlBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDbMysqlBackupCreate,
		ReadContext:   resourceTaurusDbMysqlBackupRead,
		DeleteContext: resourceTaurusDbMysqlBackupDelete,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"take_up_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTaurusDbMysqlBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	createOpts := backup.CreateOpts{
		InstanceId:  d.Get("instance_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := backup.Create(client, createOpts)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     TaurusDbInstanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating TaurusDb MySQL backup: %s", err)
	}

	createResp := r.(*backup.CreateResponse)

	d.SetId(createResp.Backup.Id)

	err = waitForTaurusDbBackupComplete(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTaurusDbMysqlBackupRead(ctx, d, meta)
}

func waitForTaurusDbBackupComplete(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"BUILDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      taurusDbBackupStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for TaurusDb MySQL backup (%s) to build complete: %s", d.Id(), err)
	}
	return nil
}

func taurusDbBackupStateRefreshFunc(client *golangsdk.ServiceClient, backupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		backup, err := getTaurusDbBackup(client, backupID)
		if err != nil {
			return nil, "", err
		}

		return backup, backup.Status, nil
	}
}

func resourceTaurusDbMysqlBackupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	var mErr *multierror.Error

	tbBackup, err := getTaurusDbBackup(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving TaurusDb MySQL backup")
	}
	mErr = multierror.Append(
		mErr,
		d.Set("region", config.GetRegion(d)),
		d.Set("instance_id", tbBackup.InstanceId),
		d.Set("name", tbBackup.Name),
		d.Set("description", tbBackup.Description),
		d.Set("begin_time", tbBackup.BeginTime),
		d.Set("end_time", tbBackup.EndTime),
		d.Set("take_up_time", tbBackup.TakeUpTime),
		d.Set("size", tbBackup.Size),
		d.Set("datastore", flattenGetTaurusDbBackupResponseBodyDatastore(tbBackup)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetTaurusDbBackupResponseBodyDatastore(backup *backup.BackupListInfo) []interface{} {
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"type":    backup.Datastore.Type,
		"version": backup.Datastore.Version,
	})
	return rst
}

func getTaurusDbBackup(client *golangsdk.ServiceClient, backupID string) (*backup.BackupListInfo, error) {
	getResp, err := backup.List(client,
		backup.ListOpts{BackupId: backupID},
	)

	if err != nil {
		return nil, fmt.Errorf("error retrieving TaurusDb MySQL backup: %s", err)
	}

	if len(getResp) == 0 {
		return nil, fmt.Errorf("TaurusDb backup doesn't exist")
	}

	return &getResp[0], nil
}

func resourceTaurusDbMysqlBackupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	_, err = backup.Delete(client, d.Id())

	if err != nil {
		return diag.Errorf("error deleting TaurusDb MySQL backup: %s", err)
	}

	return nil
}
