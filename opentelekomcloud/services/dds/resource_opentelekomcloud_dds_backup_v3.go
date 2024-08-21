package dds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/backups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDdsBackupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsBackupV3Create,
		ReadContext:   resourceDdsBackupV3Read,
		DeleteContext: resourceDdsBackupV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: ddsBackupV3ImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			"instance_name": {
				Type:     schema.TypeString,
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
						"storage_engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"type": {
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
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDdsBackupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	instanceId := d.Get("instance_id").(string)
	opts := backups.CreateOpts{
		Backup: &backups.Backup{
			InstanceId:  instanceId,
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
	}
	retryFunc := func() (interface{}, bool, error) {
		resp, err := backups.Create(client, opts)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating DDS backup: %s", err)
	}
	body := r.(*backups.Job)
	d.SetId(body.BackupId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      JobStateRefreshFunc(client, body.JobId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for job (%s) to complete: %s", body.BackupId, err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceDdsBackupV3Read(clientCtx, d, meta)
}

func resourceDdsBackupV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	listOpts := backups.ListBackupsOpts{
		BackupId: d.Id(),
	}
	backupList, err := backups.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error fetching DDS backups: %w", err)
	}
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", backupList.Backups[0].Name),
		d.Set("instance_id", backupList.Backups[0].InstanceId),
		d.Set("instance_name", backupList.Backups[0].InstanceName),
		d.Set("type", backupList.Backups[0].Type),
		d.Set("begin_time", backupList.Backups[0].BeginTime),
		d.Set("end_time", backupList.Backups[0].EndTime),
		d.Set("status", backupList.Backups[0].Status),
		d.Set("size", backupList.Backups[0].Size),
		d.Set("description", backupList.Backups[0].Description),
	)

	datastoreList := make([]map[string]interface{}, 0, 1)
	datastore := map[string]interface{}{
		"type":           backupList.Backups[0].Datastore.Type,
		"version":        backupList.Backups[0].Datastore.Version,
		"storage_engine": backupList.Backups[0].Datastore.StorageEngine,
	}
	datastoreList = append(datastoreList, datastore)
	if err := d.Set("datastore", datastoreList); err != nil {
		return fmterr.Errorf("error setting DDSv3 datastore opts: %w", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDdsBackupV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}
	instanceId := d.Get("instance_id").(string)
	retryFunc := func() (interface{}, bool, error) {
		resp, err := backups.Delete(client, d.Id())
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 5 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, parseDdsBackupError(err), "error deleting DDS backup")
	}

	body := r.(*backups.Job)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed", "Deleted"},
		Refresh:      JobStateRefreshFunc(client, body.JobId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for job (%s) to be completed: %s", d.Id(), err)
	}

	return nil
}

func parseDdsBackupError(err error) error {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return err
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return err
		}

		if errorCode == "DBS.201502" || errorCode == "DBS.201214" {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return err
}

func ddsBackupV3ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <instance_id>/<id>")
	}
	instanceId := parts[0]
	backupId := parts[1]
	d.SetId(backupId)
	err := d.Set("instance_id", instanceId)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
