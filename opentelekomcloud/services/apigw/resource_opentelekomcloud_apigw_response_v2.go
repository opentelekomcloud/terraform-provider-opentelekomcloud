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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/response"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIResponseV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResponseV2Create,
		ReadContext:   resourceResponseV2Read,
		UpdateContext: resourceResponseV2Update,
		DeleteContext: resourceResponseV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCustomResponseV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[\w-]*$`),
						"Only letters, digits, hyphens(-), and underscores (_) are allowed."),
					validation.StringLenBetween(1, 64),
				),
			},
			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"body": {
							Type:     schema.TypeString,
							Required: true,
						},
						"status_code": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(200, 599),
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
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

func buildCustomResponses(respSet *schema.Set) map[string]response.ResponseInfo {
	result := make(map[string]response.ResponseInfo)

	for _, resp := range respSet.List() {
		rule := resp.(map[string]interface{})

		result[rule["error_type"].(string)] = response.ResponseInfo{
			Body:   rule["body"].(string),
			Status: rule["status_code"].(int),
		}
	}

	return result
}

func resourceResponseV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := response.CreateOpts{
		GatewayID: d.Get("gateway_id").(string),
		GroupId:   d.Get("group_id").(string),
		Name:      d.Get("name").(string),
		Responses: buildCustomResponses(d.Get("rule").(*schema.Set)),
	}
	resp, err := response.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud APIGW custom response: %s", err)
	}
	d.SetId(resp.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceResponseV2Read(clientCtx, d, meta)
}

func flattenCustomResponses(respMap map[string]response.ResponseInfo) []map[string]interface{} {
	if len(respMap) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(respMap))
	for errorType, rule := range respMap {
		if rule.IsDefault {
			continue
		}
		result = append(result, map[string]interface{}{
			"error_type":  errorType,
			"body":        rule.Body,
			"status_code": rule.Status,
		})
	}
	return result
}

func resourceResponseV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	groupId := d.Get("group_id").(string)

	resp, err := response.Get(client, gatewayId, groupId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud APIGW custom response")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("rule", flattenCustomResponses(resp.Responses)),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud APIGW custom response (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceResponseV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opt := response.CreateOpts{
		GatewayID: d.Get("gateway_id").(string),
		GroupId:   d.Get("group_id").(string),
		Name:      d.Get("name").(string),
		Responses: buildCustomResponses(d.Get("rule").(*schema.Set)),
	}
	_, err = response.Update(client, d.Id(), opt)
	if err != nil {
		return diag.Errorf("error updating OpenTelekomCloud APIGW custom response: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceResponseV2Read(clientCtx, d, meta)
}

func resourceResponseV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	groupId := d.Get("group_id").(string)
	err = response.Delete(client, gatewayId, groupId, d.Id())
	if err != nil {
		return diag.Errorf("error deleting OpenTelekomCloud APIGW custom response (%s) from the dedicated group (%s): %s",
			d.Id(), groupId, err)
	}
	return nil
}

func resourceCustomResponseV2ImportState(ctx context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import IDs and name, " +
			"must be <gateway_id>/<group_id>/<name>")
	}

	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf(errCreationV2Client, err)
	}

	gatewayId := parts[0]
	groupId := parts[1]
	name := parts[2]
	opts := response.ListOpts{
		GatewayID: gatewayId,
		GroupID:   groupId,
	}
	resp, err := response.List(client, opts)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error getting custom response list from server: %s", err)
	}
	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmt.Errorf("unable to find any custom response from server")
	}
	for _, val := range resp {
		if val.Name == name {
			d.SetId(val.ID)
			mErr := multierror.Append(nil,
				d.Set("gateway_id", gatewayId),
				d.Set("group_id", groupId),
			)
			return []*schema.ResourceData{d}, mErr.ErrorOrNil()
		}
	}
	return []*schema.ResourceData{d}, fmt.Errorf("unable to find the custom response (%s) from server", name)
}
