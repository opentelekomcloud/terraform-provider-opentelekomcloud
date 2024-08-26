package rds

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceRdsMaintenanceV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsMaintenanceV3Create,
		ReadContext:   resourceRdsMaintenanceV3Read,
		DeleteContext: resourceRdsMaintenanceV3Delete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRdsMaintenanceV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	instanceId := d.Get("instance_id").(string)

	updateOpts := instances.ChangeOpsWindowOpts{
		InstanceId: instanceId,
		StartTime:  d.Get("start_time").(string),
		EndTime:    d.Get("end_time").(string),
	}

	err = instances.ChangeOpsWindow(client, updateOpts)
	if err != nil {
		return diag.Errorf("error setting maintenance window for RDS instance %s: %s", instanceId, err)
	}

	d.SetId(instanceId)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceRdsMaintenanceV3Read(clientCtx, d, meta)
}

func resourceRdsMaintenanceV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	rdsInstance, err := GetRdsInstance(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching RDS instance: %s", err)
	}
	if rdsInstance == nil {
		d.SetId("")
		return nil
	}

	times := strings.Split(rdsInstance.MaintenanceWindow, "-")
	if len(times) != 2 {
		return fmterr.Errorf("Invalid maintenance time window.")
	}

	mErr := multierror.Append(
		d.Set("start_time", times[0]),
		d.Set("end_time", times[1]),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceRdsMaintenanceV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	// set time to default maintenance window (22:00-00:00)
	mOpts := instances.ChangeOpsWindowOpts{
		InstanceId: d.Get("instance_id").(string),
		StartTime:  "22:00",
		EndTime:    "00:00",
	}
	err = instances.ChangeOpsWindow(client, mOpts)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud RDSv3 maintenance window: %s", err)
	}

	d.SetId("")
	return nil
}
