package apigw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	appcode "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/app_code"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIAppcodeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppcodeV2Create,
		ReadContext:   resourceAppcodeV2Read,
		DeleteContext: resourceAppcodeV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppcodeV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppcodeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	appId := d.Get("application_id").(string)
	appCode := d.Get("value").(string)
	var resp *appcode.CodeResp
	if appCode != "" {
		opts := appcode.CreateOpts{
			GatewayID: gatewayId,
			AppID:     appId,
			AppCode:   appCode,
		}
		resp, err = appcode.Create(client, opts)
		if err != nil {
			return diag.Errorf("generating OpenTelekomCloud APIGW AppCode failed: %s", err)
		}
	} else {
		resp, err = appcode.GenerateAppCode(client, gatewayId, appId)
		if err != nil {
			return diag.Errorf("auto generating OpenTelekomCloud APIGW AppCode failed: %s", err)
		}
	}

	d.SetId(resp.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAppcodeV2Read(clientCtx, d, meta)
}

func resourceAppcodeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	appId := d.Get("application_id").(string)
	appCodeId := d.Id()

	resp, err := appcode.Get(client, gatewayId, appId, appCodeId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error querying OpenTelekomCloud APIGW AppCode (%s) from specified application (%s) under dedicated instance (%s)",
				appCodeId, appId, gatewayId))
	}
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("value", resp.AppCode),
		d.Set("created_at", resp.CreateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud APIGW AppCode resource fields: %s", err)
	}
	return nil
}

func resourceAppcodeV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	gatewayId := d.Get("gateway_id").(string)
	appId := d.Get("application_id").(string)
	appCodeId := d.Id()

	err = appcode.Delete(client, gatewayId, appId, appCodeId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error deleting OpenTelekomCloud APIGW AppCode (%s) from specified application (%s) under dedicated instance (%s)",
				appCodeId, appId, gatewayId))
	}
	return nil
}

func resourceAppcodeV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<application_id>/<id>', "+
			"but got '%s'", importedId)
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("gateway_id", parts[0]),
		d.Set("application_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving OpenTelekomCloud APIGW AppCode resource fields during import: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
