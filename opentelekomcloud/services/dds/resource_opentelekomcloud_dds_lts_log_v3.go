package dds

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/logs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDdsLtsLogV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsLtsLogV3Create,
		ReadContext:   resourceDdsLtsLogV3Read,
		UpdateContext: resourceDdsLtsLogV3Create,
		DeleteContext: resourceDdsLtsLogV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lts_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lts_stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDdsLtsLogV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}
	instanceID := d.Get("instance_id").(string)
	createOpts := logs.CreateOpts{
		LtsConfigs: []logs.Configs{
			{
				InstanceID:  instanceID,
				LogType:     d.Get("log_type").(string),
				LtsGroupId:  d.Get("lts_group_id").(string),
				LtsStreamId: d.Get("lts_stream_id").(string),
			},
		},
	}
	retryFunc := func() (interface{}, bool, error) {
		err = logs.Create(client, createOpts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error associating DDS with LTS log: %s", err)
	}

	d.SetId(instanceID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating"},
		Target:     []string{"normal"},
		Refresh:    instanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to become ready: %w", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceDdsLtsLogV3Read(clientCtx, d, meta)
}

func resourceDdsLtsLogV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	list, err := logs.List(client, logs.ListOpts{})
	if err != nil {
		return diag.Errorf("error retrieving DDS LTS configs: %s", err)
	}

	var listResult logs.InstanceLtsConfigs
	for _, lr := range list.InstanceLtsConfigs {
		if lr.Instance.ID == d.Id() {
			listResult = lr
			break
		}
	}
	if len(listResult.LtsConfigs) == 0 {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find OpenTelekomCloud LTS v2 log group by its ID (%s)", d.Id()))
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("instance_id", d.Id()),
		d.Set("log_type", listResult.LtsConfigs[0].LogType),
		d.Set("lts_group_id", listResult.LtsConfigs[0].LtsGroupId),
		d.Set("lts_stream_id", listResult.LtsConfigs[0].LtsStreamId),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDdsLtsLogV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	instanceID := d.Get("instance_id").(string)
	logType := d.Get("log_type").(string)

	retryFunc := func() (interface{}, bool, error) {
		err = logs.Delete(client, logs.DeleteOpts{
			LtsConfigs: []logs.LtsConfig{
				{
					InstanceID: instanceID,
					LogType:    logType,
				},
			},
		})
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error disassociating DDS with LTS log: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating"},
		Target:     []string{"normal"},
		Refresh:    instanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to become ready: %w", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceDdsLtsLogV3Read(clientCtx, d, meta)
}
