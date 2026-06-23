package cfw

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	logs "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/logs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCfwLogConfigurationV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCfwLogConfigurationV1Create,
		ReadContext:   resourceCfwLogConfigurationV1Read,
		UpdateContext: resourceCfwLogConfigurationV1Update,
		DeleteContext: resourceCfwLogConfigurationV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lts_enable": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"lts_log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lts_attack_log_stream_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lts_attack_log_stream_enable": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"lts_access_log_stream_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lts_access_log_stream_enable": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"lts_flow_log_stream_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lts_flow_log_stream_enable": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCfwLogConfigurationV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	ltsEnable := d.Get("lts_enable").(int)
	createOpts := logs.LogConfigOpts{
		FWInstanceID:             d.Get("fw_instance_id").(string),
		LtsEnable:                &ltsEnable,
		LtsLogGroupID:            d.Get("lts_log_group_id").(string),
		LtsAttackLogStreamID:     d.Get("lts_attack_log_stream_id").(string),
		LtsAttackLogStreamEnable: d.Get("lts_attack_log_stream_enable").(int),
		LtsAccessLogStreamID:     d.Get("lts_access_log_stream_id").(string),
		LtsAccessLogStreamEnable: d.Get("lts_access_log_stream_enable").(int),
		LtsFlowLogStreamID:       d.Get("lts_flow_log_stream_id").(string),
		LtsFlowLogStreamEnable:   d.Get("lts_flow_log_stream_enable").(int),
	}

	log.Printf("[DEBUG] Create CFW log configuration options: %#v", createOpts)

	_, err = logs.CreateLogConfig(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud CFW log configuration: %w", err)
	}

	d.SetId(d.Get("fw_instance_id").(string))

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCfwLogConfigurationV1Read(clientCtx, d, meta)
}

func resourceCfwLogConfigurationV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	queryOpts := logs.QueryParameters{
		FwInstanceID: d.Id(),
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		queryOpts.EnterpriseProjectID = v.(string)
	}

	logConfig, err := logs.GetLogConfig(client, queryOpts)
	if err != nil {
		return fmterr.Errorf("error reading OpenTelekomCloud CFW log configuration: %w", err)
	}

	log.Printf("[DEBUG] Retrieved CFW log configuration %s: %#v", d.Id(), logConfig)

	mErr := multierror.Append(nil,
		d.Set("fw_instance_id", logConfig.FWInstanceID),
		d.Set("lts_enable", logConfig.LtsEnable),
		d.Set("lts_log_group_id", logConfig.LtsLogGroupID),
		d.Set("lts_attack_log_stream_id", logConfig.LtsAttackLogStreamID),
		d.Set("lts_attack_log_stream_enable", logConfig.LtsAttackLogStreamEnable),
		d.Set("lts_access_log_stream_id", logConfig.LtsAccessLogStreamID),
		d.Set("lts_access_log_stream_enable", logConfig.LtsAccessLogStreamEnable),
		d.Set("lts_flow_log_stream_id", logConfig.LtsFlowLogStreamID),
		d.Set("lts_flow_log_stream_enable", logConfig.LtsFlowLogStreamEnable),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCfwLogConfigurationV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	ltsEnable := d.Get("lts_enable").(int)
	updateOpts := logs.LogConfigOpts{
		FWInstanceID:             d.Get("fw_instance_id").(string),
		LtsEnable:                &ltsEnable,
		LtsLogGroupID:            d.Get("lts_log_group_id").(string),
		LtsAttackLogStreamID:     d.Get("lts_attack_log_stream_id").(string),
		LtsAttackLogStreamEnable: d.Get("lts_attack_log_stream_enable").(int),
		LtsAccessLogStreamID:     d.Get("lts_access_log_stream_id").(string),
		LtsAccessLogStreamEnable: d.Get("lts_access_log_stream_enable").(int),
		LtsFlowLogStreamID:       d.Get("lts_flow_log_stream_id").(string),
		LtsFlowLogStreamEnable:   d.Get("lts_flow_log_stream_enable").(int),
	}

	log.Printf("[DEBUG] Update CFW log configuration options: %#v", updateOpts)

	_, err = logs.UpdateLogConfig(client, updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud CFW log configuration: %w", err)
	}

	log.Printf("[DEBUG] CFW log configuration %s updated successfully", d.Id())

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCfwLogConfigurationV1Read(clientCtx, d, meta)
}

func resourceCfwLogConfigurationV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The action resource is only removed from the state the task remains in the cloud.",
		},
	}
}
