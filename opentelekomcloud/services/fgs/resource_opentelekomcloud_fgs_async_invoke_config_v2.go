package fgs

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/async_config"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAsyncInvokeConfigurationV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAsyncInvokeConfigurationV2Create,
		ReadContext:   resourceAsyncInvokeConfigurationV2Read,
		UpdateContext: resourceAsyncInvokeConfigurationV2Update,
		DeleteContext: resourceAsyncInvokeConfigurationV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAsyncInvokeConfigImportState,
		},

		Schema: map[string]*schema.Schema{
			"function_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_async_event_age_in_seconds": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_async_retry_attempts": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"on_success": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     destinationConfigSchemaResource(),
			},
			"on_failure": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     destinationConfigSchemaResource(),
			},
			// Setting this param leads to Internal error
			// "enable_async_status_log": {
			// 	Type:     schema.TypeBool,
			// 	Optional: true,
			// },
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func destinationConfigSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},
			"param": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func modifyAsyncInvokeConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		opts = async_config.UpdateOpts{
			FuncUrn:     d.Get("function_urn").(string),
			MaxEventAge: pointerto.Int(d.Get("max_async_event_age_in_seconds").(int)),
			MaxRetry:    pointerto.Int(d.Get("max_async_retry_attempts").(int)),
			// EnableStatusLog: pointerto.Bool(d.Get("enable_async_status_log").(bool)),
		}
		destinationConfig = async_config.DestinationConfig{}
	)

	if successConfigs, ok := d.GetOk("on_success"); ok {
		raws := successConfigs.([]interface{})
		cfgDetails := raws[0].(map[string]interface{})
		destinationConfig.OnSuccess = &async_config.Destination{
			Destination: cfgDetails["destination"].(string),
			Param:       cfgDetails["param"].(string),
		}
	}
	if failureConfigs, ok := d.GetOk("on_failure"); ok {
		raws := failureConfigs.([]interface{})
		cfgDetails := raws[0].(map[string]interface{})
		destinationConfig.OnFailure = &async_config.Destination{
			Destination: cfgDetails["destination"].(string),
			Param:       cfgDetails["param"].(string),
		}
	}
	if destinationConfig != (async_config.DestinationConfig{}) {
		opts.DestinationConfig = &destinationConfig
	}
	_, err := async_config.Update(client, opts)
	if err != nil {
		return fmt.Errorf("error modifying the async invoke configuration: %s", err)
	}
	return nil
}

func resourceAsyncInvokeConfigurationV2Create(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = modifyAsyncInvokeConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error creating the configuration of the asynchronous invocation: %s", err)
	}
	d.SetId(d.Get("function_urn").(string))

	clientCtx := common.CtxWithClient(ctx, client, fgsClientV2)
	return resourceAsyncInvokeConfigurationV2Read(clientCtx, d, meta)
}

func flattenDestinationConfig(destConfig async_config.Destination) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"destination": destConfig.Destination,
			"param":       destConfig.Param,
		},
	}
}

func resourceAsyncInvokeConfigurationV2Read(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	resp, err := async_config.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "asynchronous invocation configuration")
	}
	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("max_async_event_age_in_seconds", resp.MaxEventAge),
		d.Set("max_async_retry_attempts", resp.MaxRetry),
		d.Set("on_success", flattenDestinationConfig(*resp.DestinationConfig.OnSuccess)),
		d.Set("on_failure", flattenDestinationConfig(*resp.DestinationConfig.OnFailure)),
		// d.Set("enable_async_status_log", resp.EnableStatusLog),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving asynchronous invocation configuration fields: %s", mErr)
	}

	return nil
}

func resourceAsyncInvokeConfigurationV2Update(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = modifyAsyncInvokeConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error updating the configuration of the asynchronous invocation: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, fgsClientV2)
	return resourceAsyncInvokeConfigurationV2Read(clientCtx, d, meta)
}

func resourceAsyncInvokeConfigurationV2Delete(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = async_config.Delete(client, d.Id())
	if err != nil {
		return diag.Errorf("error deleting the configuration of the asynchronous invocation: %s", err)
	}
	return nil
}

func resourceAsyncInvokeConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	return []*schema.ResourceData{d}, d.Set("function_urn", d.Id())
}
