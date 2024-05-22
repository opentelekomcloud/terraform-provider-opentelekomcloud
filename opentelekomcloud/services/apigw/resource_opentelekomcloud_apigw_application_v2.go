package apigw

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/app"
	appcode "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/app_code"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIApplicationV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationV2Create,
		ReadContext:   resourceApplicationV2Read,
		UpdateContext: resourceApplicationV2Update,
		DeleteContext: resourceApplicationV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 255),
				),
			},
			"app_codes": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.All(
						validation.StringMatch(
							regexp.MustCompile(`^[A-Za-z0-9+=][\w!@#$%+-/=]*$`),
							"The code consists of 64 to 180 characters, starting with a letter, digit, "+
								"plus sign (+) or slash (/). Only letters, digits and following special special "+
								"characters are allowed: !@#$%+-_/="),
						validation.StringLenBetween(64, 180),
					),
				},
			},
			"secret_action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(SecretActionReset),
				}, false),
			},
			"registration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createApplicationCodes(client *golangsdk.ServiceClient, gatewayId, appId string, codes []interface{}) error {
	for _, code := range codes {
		opts := appcode.CreateOpts{
			GatewayID: gatewayId,
			AppID:     appId,
			AppCode:   code.(string),
		}
		_, err := appcode.Create(client, opts)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud APIGW application code: %s", err)
		}
	}
	return nil
}

func resourceApplicationV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	gatewayId := d.Get("gateway_id").(string)

	opts := app.CreateOpts{
		GatewayID:   gatewayId,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	resp, err := app.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud APIGW dedicated application: %s", err)
	}
	d.SetId(resp.ID)

	if v, ok := d.GetOk("app_codes"); ok {
		if err := createApplicationCodes(client, gatewayId, d.Id(), v.(*schema.Set).List()); err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceApplicationV2Read(clientCtx, d, meta)
}

func queryApplicationCodes(client *golangsdk.ServiceClient, gatewayId, appId string) ([]appcode.CodeResp, error) {
	codes, err := appcode.ListAppCodesOfApp(client, appcode.ListAppsOpts{
		GatewayID: gatewayId,
		AppID:     appId,
	})
	if err != nil {
		return nil, err
	}
	return codes, nil
}

func flattenApplicationCodes(codes []appcode.CodeResp) []interface{} {
	if len(codes) < 1 {
		return nil
	}
	result := make([]interface{}, len(codes))
	for i, v := range codes {
		result[i] = v.AppCode
	}
	return result
}

func resourceApplicationV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	resp, err := app.Get(client, d.Get("gateway_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dedicated application")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("registration_time", resp.RegisterTime),
		d.Set("updated_at", resp.UpdateTime),
		d.Set("app_key", resp.AppKey),
		d.Set("app_secret", resp.AppSecret),
	)
	if codes, err := queryApplicationCodes(client, d.Get("gateway_id").(string), d.Id()); err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		// The application code is sort by create time on server, not code.
		mErr = multierror.Append(d.Set("app_codes",
			schema.NewSet(schema.HashString, flattenApplicationCodes(codes))))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud APIGW dedicated application fields: %s", err)
	}
	return nil
}

func isCodeInApplication(codes []appcode.CodeResp, code string) (string, bool) {
	for _, s := range codes {
		if s.AppCode == code {
			return s.ID, true
		}
	}
	return "", false
}

func removeApplicationCodes(client *golangsdk.ServiceClient, gatewayId, appId string, codes []interface{}) error {
	results, err := queryApplicationCodes(client, gatewayId, appId)
	if err != nil {
		return fmt.Errorf("error retrieving OpenTelekomCloud APIGW application codes: %s", err)
	}
	for _, code := range codes {
		codeId, ok := isCodeInApplication(results, code.(string))
		if !ok {
			continue
		}
		if err := appcode.Delete(client, gatewayId, appId, codeId); err != nil {
			return fmt.Errorf("error removing code (%v) form the OpenTelekomCloud APIGW application (%s) : %s", code, appId, err)
		}
	}
	return nil
}

func updateApplicationCodes(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	gatewayId := d.Get("gateway_id").(string)
	appId := d.Id()
	oldRaws, newRaws := d.GetChange("app_codes")

	addRaws := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	removeRaws := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	if removeRaws.Len() > 0 {
		if err := removeApplicationCodes(client, gatewayId, appId, removeRaws.List()); err != nil {
			return err
		}
	}

	if addRaws.Len() > 0 {
		if err := createApplicationCodes(client, gatewayId, appId, addRaws.List()); err != nil {
			return err
		}
	}
	return nil
}

func resourceApplicationV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	appId := d.Id()

	if d.HasChanges("name", "description") {
		opts := app.CreateOpts{
			GatewayID:   gatewayId,
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}
		_, err = app.Update(client, appId, opts)
		if err != nil {
			return diag.Errorf("error updating dedicated OpenTelekomCloud APIGW application (%s): %s", appId, err)
		}
	}
	if d.HasChange("app_codes") {
		err = updateApplicationCodes(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("secret_action") {
		if v, ok := d.GetOk("secret_action"); ok && v.(string) == string(SecretActionReset) {
			if _, err := app.ResetSecret(client,
				app.ResetOpts{GatewayID: d.Get("gateway_id").(string), AppID: d.Id()}); err != nil {
				return diag.Errorf("error resetting OpenTelekomCloud APIGW application secret: %s", err)
			}
		}
	}
	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceApplicationV2Read(clientCtx, d, meta)
}

func resourceApplicationV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	appId := d.Id()
	err = app.Delete(client, gatewayId, appId)
	if err != nil {
		return diag.Errorf("error deleting OpenTelekomCloud APIGW application (%s) from the gateway (%s): %s", appId, gatewayId, err)
	}

	return nil
}

func resourceApplicationV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <gateway_id>/<id>")
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("gateway_id", parts[0])
}
